package ru.mirea.production.security;

import java.time.LocalDate;
import java.util.*;
import org.junit.jupiter.api.*;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.test.util.ReflectionTestUtils;
import ru.mirea.production.order.dao.ProductionOrderDAO;
import ru.mirea.production.order.dto.ProductionOrderRequestDTO;
import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.enums.StatusEnum;
import ru.mirea.production.order.exception.OrderNotFoundException;
import ru.mirea.production.order.service.ProductionOrderServiceImpl;
import ru.mirea.production.user.entity.Role;

import static org.junit.jupiter.api.Assertions.*;

class OrderAccessTest {
    private final Orders dao = new Orders();
    private final ProductionOrderServiceImpl service = new ProductionOrderServiceImpl();

    @BeforeEach
    void setup() { ReflectionTestUtils.setField(service, "productionOrderDAO", dao); }
    @AfterEach
    void cleanup() { SecurityContextHolder.clearContext(); }

    @Test
    void clientOnlyUsesQueriesFilteredByOwner() {
        as(1L, Role.CLIENT);
        assertEquals(List.of(dao.orders.get(1)), service.getAllOrders());
        assertEquals(1L, dao.lastClient);
        assertEquals(0, dao.unfilteredReads);
        assertSame(dao.orders.get(1), service.getOrderById(1));
        assertThrows(OrderNotFoundException.class, () -> service.getOrderById(2));
        assertThrows(OrderNotFoundException.class, () -> service.getOrderById(3));
        assertThrows(OrderNotFoundException.class, () -> service.cancelOrder(2));
        assertEquals(StatusEnum.NEW, dao.orders.get(2).getStatus());
        assertEquals(0, dao.saves);
    }

    @Test
    void plannerCanReadAllIncludingUnassignedButCannotCreateOrCancel() {
        as(8L, Role.PLANNER);
        assertEquals(3, service.getAllOrders().size());
        assertSame(dao.orders.get(3), service.getOrderById(3));
        assertThrows(AccessDeniedException.class, () -> service.cancelOrder(1));
        assertThrows(AccessDeniedException.class, () -> service.createOrder(request()));
    }

    @Test
    void newOrderOwnerComesFromPrincipalAndClientCanCancelOwnOrder() {
        as(42L, Role.CLIENT);
        assertEquals(42L, service.createOrder(request()).getClientId());
        as(1L, Role.CLIENT);
        assertEquals(StatusEnum.CANCELLED, service.cancelOrder(1).getStatus());
    }

    @Test
    void anonymousAndOtherRolesCannotReadOrders() {
        assertThrows(AccessDeniedException.class, service::getAllOrders);
        for (Role role : List.of(Role.WORKER, Role.SHOP_MANAGER, Role.REPORTING)) {
            as(5L, role);
            assertThrows(AccessDeniedException.class, service::getAllOrders);
            assertThrows(AccessDeniedException.class, () -> service.getOrderById(1));
        }
    }

    private ProductionOrderRequestDTO request() {
        // Service ownership test; HTTP bean validation covers nonempty items separately.
        return new ProductionOrderRequestDTO(LocalDate.now(), "Test", List.of());
    }

    private void as(Long id, Role role) {
        var user = new AppPrincipal(id, "Test", "test", "unused", role, true);
        SecurityContextHolder.getContext().setAuthentication(
                UsernamePasswordAuthenticationToken.authenticated(user, null, user.getAuthorities()));
    }

    static class Orders implements ProductionOrderDAO {
        final Map<Integer, ProductionOrder> orders = new LinkedHashMap<>();
        Long lastClient;
        int unfilteredReads;
        int saves;
        Orders() {
            for (int id = 1; id <= 3; id++) {
                var order = new ProductionOrder(LocalDate.now(), "Test");
                order.setClientId(id == 3 ? null : (long) id);
                orders.put(id, order);
            }
        }
        @Override public List<ProductionOrder> getAllOrders() {
            unfilteredReads++;
            return new ArrayList<>(orders.values());
        }
        @Override public ProductionOrder getOrderById(int id) { unfilteredReads++; return orders.get(id); }
        @Override public List<ProductionOrder> getOrdersByClientId(Long clientId) {
            lastClient = clientId;
            return orders.values().stream().filter(order -> clientId.equals(order.getClientId())).toList();
        }
        @Override public ProductionOrder getOrderByIdAndClientId(int id, Long clientId) {
            lastClient = clientId;
            var order = orders.get(id);
            return order != null && clientId.equals(order.getClientId()) ? order : null;
        }
        @Override public ProductionOrder saveOrder(ProductionOrder order) { saves++; return order; }
    }
}
