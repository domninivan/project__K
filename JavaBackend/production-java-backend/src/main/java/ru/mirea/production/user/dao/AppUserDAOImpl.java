package ru.mirea.production.user.dao;

import org.hibernate.SessionFactory;
import org.hibernate.exception.ConstraintViolationException;
import org.springframework.stereotype.Repository;
import ru.mirea.production.user.entity.AppUser;
import ru.mirea.production.user.exception.DuplicateUsernameException;

@Repository
public class AppUserDAOImpl implements AppUserDAO {
    private final SessionFactory sessionFactory;

    public AppUserDAOImpl(SessionFactory sessionFactory) { this.sessionFactory = sessionFactory; }

    @Override
    public AppUser save(AppUser user) {
        var session = sessionFactory.getCurrentSession();
        try {
            session.persist(user);
            session.flush();
            return user;
        } catch (RuntimeException ex) {
            // The UNIQUE constraint also protects concurrent requests for the same login.
            for (Throwable cause = ex; cause != null; cause = cause.getCause()) {
                if (cause instanceof ConstraintViolationException constraint
                        && "23505".equals(constraint.getSQLState())
                        && "app_user_username_key".equals(constraint.getConstraintName())) {
                    throw new DuplicateUsernameException();
                }
            }
            throw ex;
        }
    }
}
