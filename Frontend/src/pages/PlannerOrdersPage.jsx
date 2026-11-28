import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useOrders } from "../context/OrdersContext";
import { formatAgreedDate, getOrderStatus } from "../utils/orderDisplay";

export default function PlannerOrdersPage() {
  const { orders, calculatePlan } = useOrders();
  const navigate = useNavigate();
  const [filter, setFilter] = useState("all");
  const [pendingId, setPendingId] = useState(null);
  const [errors, setErrors] = useState({});
  const filteredOrders = orders.filter((order) => {
    if (filter === "unplanned") return order.status === "NEW" && !order.plan;
    if (filter === "decision") return order.status === "NEW" &&
      ["ACCEPTED", "REJECTED"].includes(order.proposal?.status);
    return true;
  });

  async function calculate(id) {
    setPendingId(id);
    setErrors((previous) => ({ ...previous, [id]: null }));
    try {
      await calculatePlan(id);
      navigate(`/planner/orders/${id}`);
    } catch (error) {
      setErrors((previous) => ({ ...previous, [id]: error.message }));
    } finally {
      setPendingId(null);
    }
  }

  return (
    <>
      <div className="page-header">
        <h2>Заказы планировщика</h2>
        <label>Фильтр заказов <select value={filter} onChange={(event) => setFilter(event.target.value)}>
          <option value="all">Все</option>
          <option value="unplanned">Без плана</option>
          <option value="decision">Требуют решения</option>
        </select></label>
      </div>
      {filteredOrders.length === 0 ? <p>Нет заказов по выбранному фильтру.</p> : (
        <div className="card table-wrap">
          <table>
            <thead><tr><th>Номер</th><th>Клиент</th><th>Состав</th><th>Согласованный срок</th>
              <th>Состояние</th><th>Действие</th></tr></thead>
            <tbody>
              {filteredOrders.map((order) => (
                <tr key={order.id}>
                  <td>{order.id}</td><td>{order.customerName || "Тестовый клиент"}</td>
                  <td>{order.items.map((item) => <div key={item.id}>{item.productTitle} — {item.quantity} шт.</div>)}</td>
                  <td>{formatAgreedDate(order.agreedDate)}</td><td>{getOrderStatus(order)}</td>
                  <td>
                    {order.status === "NEW" && !order.plan ? (
                      <button type="button" disabled={pendingId !== null} onClick={() => calculate(order.id)}>
                        {pendingId === order.id ? "Расчёт…" : errors[order.id] ? "Повторить расчёт" : "Рассчитать"}
                      </button>
                    ) : <Link to={`/planner/orders/${order.id}`}>Открыть</Link>}
                    {errors[order.id] && <p className="error" role="alert">{errors[order.id]}</p>}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </>
  );
}
