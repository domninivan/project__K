package ru.mirea.production.user.dao;

import ru.mirea.production.user.entity.AppUser;

public interface AppUserDAO {
    AppUser save(AppUser user);
}
