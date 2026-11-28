package ru.mirea.production.planning;

import java.time.*;
import java.util.List;
import org.junit.jupiter.api.Test;
import com.fasterxml.jackson.databind.SerializationFeature;
import org.springframework.http.converter.json.Jackson2ObjectMapperBuilder;
import org.springframework.test.web.client.MockRestServiceServer;
import ru.mirea.production.planning.client.GoRestClientFactory;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.*;
import static org.springframework.test.web.client.response.MockRestResponseCreators.*;
import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.order.entity.enums.PriorityEnum;
import static org.junit.jupiter.api.Assertions.*;

class PlanningJsonTest {
    @Test void datesMustBeIsoStringsForGo() throws Exception {
        // Same date setting as Spring Boot's jacksonObjectMapper injected by XML.
        var mapper = Jackson2ObjectMapperBuilder.json()
                .featuresToDisable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS).build();
        var builder = GoRestClientFactory.builder("http://go.test", mapper);
        var server = MockRestServiceServer.bindTo(builder).build();
        var request = new PlanningRequest(5, OffsetDateTime.parse("2026-11-24T09:00:00+03:00"),
                LocalDate.of(2026, 12, 20), PriorityEnum.MEDIUM, List.of());
        var body = mapper.readTree(mapper.writeValueAsBytes(request));
        assertEquals("2026-12-20", body.get("dueDate").asText());
        assertEquals("2026-11-24T09:00:00+03:00", body.get("earliestStartAt").asText());
        server.expect(requestTo("http://go.test/api/plans"))
                .andExpect(content().json("{\"dueDate\":\"2026-12-20\",\"earliestStartAt\":\"2026-11-24T09:00:00+03:00\"}"))
                .andRespond(withSuccess());
        builder.build().post().uri("/api/plans").body(request).retrieve().toBodilessEntity();
        server.verify();
    }
}
