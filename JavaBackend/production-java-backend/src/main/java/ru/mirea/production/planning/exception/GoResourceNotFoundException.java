package ru.mirea.production.planning.exception;

public class GoResourceNotFoundException extends RuntimeException {
    public GoResourceNotFoundException(String message) {
        super(message);
    }
}
