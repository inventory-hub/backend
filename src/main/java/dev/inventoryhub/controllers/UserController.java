package dev.inventoryhub.controllers;

import dev.inventoryhub.models.User;
import dev.inventoryhub.exceptions.UserNotFoundException;
import dev.inventoryhub.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;

import java.security.Principal;

@RestController
public class UserController {

    @Autowired
    UserRepository repository;

    @GetMapping("/users/{id}")
    public User getUserInfoById(@PathVariable Long id) {
        return repository.findById(id)
                .orElseThrow(() -> new UserNotFoundException(id));
    }

    @GetMapping("/user/login")
    public String home(Principal principal) {
        return "Hello, " + principal.getName();
    }
}
