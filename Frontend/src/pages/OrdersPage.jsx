import { Link } from "react-router-dom";
import { useOrders } from "../context/OrdersContext";
import { getOrderStatus, formatAgreedDate } from "../utils/orderDisplay";

const priorityNames = {
  HIGH: "Высокий",
  MEDIUM: "Средний",
  LOW: "Низкий",
};

function OrdersPage() {
  const { orders } = useOrders();
  return (
    <div>
      <div className="page-header">
        <div>
          <h2>Мои заказы</h2>
          <p>Всего заказов: {orders.length}</p>
        </div>
        <Link to="/orders/new" className="button-link">
          + Создать заказ
        </Link>
      </div>

      <div className="card table-wrap">
        <table>
          <thead>
            <tr>
              <th>Номер</th>
              <th>Состав заказа</th>
              <th>Согласованная дата</th>
              <th>Статус</th>
              <th>Приоритет</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {orders.map((order) => (
              <tr key={order.id}>
                <td>{order.id}</td>
                <td>
                  {order.items.map((item) => (
                    <div key={item.id}>
                      {item.productTitle} — {item.quantity} шт.
                    </div>
                  ))}
                </td>
                <td>{formatAgreedDate(order.agreedDate)}</td>
                <td>
                  <span
                    className={"badge status-" + order.status.toLowerCase()}
                  >
                    {getOrderStatus(order)}
                  </span>
                </td>
                <td>
                  <span
                    className={"badge priority-" + order.priority.toLowerCase()}
                  >
                    {priorityNames[order.priority]}
                  </span>
                </td>
                <td><Link to={`/orders/${order.id}`}>Открыть</Link></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default OrdersPage;
