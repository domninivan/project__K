package ru.mirea.production.product.dao;

import ru.mirea.production.product.entity.Product;
import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public class ProductDAOImpl implements ProductDAO{

    @Autowired
    private SessionFactory sessionFactory;

    @Override
    public List<Product> getAllProducts() {
        Session session = sessionFactory.getCurrentSession();
        List<Product> products = session
                .createQuery("from Product", Product.class)
                .getResultList();
        return products;
    }


    @Override
    public Product getProduct(int id) {
        Session session = sessionFactory.getCurrentSession();
        return session.find(Product.class, id);
    }
}
