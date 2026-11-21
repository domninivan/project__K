package ru.mirea.production.security;

import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.core.context.SecurityContextHolder;
import ru.mirea.production.user.entity.Role;
import java.util.Arrays;

public final class CurrentActor {
    private CurrentActor() {}

    public static AppPrincipal require(Role... roles) {
        var auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth == null || !auth.isAuthenticated() || !(auth.getPrincipal() instanceof AppPrincipal user)
                || !user.isEnabled() || Arrays.stream(roles).noneMatch(role -> role == user.getRole())) {
            throw new AccessDeniedException("Действие запрещено");
        }
        return user;
    }
}
