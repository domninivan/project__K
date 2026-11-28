package ru.mirea.production.planning.exception;

public class GoConflictException extends RuntimeException {
    public GoConflictException(String message) {
        super(message);
    }
}
