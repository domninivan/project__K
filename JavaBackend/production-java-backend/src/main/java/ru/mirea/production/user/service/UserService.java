package ru.mirea.production.user.service;

import java.nio.charset.StandardCharsets;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.http.HttpStatus;
import ru.mirea.production.user.dao.AppUserDAO;
import ru.mirea.production.user.dto.*;
import ru.mirea.production.user.entity.*;

@Service
public class UserService {
    private final AppUserDAO users;
    private final PasswordEncoder passwords;

    public UserService(AppUserDAO users, PasswordEncoder passwords) {
        this.users = users;
        this.passwords = passwords;
    }

    @Transactional
    public UserResponse create(CreateUserRequest request) {
        if (request.password().getBytes(StandardCharsets.UTF_8).length > 72) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Пароль превышает 72 байта UTF-8.");
        }
        AppUser user = users.save(new AppUser(request.fullName().strip(), request.username(),
                passwords.encode(request.password()), request.role()));
        return new UserResponse(user.getId(), user.getFullName(), user.getUsername(), user.getRole(), user.isActive());
    }
}
