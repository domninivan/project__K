package ru.mirea.production.planning;

import com.fasterxml.jackson.databind.*;
import java.time.*;
import java.util.*;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.condition.EnabledIfEnvironmentVariable;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.web.servlet.*;
import static org.junit.jupiter.api.Assertions.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

/** Requires isolated migrated Go DB + real Go server and an isolated Java DB. */
@SpringBootTest
@AutoConfigureMockMvc
@EnabledIfEnvironmentVariable(named = "JAVA_GO_TEST_URL", matches = ".+")
class JavaGoIntegrationTest {
    @Autowired MockMvc mvc;
    @Autowired ObjectMapper json;
    @DynamicPropertySource static void properties(DynamicPropertyRegistry r) {
        String testDatabase = Objects.requireNonNull(System.getenv("JAVA_MIGRATION_TEST_URL"),
                "Set JAVA_MIGRATION_TEST_URL to an isolated Java test database");
        r.add("DB_URL", () -> testDatabase);
        r.add("DB_USERNAME", () -> System.getenv().getOrDefault("JAVA_MIGRATION_TEST_USER", "postgres"));
        r.add("DB_PASSWORD", () -> System.getenv().getOrDefault("JAVA_MIGRATION_TEST_PASSWORD", ""));
        r.add("go.base-url", () -> System.getenv("JAVA_GO_TEST_URL"));
        r.add("JWT_SECRET_BASE64", () -> Base64.getEncoder().encodeToString(new byte[32]));
    }

    @Test void allExistingPlanningMethodsCommunicateWithRealGo() throws Exception {
        String planner = user("PLANNER");
        String client = user("CLIENT");
        var centers = read(mvc.perform(get("/api/planner/work-centers").header("Authorization", planner))
                .andExpect(status().isOk()));
        assertFalse(centers.isEmpty());
        int center = centers.get(0).get("id").asInt();
        var date = LocalDate.now(ZoneId.of("Europe/Moscow"));
        mvc.perform(get("/api/planner/work-centers/" + center + "/availability").header("Authorization", planner)
                        .param("from", date.toString()).param("to", date.plusDays(14).toString()))
                .andExpect(status().isOk()).andExpect(jsonPath("$.intervals[0].startAt").exists());
        var products = read(mvc.perform(get("/api/products").header("Authorization", client)).andExpect(status().isOk()));
        int order = order(client, products.get(0).get("id").asInt(), date.plusDays(30));
        var plan = read(mvc.perform(post("/api/planner/orders/" + order).header("Authorization", planner))
                .andExpect(status().isOk()).andExpect(jsonPath("$.status").value("PROPOSED"))
                .andExpect(jsonPath("$.dueDate").value(date.plusDays(30).toString()))
                .andExpect(jsonPath("$.items[0].operations[0].operationId").isNumber()));
        int planId = plan.get("planId").asInt();
        mvc.perform(get("/api/planner/orders/" + order).header("Authorization", planner))
                .andExpect(status().isOk()).andExpect(jsonPath("$.planId").value(planId));
        // Explicit +03:00 catches incorrect '+' encoding in query parameters.
        mvc.perform(get("/api/planner/schedule").header("Authorization", planner)
                        .param("from", date + "T00:00:00+03:00")
                        .param("to", date.plusDays(30) + "T00:00:00+03:00"))
                .andExpect(status().isOk()).andExpect(jsonPath("$.workCenters[0].operations[0]").exists());
        mvc.perform(post("/api/planner/plans/" + planId + "/accept").header("Authorization", planner))
                .andExpect(status().isOk()).andExpect(jsonPath("$.status").value("ACCEPTED"));
        mvc.perform(post("/api/planner/plans/" + planId + "/reject").header("Authorization", planner))
                .andExpect(status().isConflict());
        mvc.perform(post("/api/planner/orders/" + order + "/cancel-plan").header("Authorization", planner))
                .andExpect(status().isOk()).andExpect(jsonPath("$.status").value("CANCELLED"));
        int second = order(client, products.get(0).get("id").asInt(), date.plusDays(30));
        int draft = read(mvc.perform(post("/api/planner/orders/" + second).header("Authorization", planner))
                .andExpect(status().isOk())).get("planId").asInt();
        mvc.perform(post("/api/planner/plans/" + draft + "/reject").header("Authorization", planner))
                .andExpect(status().isOk()).andExpect(jsonPath("$.status").value("REJECTED"));
        mvc.perform(get("/api/planner/orders/2147480000").header("Authorization", planner))
                .andExpect(status().isNotFound());
        mvc.perform(get("/api/planner/work-centers/" + center + "/availability").header("Authorization", planner)
                        .param("from", date.plusDays(1).toString()).param("to", date.toString()))
                .andExpect(status().isBadRequest());
    }

    private String user(String role) throws Exception {
        String username = role + "-" + UUID.randomUUID();
        mvc.perform(post("/api/users").contentType("application/json").content(json.writeValueAsString(Map.of(
                "fullName", "Integration User", "username", username, "password", "test-password", "role", role))))
                .andExpect(status().isCreated());
        return "Bearer " + read(mvc.perform(post("/api/auth/login").contentType("application/json")
                .content(json.writeValueAsString(Map.of("username", username, "password", "test-password"))))
                .andExpect(status().isOk())).get("accessToken").asText();
    }
    private int order(String token, int product, LocalDate due) throws Exception {
        return read(mvc.perform(post("/api/orders").header("Authorization", token).contentType("application/json")
                .content(json.writeValueAsString(Map.of("dueDate", due.toString(), "comment", "Go integration",
                        "items", List.of(Map.of("productId", product, "quantity", 1))))))
                .andExpect(status().isOk())).get("id").asInt();
    }
    private JsonNode read(ResultActions action) throws Exception {
        return json.readTree(action.andReturn().getResponse().getContentAsString());
    }
}
