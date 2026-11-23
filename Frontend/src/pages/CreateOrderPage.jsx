import { useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { products } from "../mocks/mockData";
import { useOrders } from "../context/OrdersContext";
import OrderForm from "../components/orders/OrderForm";

function CreateOrderPage() {
  const navigate = useNavigate();
  const { createOrder } = useOrders();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState(null);
  const submitting = useRef(false);

  async function handleSubmit(payload) {
    if (submitting.current) return;
    submitting.current = true;
    setIsSubmitting(true);
    setSubmitError(null);
    try {
      const order = await createOrder(payload);
      navigate(`/orders/${order.id}`, { replace: true });
    } catch {
      setSubmitError("Не удалось создать заказ. Попробуйте ещё раз.");
    } finally {
      submitting.current = false;
      setIsSubmitting(false);
    }
  }

  return (
    <div className="order-page">
      <Link to="/orders" className="back-link">← К заказам</Link>
      <div className="page-header"><h2>Новый заказ</h2></div>
      <OrderForm products={products} onSubmit={handleSubmit}
        onCancel={() => navigate("/orders")} isSubmitting={isSubmitting} submitError={submitError} />
    </div>
  );
}

export default CreateOrderPage;
