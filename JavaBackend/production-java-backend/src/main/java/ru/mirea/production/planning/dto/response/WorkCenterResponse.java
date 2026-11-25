package ru.mirea.production.planning.dto.response;

import java.util.List;

public record WorkCenterResponse(
        Integer id,
        String code,
        String title,
        Boolean active,
        List<String> operationTypes
) {}
