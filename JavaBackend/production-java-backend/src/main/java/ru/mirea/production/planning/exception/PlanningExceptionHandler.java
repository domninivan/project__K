package ru.mirea.production.planning.exception;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.util.Map;

@RestControllerAdvice
public class PlanningExceptionHandler {

    @ExceptionHandler(GoServiceUnavailableException.class)
    public ResponseEntity<Map<String, Object>> handleGoServiceUnavailable(
            GoServiceUnavailableException ex
    ) {
        return ResponseEntity
                .status(HttpStatus.SERVICE_UNAVAILABLE)
                .body(Map.of(
                        "status", HttpStatus.SERVICE_UNAVAILABLE.value(),
                        "error", "Service Unavailable",
                        "message", ex.getMessage()
                ));
    }

    @ExceptionHandler(GoServiceException.class)
    public ResponseEntity<Map<String, Object>> handleGoServiceException(
            GoServiceException ex
    ) {
        return ResponseEntity
                .status(HttpStatus.BAD_GATEWAY)
                .body(Map.of(
                        "status", HttpStatus.BAD_GATEWAY.value(),
                        "error", "Bad Gateway",
                        "message", ex.getMessage()
                ));
    }
    @ExceptionHandler(ProductionPlanNotFoundException.class)
    public ResponseEntity<Map<String, Object>> handleProductionPlanNotFound(
            ProductionPlanNotFoundException ex
    ) {
        return ResponseEntity
                .status(HttpStatus.NOT_FOUND)
                .body(Map.of(
                        "status", HttpStatus.NOT_FOUND.value(),
                        "error", "Not Found",
                        "message", ex.getMessage()
                ));
    }

    @ExceptionHandler(GoBadRequestException.class)
    public ResponseEntity<Map<String, Object>> handleBadRequest(
            GoBadRequestException ex
    ) {
        return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .body(Map.of(
                        "status", HttpStatus.BAD_REQUEST.value(),
                        "error", "Bad Request",
                        "message", ex.getMessage()
                ));
    }

    @ExceptionHandler(GoResourceNotFoundException.class)
    public ResponseEntity<Map<String, Object>> handleNotFound(
            GoResourceNotFoundException ex
    ) {
        return ResponseEntity
                .status(HttpStatus.NOT_FOUND)
                .body(Map.of(
                        "status", HttpStatus.NOT_FOUND.value(),
                        "error", "Not Found",
                        "message", ex.getMessage()
                ));
    }

    @ExceptionHandler(GoConflictException.class)
    public ResponseEntity<Map<String, Object>> handleConflict(
            GoConflictException ex
    ) {
        return ResponseEntity
                .status(HttpStatus.CONFLICT)
                .body(Map.of(
                        "status", HttpStatus.CONFLICT.value(),
                        "error", "Conflict",
                        "message", ex.getMessage()
                ));
    }
}