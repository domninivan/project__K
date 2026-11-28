import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { getOrderStatus, formatAgreedDate } from "../utils/orderDisplay";
import { useOrders } from "../context/OrdersContext";
import { isApi, isMock } from "../config/dataMode";
import { cancelOrder } from "../api/ordersApi";

const priorityNames = {
  HIGH: "Высокий",
  MEDIUM: "Средний",
  LOW: "Низкий",
};

function formatDate(date) {
  if (!date) return "—";
  const [year, month, day] = date.split("-");
  return `${day}.${month}.${year}`;
}

function OrdersPage() {
  const location = useLocation();
  const { orders, setOrders, loading, error, loadOrders } = useOrders();

  const [cancellingId, setCancellingId] = useState(null);
  const [cancelError, setCancelError] = useState(null);

  const createdOrderId = location.state?.createdOrderId;

  async function handleCancel(order) {
    if (!window.confirm(`Отменить заказ № ${order.id}?`)) return;

    setCancelError(null);
    setCancellingId(order.id);
    try {
      const updated = await cancelOrder(order.id);
      // Заменяем строку объектом из ответа сервера (status = CANCELLED).
      setOrders((prev) =>
        prev.map((item) => (item.id === order.id ? updated : item))
      );
    } catch (err) {
      setCancelError(err.message);
    } finally {
      setCancellingId(null);
    }
  }

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

      {createdOrderId && <p role="status">Заказ № {createdOrderId} создан</p>}

      {isApi && loading && <p>Загрузка заказов…</p>}

      {isApi && error && !loading && (
        <div role="alert">
          <p className="error">Не удалось загрузить заказы: {error}</p>
          <button type="button" className="outline-btn" onClick={loadOrders}>
            Повторить
          </button>
        </div>
      )}

      {cancelError && (
        <p className="error" role="alert">
          Не удалось отменить заказ: {cancelError}
        </p>
      )}

      {!loading && !error && orders.length === 0 && (
        <p className="empty-state">Заказов нет.</p>
      )}

      {!loading && !error && orders.length > 0 && (
        <div className="card table-wrap">
          <table>
            <thead>
              <tr>
                <th>Номер</th>
                <th>Состав заказа</th>
                <th>Желаемая дата</th>
                <th>Согласованная дата</th>
                <th>Статус</th>
                <th>Приоритет</th>
                <th></th>
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
                  <td>{formatDate(order.dueDate)}</td>
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
                  <td>
                    <Link to={`/orders/${order.id}`}>Открыть</Link>{" "}
                    <button
                      type="button"
                      className="delete-btn"
                      disabled={
                        isMock ||
                        order.status === "CANCELLED" ||
                        cancellingId === order.id
                      }
                      onClick={() => handleCancel(order)}
                    >
                      {cancellingId === order.id ? "Отмена…" : "Отменить"}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default OrdersPage;
