package ru.mirea.production.planning.client;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.RestClient;

public final class GoRestClientFactory {
    private GoRestClientFactory() {}

    public static RestClient.Builder builder(String baseUrl, ObjectMapper mapper) {
        return RestClient.builder().baseUrl(baseUrl).messageConverters(converters -> {
            converters.removeIf(MappingJackson2HttpMessageConverter.class::isInstance);
            converters.add(new MappingJackson2HttpMessageConverter(mapper));
        });
    }
}
