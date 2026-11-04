package ru.mirea.production.order.dao;

import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Repository;
import ru.mirea.production.order.entity.ProductionOrder;
import java.util.List;

@Repository
public class ProductionOrderDAOImpl implements ProductionOrderDAO{

    @Autowired
    private SessionFactory sessionFactory;

    @Override
    public List<ProductionOrder> getAllOrders() {
        Session session = sessionFactory.getCurrentSession();
        return session.createQuery(
                "select distinct po from ProductionOrder po " +  
                "inner join fetch po.items i " + 
                "inner join fetch i.product", 
                 ProductionOrder.class)
                        .getResultList();

    }

    @Override
    public ProductionOrder getOrderById(int id) {
        Session session = sessionFactory.getCurrentSession();
        return session.createQuery(
                        "select distinct po from ProductionOrder po " +
                                "inner join fetch po.items i " +
                                "inner join fetch i.product " + 
                                "where po.id = :id", ProductionOrder.class)
                .setParameter("id", id)
                .uniqueResult();
    }

    @Override
    public ProductionOrder saveOrder(ProductionOrder order) {
        Session session = sessionFactory.getCurrentSession();
        return (ProductionOrder) session.merge(order);
    }

}
