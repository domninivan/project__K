import { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { products as mockProducts } from "../mocks/mockData";
import { isApi } from "../config/dataMode";
import { getProducts } from "../api/productsApi";
import { useOrders } from "../context/OrdersContext";
import OrderForm from "../components/orders/OrderForm";

function CreateOrderPage() {
  const navigate = useNavigate();
  const { createOrder } = useOrders();

  const [products, setProducts] = useState(isApi ? [] : mockProducts);
  const [loading, setLoading] = useState(isApi);
  const [loadError, setLoadError] = useState(null);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState(null);
  const submitting = useRef(false);

  async function loadProducts() {
    setLoading(true);
    setLoadError(null);
    try {
      setProducts(await getProducts());
    } catch (err) {
      setLoadError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    if (isApi) {
      loadProducts();
    }
  }, []);

  async function handleSubmit(payload) {
    if (submitting.current) return;

    submitting.current = true;
    setIsSubmitting(true);
    setSubmitError(null);
    try {
      const created = await createOrder(payload);
      navigate(`/orders/${created.id}`, { replace: true });
    } catch (err) {
      setSubmitError(err.message);
    } finally {
      submitting.current = false;
      setIsSubmitting(false);
    }
  }

  return (
    <div className="order-page">
      <Link to="/orders" className="back-link">← К заказам</Link>
      <div className="page-header"><h2>Новый заказ</h2></div>

      {isApi && loading && <p>Загрузка каталога…</p>}

      {isApi && loadError && !loading && (
        <div role="alert">
          <p className="error">Не удалось загрузить каталог: {loadError}</p>
          <button type="button" className="outline-btn" onClick={loadProducts}>
            Повторить
          </button>
        </div>
      )}

      {!loading && !loadError && products.length === 0 && (
        <p className="empty-state">Каталог пуст — оформление заказа пока недоступно.</p>
      )}

      {!loading && !loadError && products.length > 0 && (
        <OrderForm
          products={products}
          onSubmit={handleSubmit}
          onCancel={() => navigate("/orders")}
          isSubmitting={isSubmitting}
          submitError={submitError}
        />
      )}
    </div>
  );
}

export default CreateOrderPage;
