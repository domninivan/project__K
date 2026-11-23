package ru.mirea.production.planning;

import java.lang.reflect.Proxy;
import java.util.UUID;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicReference;
import javax.sql.DataSource;
import org.junit.jupiter.api.*;
import org.junit.jupiter.api.condition.EnabledIfEnvironmentVariable;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.util.ReflectionTestUtils;
import org.springframework.transaction.support.TransactionSynchronizationManager;
import org.springframework.web.client.ResourceAccessException;
import org.springframework.web.server.ResponseStatusException;
import ru.mirea.production.planning.client.GoPlanningClient;
import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.planning.exception.GoServiceUnavailableException;
import ru.mirea.production.planning.service.*;
import ru.mirea.production.security.AppPrincipal;
import ru.mirea.production.user.entity.Role;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
@EnabledIfEnvironmentVariable(named = "JAVA_MIGRATION_TEST_URL", matches = ".+")
class OrderRoutePostgresTest {
    @Autowired DataSource dataSource;
    @Autowired OrderRouteService routes;
    @Autowired PlanningServiceImpl planning;
    JdbcTemplate jdbc;

    @DynamicPropertySource
    static void properties(DynamicPropertyRegistry r) {
        r.add("DB_URL", () -> System.getenv("JAVA_MIGRATION_TEST_URL"));
        r.add("DB_USERNAME", () -> System.getenv().getOrDefault("JAVA_MIGRATION_TEST_USER", "postgres"));
        r.add("DB_PASSWORD", () -> System.getenv().getOrDefault("JAVA_MIGRATION_TEST_PASSWORD", ""));
        r.add("JWT_SECRET_BASE64", () -> java.util.Base64.getEncoder().encodeToString(new byte[32]));
    }
    @BeforeEach void setup() { jdbc = new JdbcTemplate(dataSource); asPlanner(); }
    @AfterEach void clear() { SecurityContextHolder.clearContext(); }

    @Test
    void snapshotSurvivesCatalogueChangesAndGoFailureWithoutNetworkTransaction() {
        Fixture fixture = fixture();
        Object original = ReflectionTestUtils.getField(planning, "goPlanningClient");
        AtomicReference<PlanningRequest> sent = new AtomicReference<>();
        GoPlanningClient fake = (GoPlanningClient) Proxy.newProxyInstance(getClass().getClassLoader(),
                new Class<?>[]{GoPlanningClient.class}, (proxy, method, args) -> {
                    assertFalse(TransactionSynchronizationManager.isActualTransactionActive());
                    sent.set((PlanningRequest) args[0]);
                    assertEquals(2, count(fixture.order())); // visible after snapshot transaction committed
                    throw new ResourceAccessException("Simulated upstream failure");
                });
        ReflectionTestUtils.setField(planning, "goPlanningClient", fake);
        try {
            assertThrows(GoServiceUnavailableException.class, () -> planning.planOrder(fixture.order()));
            var first = sent.get();
            assertEquals(java.time.LocalDate.of(2099, 1, 1), first.dueDate());
            assertEquals("MEDIUM", first.priority().name());
            assertEquals(fixture.operation().longValue(), first.items().get(0).operations().get(0).operationId());
            assertEquals(600L, first.items().get(0).operations().get(0).timePerUnitSeconds());
            assertEquals("Точение", jdbc.queryForObject(
                    "SELECT o.title FROM order_item_operation o JOIN production_order_item i ON i.id = o.order_item_id "
                            + "WHERE i.order_id = ? AND o.source_operation_id = ?",
                    String.class, fixture.order(), fixture.operation()));
            jdbc.update("UPDATE product_operation SET time_per_unit_seconds = 700 WHERE id = ?", fixture.operation());
            assertThrows(GoServiceUnavailableException.class, () -> planning.planOrder(fixture.order()));
            assertEquals(first.items(), sent.get().items());
            assertEquals(2, count(fixture.order()));
            jdbc.update("DELETE FROM product_operation WHERE product_id = ?", fixture.product());
            assertEquals(first.items(), routes.preparePlan(fixture.order()).items());
        } finally {
            ReflectionTestUtils.setField(planning, "goPlanningClient", original);
        }
    }

    @Test
    void concurrentPreparationCreatesOnlyOneRoute() throws Exception {
        Fixture fixture = fixture();
        try (var executor = Executors.newFixedThreadPool(2)) {
            Callable<PlanningRequest> prepare = () -> {
                asPlanner();
                try { return routes.preparePlan(fixture.order()); }
                finally { SecurityContextHolder.clearContext(); }
            };
            Future<PlanningRequest> a = executor.submit(prepare);
            Future<PlanningRequest> b = executor.submit(prepare);
            assertEquals(a.get(20, TimeUnit.SECONDS).items(), b.get(20, TimeUnit.SECONDS).items());
            assertEquals(2, count(fixture.order()));
        }
    }

    @Test
    void invalidRouteRollsBackAndClosedOrderCannotCalculate() {
        Fixture fixture = fixture();
        jdbc.update("UPDATE product_operation SET operation_type = 'UNKNOWN' WHERE id = ?", fixture.operation());
        assertThrows(ResponseStatusException.class, () -> routes.preparePlan(fixture.order()));
        assertEquals(0, count(fixture.order()));
        jdbc.update("UPDATE product_operation SET operation_type = 'TURNING' WHERE id = ?", fixture.operation());
        jdbc.update("UPDATE production_order SET status = 'CANCELLED' WHERE id = ?", fixture.order());
        assertThrows(ResponseStatusException.class, () -> routes.preparePlan(fixture.order()));
        assertEquals(0, count(fixture.order()));
    }

    private int count(int order) {
        return jdbc.queryForObject("SELECT count(*) FROM order_item_operation o JOIN production_order_item i "
                + "ON i.id = o.order_item_id WHERE i.order_id = ?", Integer.class, order);
    }
    private Fixture fixture() {
        int product = jdbc.queryForObject("INSERT INTO product(title) VALUES (?) RETURNING id",
                Integer.class, "route-" + UUID.randomUUID());
        int op = jdbc.queryForObject("INSERT INTO product_operation(product_id, sequence_number, operation_type, "
                + "time_per_unit_seconds) VALUES (?, 1, 'TURNING', 600) RETURNING id", Integer.class, product);
        jdbc.update("INSERT INTO product_operation(product_id, sequence_number, operation_type, time_per_unit_seconds) "
                + "VALUES (?, 2, 'GRINDING', 300)", product);
        int order = jdbc.queryForObject("INSERT INTO production_order(due_date, comment) VALUES ('2099-01-01', 'Snapshot test') RETURNING id", Integer.class);
        jdbc.update("INSERT INTO production_order_item(order_id, product_id, quantity) VALUES (?, ?, 20)", order, product);
        return new Fixture(order, product, op);
    }
    private void asPlanner() {
        var user = new AppPrincipal(1L, "Planner", "planner", "unused", Role.PLANNER, true);
        SecurityContextHolder.getContext().setAuthentication(
                UsernamePasswordAuthenticationToken.authenticated(user, null, user.getAuthorities()));
    }
    record Fixture(int order, int product, Integer operation) {}
}
