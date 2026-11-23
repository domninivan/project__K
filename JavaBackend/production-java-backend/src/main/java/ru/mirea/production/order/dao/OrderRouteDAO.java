package ru.mirea.production.order.dao;

import jakarta.persistence.LockModeType;
import java.util.List;
import org.hibernate.SessionFactory;
import org.springframework.stereotype.Repository;
import ru.mirea.production.order.entity.*;

@Repository
public class OrderRouteDAO {
    private final SessionFactory sessions;
    public OrderRouteDAO(SessionFactory sessions) { this.sessions = sessions; }

    public ProductionOrder lockOrder(int orderId) {
        var session = sessions.getCurrentSession();
        ProductionOrder order = session.find(ProductionOrder.class, orderId, LockModeType.PESSIMISTIC_WRITE);
        if (order != null) order.getItems().size();
        return order;
    }

    public List<OrderItemOperation> find(int orderId) {
        return sessions.getCurrentSession().createQuery(
                "from OrderItemOperation o where o.orderItem.order.id = :id "
                        + "order by o.orderItem.id, o.sequenceNumber", OrderItemOperation.class)
                .setParameter("id", orderId).getResultList();
    }
    public void save(OrderItemOperation operation) { sessions.getCurrentSession().persist(operation); }
}
