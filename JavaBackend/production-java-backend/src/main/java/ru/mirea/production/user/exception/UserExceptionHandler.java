package ru.mirea.production.user.exception;

import java.util.LinkedHashMap;
import java.util.Map;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;
import ru.mirea.production.user.controller.UserController;

@RestControllerAdvice(assignableTypes = UserController.class)
public class UserExceptionHandler {
    @ExceptionHandler(DuplicateUsernameException.class)
    public ResponseEntity<?> duplicate(DuplicateUsernameException ex) {
        return error(409, "USERNAME_TAKEN", ex.getMessage(), Map.of());
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<?> validation(MethodArgumentNotValidException ex) {
        Map<String, String> fields = new LinkedHashMap<>();
        ex.getBindingResult().getFieldErrors().forEach(field ->
                fields.putIfAbsent(field.getField(), field.getDefaultMessage()));
        return error(400, "VALIDATION_ERROR", "Проверьте поля запроса.", fields);
    }

    @ExceptionHandler(HttpMessageNotReadableException.class)
    public ResponseEntity<?> invalidJson(HttpMessageNotReadableException ex) {
        return error(400, "VALIDATION_ERROR", "Некорректный JSON или неизвестная роль.", Map.of());
    }

    @ExceptionHandler(ResponseStatusException.class)
    public ResponseEntity<?> invalidPassword(ResponseStatusException ex) {
        return error(ex.getStatusCode().value(), "VALIDATION_ERROR", ex.getReason(), Map.of());
    }

    private ResponseEntity<?> error(int status, String code, String message, Map<String, String> fields) {
        return ResponseEntity.status(status).body(Map.of("code", code, "message", message,
                "fieldErrors", fields, "retryable", false));
    }
}
