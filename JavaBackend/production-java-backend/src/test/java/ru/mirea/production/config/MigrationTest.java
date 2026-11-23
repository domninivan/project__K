package ru.mirea.production.config;

import org.flywaydb.core.Flyway;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.condition.EnabledIfEnvironmentVariable;
import ru.mirea.production.user.entity.Role;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;

@EnabledIfEnvironmentVariable(named = "JAVA_MIGRATION_TEST_URL", matches = ".+")
class MigrationTest {
    @Test
    void emptyDatabaseMigratesAndEnforcesRoles() throws Exception {
        inSchema((connection, schema) -> {
            Flyway flyway = flyway(schema, false);
            assertEquals(5, flyway.migrate().migrationsExecuted);
            assertEquals(0, flyway.migrate().migrationsExecuted);
            try (var statement = connection.createStatement()) {
                try (var result = statement.executeQuery("SELECT count(*) FROM product")) {
                    assertTrue(result.next());
                    assertEquals(4, result.getInt(1));
                }
                try (var result = statement.executeQuery("SELECT count(*) FROM production_order")) {
                    assertTrue(result.next());
                    assertEquals(2, result.getInt(1));
                }
                for (Role role : Role.values()) {
                    statement.executeUpdate("INSERT INTO app_user(full_name, username, password_hash, role) "
                            + "VALUES ('Test', '" + role.name() + "', 'test-only-placeholder', '" + role.name() + "')");
                }
                try (var result = statement.executeQuery("SELECT count(*) FROM app_user WHERE active")) {
                    assertTrue(result.next());
                    assertEquals(5, result.getInt(1));
                }
                SQLException invalidRole = assertThrows(SQLException.class, () -> statement.executeUpdate(
                        "INSERT INTO app_user(full_name, username, password_hash, role) "
                                + "VALUES ('Test', 'invalid', 'test-only-placeholder', 'ADMIN')"));
                assertEquals("23514", invalidRole.getSQLState());
                SQLException duplicate = assertThrows(SQLException.class, () -> statement.executeUpdate(
                        "INSERT INTO app_user(full_name, username, password_hash, role) "
                                + "VALUES ('Test', 'CLIENT', 'test-only-placeholder', 'CLIENT')"));
                assertEquals("23505", duplicate.getSQLState());
            }
        });
    }

    @Test
    void baselinePreservesExistingOrders() throws Exception {
        inSchema((connection, schema) -> {
            Flyway.configure().dataSource(url(), user(), password()).schemas(schema).defaultSchema(schema)
                    .locations("classpath:db/migration").target("1").load().migrate();
            try (var statement = connection.createStatement()) {
                // Simulate the old Docker-created schema without Flyway history.
                statement.execute("DROP TABLE flyway_schema_history");
                statement.executeUpdate("INSERT INTO production_order(due_date, comment) "
                        + "VALUES ('2026-11-23', 'preserve-me')");
            }
            assertThrows(org.flywaydb.core.api.FlywayException.class,
                    () -> flyway(schema, false).migrate());
            assertEquals(4, flyway(schema, true).migrate().migrationsExecuted);
            try (var statement = connection.createStatement();
                 var result = statement.executeQuery("SELECT comment, client_id FROM production_order WHERE comment = 'preserve-me'")) {
                assertTrue(result.next());
                assertEquals("preserve-me", result.getString(1));
                assertNull(result.getObject(2));
            }
        });
    }

    private Flyway flyway(String schema, boolean baseline) {
        return Flyway.configure().dataSource(url(), user(), password())
                .schemas(schema).defaultSchema(schema).locations("classpath:db/migration")
                .baselineVersion("1").baselineOnMigrate(baseline).cleanDisabled(true).load();
    }

    @Test
    void upgradingExistingSeedDoesNotDuplicateOrOverwriteData() throws Exception {
        inSchema((connection, schema) -> {
            Flyway.configure().dataSource(url(), user(), password()).schemas(schema).defaultSchema(schema)
                    .locations("classpath:db/migration").target("3").load().migrate();
            // Existing installations already have the former Docker seed.
            try (var input = getClass().getResourceAsStream("/db/migration/V4__seed_demo_data.sql");
                 var statement = connection.createStatement()) {
                assertNotNull(input);
                statement.execute(new String(input.readAllBytes(), java.nio.charset.StandardCharsets.UTF_8));
                statement.executeUpdate("UPDATE product SET description = 'keep-existing-description'");
            }
            assertEquals(2, flyway(schema, false).migrate().migrationsExecuted);
            assertEquals(0, flyway(schema, false).migrate().migrationsExecuted);
            String[] tables = {"product", "product_operation", "product_material_requirement",
                    "production_order", "production_order_item"};
            int[] counts = {4, 16, 8, 2, 4};
            try (var statement = connection.createStatement()) {
                for (int i = 0; i < tables.length; i++) {
                    try (var result = statement.executeQuery("SELECT count(*) FROM " + tables[i])) {
                        assertTrue(result.next());
                        assertEquals(counts[i], result.getInt(1), tables[i]);
                    }
                }
                try (var result = statement.executeQuery(
                        "SELECT count(*) FROM product WHERE description = 'keep-existing-description'")) {
                    assertTrue(result.next());
                    assertEquals(4, result.getInt(1));
                }
            }
        });
    }

    private void inSchema(SchemaCheck check) throws Exception {
        String schema = "migration_test_" + UUID.randomUUID().toString().replace("-", "");
        try (Connection connection = DriverManager.getConnection(url(), user(), password());
             var statement = connection.createStatement()) {
            statement.execute("CREATE SCHEMA " + schema);
            try {
                statement.execute("SET search_path TO " + schema);
                check.run(connection, schema);
            } finally {
                statement.execute("DROP SCHEMA " + schema + " CASCADE");
            }
        }
    }

    private String url() { return System.getenv("JAVA_MIGRATION_TEST_URL"); }
    private String user() { return System.getenv().getOrDefault("JAVA_MIGRATION_TEST_USER", "postgres"); }
    private String password() { return System.getenv().getOrDefault("JAVA_MIGRATION_TEST_PASSWORD", ""); }

    @FunctionalInterface
    private interface SchemaCheck {
        void run(Connection connection, String schema) throws Exception;
    }
}
