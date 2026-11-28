package ru.mirea.production.security;

import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import org.springframework.security.authentication.*;
import org.springframework.security.core.Authentication;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.security.core.AuthenticationException;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.http.HttpStatus;
import java.nio.charset.StandardCharsets;
import java.util.Map;
import java.util.LinkedHashMap;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.beans.factory.annotation.Qualifier;
import ru.mirea.production.user.entity.Role;

@RestController
public class AuthController {
    private final AuthenticationManager manager;
    private final TokenService tokens;

    public AuthController(@Qualifier("authenticationManager") AuthenticationManager manager, TokenService tokens) {
        this.manager = manager;
        this.tokens = tokens;
    }

    @PostMapping("/api/auth/login")
    public ResponseEntity<TokenService.TokenResponse> login(@Valid @RequestBody LoginRequest request) {
        if (request.password().getBytes(StandardCharsets.UTF_8).length > 72) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Пароль превышает 72 байта UTF-8");
        }
        var result = manager.authenticate(UsernamePasswordAuthenticationToken.unauthenticated(
                request.username(), request.password()));
        return ResponseEntity.ok().header("Cache-Control", "no-store").header("Pragma", "no-cache")
                .body(tokens.issue((AppPrincipal) result.getPrincipal()));
    }

    @GetMapping("/api/me")
    public Me me(Authentication authentication) {
        AppPrincipal user = (AppPrincipal) authentication.getPrincipal();
        return new Me(user.getId(), user.getFullName(), user.getRole());
    }

    public record LoginRequest(@NotBlank @Size(max = 100) String username,
                               @NotBlank @Size(max = 72) String password) {
        @Override public String toString() { return "LoginRequest[credentials redacted]"; }
    }
    public record Me(Long id, String fullName, Role role) {}

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Map<String, Object>> invalidLoginRequest(MethodArgumentNotValidException ex) {
        // Do not pass this exception to the default resolver: its message includes rejected passwords.
        Map<String, String> fields = new LinkedHashMap<>();
        ex.getBindingResult().getFieldErrors().forEach(error ->
                fields.putIfAbsent(error.getField(), error.getDefaultMessage()));
        return ResponseEntity.badRequest().body(Map.of("code", "VALIDATION_ERROR",
                "message", "Проверьте поля запроса.", "fieldErrors", fields, "retryable", false));
    }

    @ExceptionHandler(AuthenticationException.class)
    public ResponseEntity<Map<String, Object>> authenticationFailed(AuthenticationException ex) {
        return ResponseEntity.status(401).header("WWW-Authenticate", "Bearer")
                .header("Cache-Control", "no-store")
                .body(Map.of("code", "UNAUTHORIZED", "message", "Неверные учётные данные.",
                        "fieldErrors", Map.of(), "retryable", false));
    }
}
