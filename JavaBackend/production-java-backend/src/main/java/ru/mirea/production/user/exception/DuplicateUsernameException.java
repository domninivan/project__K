package ru.mirea.production.user.exception;

public class DuplicateUsernameException extends RuntimeException {
    public DuplicateUsernameException() { super("Пользователь с таким логином уже существует."); }
}
