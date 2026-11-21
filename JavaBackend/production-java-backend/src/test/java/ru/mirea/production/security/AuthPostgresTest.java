package ru.mirea.production.security;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.boot.test.system.CapturedOutput;
import org.springframework.boot.test.system.OutputCaptureExtension;
import org.junit.jupiter.api.condition.EnabledIfEnvironmentVariable;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.web.servlet.MockMvc;
import com.fasterxml.jackson.databind.ObjectMapper;
import javax.sql.DataSource;
import java.util.UUID;
import ru.mirea.production.user.entity.Role;
import static org.junit.jupiter.api.Assertions.*;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
@EnabledIfEnvironmentVariable(named = "JAVA_MIGRATION_TEST_URL", matches = ".+")
class AuthPostgresTest {
    @Autowired MockMvc mvc;
    @Autowired DataSource dataSource;
    @Autowired PasswordEncoder encoder;
    @Autowired ObjectMapper json;

    @Test
    @ExtendWith(OutputCaptureExtension.class)
    void invalidLoginDoesNotLogPassword(CapturedOutput output) throws Exception {
        String password = "do-not-log-this-password-".repeat(5);
        mvc.perform(post("/api/auth/login").contentType("application/json")
                        .content(json.writeValueAsString(java.util.Map.of("username", "test", "password", password))))
                .andExpect(status().isBadRequest());
        assertFalse(output.getAll().contains(password), "Login validation leaked password to logs");
    }

    @DynamicPropertySource
    static void properties(DynamicPropertyRegistry registry) {
        registry.add("DB_URL", () -> System.getenv("JAVA_MIGRATION_TEST_URL"));
        registry.add("DB_USERNAME", () -> System.getenv().getOrDefault("JAVA_MIGRATION_TEST_USER", "postgres"));
        registry.add("DB_PASSWORD", () -> System.getenv().getOrDefault("JAVA_MIGRATION_TEST_PASSWORD", ""));
        registry.add("JWT_SECRET_BASE64", () -> java.util.Base64.getEncoder().encodeToString(new byte[32]));
    }

    @Test
    void realLoginCreationQueriesAndOwnerIsolation() throws Exception {
        var jdbc = new JdbcTemplate(dataSource);
        String suffix = UUID.randomUUID().toString();
        String first = "first-" + suffix;
        String second = "second-" + suffix;
        String planner = "planner-" + suffix;
        register(planner, Role.PLANNER);
        String plannerToken = login(planner);
        long firstClientId = register(first, Role.CLIENT);
        register(second, Role.CLIENT);
        Integer product = jdbc.queryForObject("INSERT INTO product(title) VALUES (?) RETURNING id",
                Integer.class, "Test " + suffix);
        String firstToken = login(first);
        String secondToken = login(second);
        mvc.perform(get("/api/me").header("Authorization", firstToken))
                .andExpect(status().isOk()).andExpect(jsonPath("$.role").value("CLIENT"));
        var created = mvc.perform(post("/api/orders").header("Authorization", firstToken)
                        .contentType("application/json").content("{\"dueDate\":\"2099-01-01\",\"comment\":\"Test\","
                                + "\"clientId\":999999,\"items\":[{\"productId\":" + product + ",\"quantity\":2}]}"))
                .andExpect(status().isOk()).andReturn();
        int order = json.readTree(created.getResponse().getContentAsString()).get("id").asInt();
        assertEquals(firstClientId, jdbc.queryForObject(
                "SELECT client_id FROM production_order WHERE id = ?", Long.class, order));
        mvc.perform(get("/api/orders/" + order).header("Authorization", firstToken))
                .andExpect(status().isOk()).andExpect(jsonPath("$.id").value(order))
                .andExpect(jsonPath("$.status").value("NEW"))
                .andExpect(jsonPath("$.items[0].productId").value(product))
                .andExpect(jsonPath("$.items[0].quantity").value(2));
        mvc.perform(get("/api/orders/" + order).header("Authorization", secondToken)).andExpect(status().isNotFound());
        mvc.perform(delete("/api/orders/" + order).header("Authorization", secondToken)).andExpect(status().isNotFound());
        mvc.perform(get("/api/orders").header("Authorization", secondToken))
                .andExpect(status().isOk()).andExpect(content().json("[]"));
        mvc.perform(get("/api/orders").header("Authorization", firstToken))
                .andExpect(status().isOk()).andExpect(jsonPath("$[0].id").value(order));
        mvc.perform(get("/api/orders/" + order)).andExpect(status().isUnauthorized());
        mvc.perform(get("/api/orders/" + order).header("Authorization", plannerToken)).andExpect(status().isOk());
        mvc.perform(post("/api/orders").header("Authorization", plannerToken).contentType("application/json").content("{}"))
                .andExpect(status().isForbidden());
        jdbc.update("UPDATE app_user SET active = false WHERE username = ?", first);
        mvc.perform(get("/api/me").header("Authorization", firstToken)).andExpect(status().isUnauthorized());
    }

    @Test
    void anonymousRegistrationAcceptsEveryRoleAndHashesPasswords() throws Exception {
        var jdbc = new JdbcTemplate(dataSource);
        String suffix = UUID.randomUUID().toString();
        for (Role role : Role.values()) {
            String username = role.name() + "-" + suffix;
            String body = json.writeValueAsString(java.util.Map.of("fullName", "New User", "username", username,
                    "password", "test-password", "role", role.name()));
            mvc.perform(post("/api/users").contentType("application/json").content(body))
                    .andExpect(status().isCreated()).andExpect(jsonPath("$.username").value(username))
                    .andExpect(jsonPath("$.role").value(role.name())).andExpect(jsonPath("$.active").value(true))
                    .andExpect(jsonPath("$.password").doesNotExist()).andExpect(jsonPath("$.passwordHash").doesNotExist());
            String hash = jdbc.queryForObject("SELECT password_hash FROM app_user WHERE username = ?", String.class, username);
            assertNotEquals("test-password", hash);
            assertTrue(encoder.matches("test-password", hash));
            String newToken = login(username);
            mvc.perform(get("/api/me").header("Authorization", newToken))
                    .andExpect(status().isOk()).andExpect(jsonPath("$.role").value(role.name()));
            mvc.perform(post("/api/users").contentType("application/json").content(body))
                    .andExpect(status().isConflict()).andExpect(jsonPath("$.code").value("USERNAME_TAKEN"));
            if (role != Role.PLANNER) {
                mvc.perform(post("/api/users").header("Authorization", newToken).contentType("application/json").content(body))
                        .andExpect(status().isConflict());
            }
        }
        mvc.perform(post("/api/users").contentType("application/json").content("{}"))
                .andExpect(status().isBadRequest());
        for (String body : new String[]{"{}",
                "{\"fullName\":\"User\",\"username\":\"x\",\"password\":\"test-password\",\"role\":\"ADMIN\"}",
                "{\"fullName\":\"User\",\"username\":\"x\",\"password\":\"short\",\"role\":\"CLIENT\"}",
                json.writeValueAsString(java.util.Map.of("fullName", "User", "username", "x",
                        "password", "я".repeat(40), "role", "CLIENT"))}) {
            mvc.perform(post("/api/users").contentType("application/json").content(body))
                    .andExpect(status().isBadRequest()).andExpect(jsonPath("$.code").value("VALIDATION_ERROR"));
        }
    }

    private long register(String username, Role role) throws Exception {
        var result = mvc.perform(post("/api/users")
                        .contentType("application/json").content(json.writeValueAsString(java.util.Map.of(
                                "fullName", "Test Client", "username", username,
                                "password", "test-password", "role", role.name()))))
                .andExpect(status().isCreated()).andExpect(jsonPath("$.role").value(role.name()))
                .andReturn();
        return json.readTree(result.getResponse().getContentAsString()).get("id").asLong();
    }

    private String login(String username) throws Exception {
        var result = mvc.perform(post("/api/auth/login").contentType("application/json")
                        .content(json.writeValueAsString(java.util.Map.of("username", username, "password", "test-password"))))
                .andExpect(status().isOk()).andReturn();
        return "Bearer " + json.readTree(result.getResponse().getContentAsString()).get("accessToken").asText();
    }
}
