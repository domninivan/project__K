package ru.mirea.production.security;

import java.util.List;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.User;
import ru.mirea.production.user.entity.Role;

public class AppPrincipal extends User {
    private final Long id;
    private final String fullName;
    private final Role role;

    public AppPrincipal(Long id, String fullName, String username, String hash, Role role, boolean active) {
        super(username, hash, active, true, true, true,
                List.of(new SimpleGrantedAuthority("ROLE_" + role.name())));
        this.id = id;
        this.fullName = fullName;
        this.role = role;
    }

    public Long getId() { return id; }
    public String getFullName() { return fullName; }
    public Role getRole() { return role; }
}
