package ru.mirea.production.security;

import java.util.Arrays;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

/** Local console helper; no HTTP endpoint and no database changes. */
public final class PasswordHash {
    private PasswordHash() {}
    public static void main(String[] args) {
        var console = System.console();
        if (console == null) throw new IllegalStateException("Run from an interactive terminal");
        char[] password = console.readPassword("Password: ");
        if (password == null) throw new IllegalArgumentException("Password is required");
        try {
            if (password.length == 0) throw new IllegalArgumentException("Password is required");
            System.out.println(new BCryptPasswordEncoder().encode(new String(password)));
        } finally {
            Arrays.fill(password, '\0');
        }
    }
}
