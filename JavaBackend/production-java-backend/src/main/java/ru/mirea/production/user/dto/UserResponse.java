package ru.mirea.production.user.dto;

import ru.mirea.production.user.entity.Role;

public record UserResponse(Long id, String fullName, String username, Role role, boolean active) {}
