package ru.mirea.production.planning.exception;

public class GoServiceUnavailableException extends RuntimeException {

    public GoServiceUnavailableException() {
        super("Сервис планирования временно недоступен");
    }
}