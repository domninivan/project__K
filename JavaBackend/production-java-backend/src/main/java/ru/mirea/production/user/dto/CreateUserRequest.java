package ru.mirea.production.user.dto;

import jakarta.validation.constraints.*;
import ru.mirea.production.user.entity.Role;

public record CreateUserRequest(
        @NotBlank @Size(max = 200) String fullName,
        @NotBlank @Size(max = 100) @Pattern(regexp = "\\S+") String username,
        @NotBlank @Size(min = 8, max = 72) String password,
        @NotNull Role role) {
    @Override
    public String toString() {
        return "CreateUserRequest[credentials redacted]";
    }
}
