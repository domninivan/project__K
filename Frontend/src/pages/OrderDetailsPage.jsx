import { Link, useParams } from "react-router-dom";
import { useOrders } from "../context/OrdersContext";
import { getOrderStatus, formatAgreedDate } from "../utils/orderDisplay";
import { getPlanningActions } from "../mocks/orderPlanning";

function DateProposal({ order, onRespond }) {
  const { proposal } = order;
  if (!proposal) return null;
  const canRespond = getPlanningActions(order).respond;

  function respond(event) {
    event.preventDefault();
    const accepted = event.nativeEvent.submitter?.value === "accept";
    const response = new FormData(event.currentTarget).get("response").trim();
    onRespond({ version: proposal.version, accepted, response });
  }

  return (
    <section className="card order-details planning-panel">
      <h3>Предложение срока</h3>
      <p>{formatAgreedDate(proposal.date)}</p>
      {proposal.reason && <p className="order-comment">{proposal.reason}</p>}
      {proposal.status === "ACCEPTED" && <p>Дата принята.</p>}
      {proposal.status === "REJECTED" && <p>Дата отклонена. Можно дождаться нового предложения.</p>}
      {canRespond && (
        <form className="planning-form" onSubmit={respond}>
          <label>Пояснение при отказе (необязательно)<textarea name="response" rows="2" maxLength="500" /></label>
          <div className="planning-actions">
            <button type="submit" value="reject" className="secondary">Отклонить дату</button>
            <button type="submit" value="accept">Принять дату</button>
          </div>
        </form>
      )}
    </section>
  );
}

function OrderDetailsPage() {
  const { orderId } = useParams();
  const { orders, changePlanning } = useOrders();
  const order = orders.find((item) => String(item.id) === orderId);

  return (
    <div className="order-page">
      <Link to="/orders" className="back-link">← Назад к заказам</Link>
      {!order ? (
        <h2>Заказ не найден</h2>
      ) : (
        <>
          <div className="page-header"><h2>Заказ № {order.id}</h2></div>
          <div className="card order-details">
            <dl className="order-properties">
              <div>
                <dt>Статус</dt>
                <dd><span className={`badge status-${order.status.toLowerCase()}`}>
                  {getOrderStatus(order)}
                </span></dd>
              </div>
              <div>
                <dt>Согласованная дата</dt>
                <dd>{formatAgreedDate(order.agreedDate)}</dd>
              </div>
            </dl>
            <h3>Состав заказа</h3>
            <div className="table-wrap">
              <table>
                <thead><tr><th>Изделие</th><th>Количество</th></tr></thead>
                <tbody>
                  {order.items.map((item) => (
                    <tr key={item.id}><td>{item.productTitle}</td><td>{item.quantity} шт.</td></tr>
                  ))}
                </tbody>
              </table>
            </div>
            <h3>Комментарий</h3>
            <p className="order-comment">{order.comment || "Не указан"}</p>
          </div>
          <DateProposal key={order.proposal?.version} order={order}
            onRespond={(data) => changePlanning(order.id, "respond", data)} />
        </>
      )}
    </div>
  );
}

export default OrderDetailsPage;
