import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { orders as mockOrders } from "../mocks/mockData";
import { isApi, isMock } from "../config/dataMode";
import { getOrders, cancelOrder } from "../api/ordersApi";

const statusNames = {
  NEW: "Новый",
  READY_FOR_PLANNING: "Готов к планированию",
  PLAN_PROPOSED: "План предложен",
  PLANNED: "Запланирован",
  IN_PRODUCTION: "В производстве",
  COMPLETED: "Выполнен",
  CANCELLED: "Отменён",
};

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

  const [orders, setOrders] = useState(isApi ? [] : mockOrders);
  const [loading, setLoading] = useState(isApi);
  const [error, setError] = useState(null);

  const [cancellingId, setCancellingId] = useState(null);
  const [cancelError, setCancelError] = useState(null);

  const createdOrderId = location.state?.createdOrderId;

  async function loadOrders() {
    setLoading(true);
    setError(null);
    try {
      setOrders(await getOrders());
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    if (isApi) {
      loadOrders();
    }
  }, []);

  async function handleCancel(order) {
    if (!window.confirm(`Отменить заказ № ${order.id}?`)) {
      return;
    }

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
      <div className="page-title">
        <h2>Список заказов</h2>
        <p>Всего заказов: {orders.length}</p>
      </div>

      {createdOrderId && (
        <div role="status">Заказ № {createdOrderId} создан</div>
      )}

      {isApi && loading && <p>Загрузка заказов…</p>}

      {isApi && error && !loading && (
        <div role="alert">
          <p>Не удалось загрузить заказы: {error}</p>
          <button type="button" onClick={loadOrders}>
            Повторить
          </button>
        </div>
      )}

      {cancelError && (
        <div role="alert">
          <p>Не удалось отменить заказ: {cancelError}</p>
        </div>
      )}

      {!loading && !error && orders.length === 0 && <p>Заказов нет.</p>}

      {!loading && !error && orders.length > 0 && (
        <div className="card table-card">
          <table>
            <thead>
              <tr>
                <th>Номер</th>
                <th>Состав</th>
                <th>Желаемая дата</th>
                <th>Статус</th>
                <th>Приоритет</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr key={order.id}>
                  <td>
                    <b>№ {order.id}</b>
                  </td>
                  <td>
                    {order.items.map((item) => (
                      <div key={item.id}>
                        {item.productTitle} — {item.quantity} шт.
                      </div>
                    ))}
                  </td>
                  <td>{formatDate(order.dueDate)}</td>
                  <td>
                    <span
                      className={"badge status-" + order.status.toLowerCase()}
                    >
                      {statusNames[order.status]}
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
                    <button
                      type="button"
                      className="delete-btn"
                      disabled={
                        order.status === "CANCELLED" ||
                        isMock ||
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
