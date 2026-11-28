package ru.mirea.production.user.controller;

import jakarta.validation.Valid;
import java.net.URI;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import ru.mirea.production.user.dto.*;
import ru.mirea.production.user.service.UserService;

@RestController
@RequestMapping("/api/users")
public class UserController {
    private final UserService users;
    public UserController(UserService users) { this.users = users; }

    @PostMapping
    public ResponseEntity<UserResponse> create(@Valid @RequestBody CreateUserRequest request) {
        UserResponse response = users.create(request);
        return ResponseEntity.created(URI.create("/api/users/" + response.id())).body(response);
    }
}
