package ru.mirea.production.security;

import javax.sql.DataSource;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import ru.mirea.production.user.entity.Role;

public class AccountDetailsService implements UserDetailsService {
    private final JdbcTemplate jdbc;

    public AccountDetailsService(DataSource dataSource) { this.jdbc = new JdbcTemplate(dataSource); }

    @Override
    public AppPrincipal loadUserByUsername(String username) {
        return find("username", username);
    }

    public AppPrincipal loadById(Long id) { return find("id", id); }

    private AppPrincipal find(String column, Object value) {
        var users = jdbc.query("SELECT id, full_name, username, password_hash, role, active FROM app_user WHERE "
                        + column + " = ?",
                (rs, row) -> new AppPrincipal(rs.getLong("id"), rs.getString("full_name"),
                        rs.getString("username"), rs.getString("password_hash"),
                        Role.valueOf(rs.getString("role")), rs.getBoolean("active")), value);
        if (users.isEmpty()) throw new UsernameNotFoundException("Учётная запись не найдена");
        return users.get(0);
    }
}
