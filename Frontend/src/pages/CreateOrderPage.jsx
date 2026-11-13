import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import OrderForm from "../components/orders/OrderForm";
import { products as mockProducts } from "../mocks/mockData";
import { isApi, isMock } from "../config/dataMode";
import { getProducts } from "../api/productsApi";
import { createOrder } from "../api/ordersApi";

function CreateOrderPage() {
  const navigate = useNavigate();

  const [products, setProducts] = useState(isApi ? [] : mockProducts);
  const [loading, setLoading] = useState(isApi);
  const [error, setError] = useState(null);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState(null);

  async function loadProducts() {
    setLoading(true);
    setError(null);
    try {
      setProducts(await getProducts());
    } catch (err) {
      setError(err.message);
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
    // Защита от повторной отправки.
    if (isSubmitting) {
      return;
    }

    // Демонстрационный режим: запросы не выполняются.
    if (isMock) {
      setSubmitError("Демонстрационный режим: сохранение недоступно");
      return;
    }

    setSubmitError(null);
    setIsSubmitting(true);
    try {
      const created = await createOrder(payload);
      navigate("/orders", { state: { createdOrderId: created.id } });
    } catch (err) {
      setSubmitError(err.message);
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div>
      <div className="page-title">
        <h2>Новый заказ</h2>
        <p>Добавьте изделия в заказ, укажите количество и желаемую дату</p>
      </div>

      {isApi && loading && <p>Загрузка каталога…</p>}

      {isApi && error && !loading && (
        <div role="alert">
          <p>Не удалось загрузить каталог: {error}</p>
          <button type="button" onClick={loadProducts}>
            Повторить
          </button>
        </div>
      )}

      {!loading && !error && products.length === 0 && (
        <p>Каталог пуст — оформление заказа пока недоступно.</p>
      )}

      {!loading && !error && products.length > 0 && (
        <OrderForm
          products={products}
          onSubmit={handleSubmit}
          isSubmitting={isSubmitting}
          submitError={submitError}
        />
      )}
    </div>
  );
}

export default CreateOrderPage;
