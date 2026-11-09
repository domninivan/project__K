import { orders } from "../mocks/mockData";

const statusNames = {
  NEW: "Новый",
  READY_FOR_PLANNING: "Готов к планированию",
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
  const [year, month, day] = date.split("-");
  return `${day}.${month}.${year}`;
}

function OrdersPage() {
  return (
    <div>
      <div className="page-title">
        <h2>Список заказов</h2>
        <p>Всего заказов: {orders.length}</p>
      </div>

      <div className="card table-card">
        <table>
          <thead>
            <tr>
              <th>Номер</th>
              <th>Состав</th>
              <th>Желаемая дата</th>
              <th>Статус</th>
              <th>Приоритет</th>
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
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default OrdersPage;
