package ru.mirea.production.product.dao;

import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Repository;
import ru.mirea.production.product.entity.ProductOperation;
import java.util.List;

@Repository
public class ProductOperationDAOImpl implements ProductOperationDAO{

    @Autowired
    private SessionFactory sessionFactory;

    @Override
    public List<ProductOperation> getProductOperationByProductId(int productId) {
        Session session = sessionFactory.getCurrentSession();
        return session
                .createQuery("from ProductOperation where product.id =: product_id order by sequenceNumber", ProductOperation.class)
                .setParameter("product_id", productId)
                .getResultList();
    }
}
