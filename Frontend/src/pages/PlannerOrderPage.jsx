import { isApi } from "../config/dataMode";
import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { useOrders } from "../context/OrdersContext";
import { getPlanningActions } from "../mocks/orderPlanning";
import { formatAgreedDate, getOrderStatus } from "../utils/orderDisplay";
import PlanGantt from "../components/orders/PlanGantt";

export default function PlannerOrderPage() {
  const { orderId } = useParams();
  const { orders, changePlanning, calculatePlan } = useOrders();
  const [historyVersion, setHistoryVersion] = useState("");
  const [calculating, setCalculating] = useState(false);
  const [error, setError] = useState("");
  const order = orders.find((item) => String(item.id) === orderId);

  if (!order) {
    return <div className="order-page"><Link to="/planner/orders">Назад к заказам</Link><h2>Заказ не найден</h2></div>;
  }

  const { plan, proposal } = order;
  const history = order.planHistory ?? [];
  const selectedPlan = history.find((item) => String(item.version) === historyVersion) ?? plan;
  const actions = isApi ? {} : getPlanningActions(order);
  async function calculate() {
    setError("");
    setCalculating(true);
    try {
      await calculatePlan(order.id);
      setHistoryVersion("");
    } catch (failure) {
      setError(failure.message);
    } finally {
      setCalculating(false);
    }
  }
  function proposeDate(event) {
    event.preventDefault();
    const reason = new FormData(event.currentTarget).get("reason").trim();
    changePlanning(order.id, "propose", { version: plan.version, reason });
  }

  return (
    <div className="order-page">
      <Link to="/planner/orders" className="back-link">← Назад к заказам</Link>
      <div className="page-header"><h2>Заказ № {order.id}</h2></div>
      <section className="card order-details">
        <p>Клиент: {order.customerName || "Тестовый клиент"}</p>
        <p>Статус: {getOrderStatus(order)}</p>
        <p>Первоначально согласованная дата: {formatAgreedDate(order.originalAgreedDate)}</p>
        <p>Согласованная дата: {formatAgreedDate(order.agreedDate)}</p>
        <h3>Состав заказа</h3>
        <ul>
          {order.items.map((item) => <li key={item.id}>{item.productTitle} — {item.quantity} шт.</li>)}
        </ul>
        <h3>Комментарий к заказу</h3>
        <p className="order-comment">{order.comment || "Не указан"}</p>
      </section>
      <section className="card order-details planning-panel">
        <h3>Производственный план</h3>
        {isApi && <p>Согласование плана доступно в демонстрационном режиме.</p>}
        {!plan && <p>План ещё не рассчитан.</p>}
        {selectedPlan && (
          <>
            {history.length > 0 && (
              <label>Версия плана <select value={historyVersion} onChange={(event) => setHistoryVersion(event.target.value)}>
                <option value="">{plan.status === "ACCEPTED" ? "Действующий" : "Черновик"} · версия {plan.version}</option>
                {history.map((item) => <option key={item.version} value={item.version}>История · версия {item.version}</option>)}
              </select></label>
            )}
            <p>{historyVersion ? "История" : plan.status === "ACCEPTED" ? "Действующий план" : plan.status === "REJECTED" ? "Отклонённый план" : "Черновик"} · версия {selectedPlan.version} · Прогноз завершения: {formatAgreedDate(selectedPlan.finishDate)}</p>
            <div className="table-wrap">
              <table>
                <thead><tr><th>Изделие</th><th>Операция</th><th>Станок</th><th>Время (МСК)</th><th>Количество</th></tr></thead>
                <tbody>
                  {selectedPlan.operations.map((operation) => (
                    <tr key={operation.id}>
                      <td>{operation.productTitle}</td><td>{operation.name}</td><td>{operation.machine}</td>
                      <td>{operation.start}–{operation.end}</td><td>{operation.quantity} шт.</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <PlanGantt plan={selectedPlan} />
          </>
        )}
        {proposal?.status === "PENDING" && <p>Ожидаем ответа клиента.</p>}
        {proposal?.status === "REJECTED" && (
          <p>Клиент отклонил дату {formatAgreedDate(proposal.date)}.
            {proposal.response && <> Причина: {proposal.response}</>}
          </p>
        )}
        {actions.approve && <p>Клиент согласовал дату. План можно утвердить.</p>}
        {plan?.status === "ACCEPTED" && <p>План утверждён.</p>}
        {actions.calculate && (
          <button type="button" disabled={calculating} onClick={calculate}>
            {calculating ? "Расчёт…" : plan ? "Пересчитать план" : "Рассчитать план"}
          </button>
        )}
        {error && <p className="error" role="alert">{error}</p>}
        {!historyVersion && actions.propose && (
          <form className="planning-form" key={plan.version} onSubmit={proposeDate}>
            <p>Предлагаемая дата: {formatAgreedDate(plan.finishDate)}</p>
            <label>Комментарий к предложению<textarea name="reason" rows="2" maxLength="500" /></label>
            <button type="submit">Предложить срок клиенту</button>
          </form>
        )}
        {!historyVersion && actions.approve && (
          <button type="button" onClick={() => changePlanning(order.id, "approve", { version: plan.version })}>
            Утвердить план
          </button>
        )}
      </section>
    </div>
  );
}
