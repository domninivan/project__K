package ru.mirea.production.planning.exception;

public class GoBadRequestException extends RuntimeException {
    public GoBadRequestException(String message) {
        super(message);
    }
}
