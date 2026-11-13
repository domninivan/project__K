import { useState } from "react";

// Возвращает сегодняшнюю дату в локальном формате "YYYY-MM-DD".
function todayStr() {
  const now = new Date();
  const y = now.getFullYear();
  const m = String(now.getMonth() + 1).padStart(2, "0");
  const d = String(now.getDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
}

// Базовая форма заказа.
//
// Принимает изделия из каталога и по сабмиту вызывает onSubmit(payload),
// где payload = { dueDate, comment, items: [{ productId, quantity }] }.
// Внутренняя логика (выбор изделий, добавление/удаление позиций,
// изменение количества, объединение дублей, валидация) живёт здесь,
// а не на странице.
function OrderForm({ products, onSubmit, isSubmitting, submitError }) {
  const [selectedProductId, setSelectedProductId] = useState("");
  const [items, setItems] = useState([]);
  const [dueDate, setDueDate] = useState("");
  const [comment, setComment] = useState("");

  const today = todayStr();
  const dateIsValid = dueDate !== "" && dueDate >= today;
  const itemsAreValid =
    items.length > 0 && items.every((item) => item.quantity >= 1);
  const canSubmit = dateIsValid && itemsAreValid && !isSubmitting;

  function productTitle(productId) {
    const product = products.find((p) => Number(p.id) === Number(productId));
    return product ? product.title : `Изделие #${productId}`;
  }

  function handleAdd(event) {
    event.preventDefault();
    const id = Number(selectedProductId);
    if (!id) return;

    setItems((prev) => {
      const existing = prev.find((item) => item.productId === id);
      if (existing) {
        // Повторно выбранное изделие: увеличиваем количество.
        return prev.map((item) =>
          item.productId === id
            ? { ...item, quantity: item.quantity + 1 }
            : item
        );
      }
      return [...prev, { productId: id, quantity: 1 }];
    });
    setSelectedProductId("");
  }

  function handleRemove(productId) {
    setItems((prev) => prev.filter((item) => item.productId !== productId));
  }

  function handleQuantityChange(productId, value) {
    const qty = Number(value);
    const next = Number.isFinite(qty) && qty >= 1 ? qty : 1;
    setItems((prev) =>
      prev.map((item) =>
        item.productId === productId ? { ...item, quantity: next } : item
      )
    );
  }

  function handleSubmit(event) {
    event.preventDefault();
    if (!canSubmit) return;

    onSubmit({
      dueDate,
      comment,
      items: items.map(({ productId, quantity }) => ({ productId, quantity })),
    });
  }

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label htmlFor="order-product">Изделие</label>
        <select
          id="order-product"
          value={selectedProductId}
          onChange={(event) => setSelectedProductId(event.target.value)}
        >
          <option value="" disabled>
            Выберите изделие
          </option>
          {products.map((product) => (
            <option key={product.id} value={product.id}>
              {product.title}
            </option>
          ))}
        </select>
        <button type="button" onClick={handleAdd} disabled={!selectedProductId}>
          Добавить
        </button>
      </div>

      {items.length === 0 ? (
        <p>Позиции не добавлены.</p>
      ) : (
        <ul>
          {items.map((item) => (
            <li key={item.productId}>
              {productTitle(item.productId)}{" "}
              <input
                type="number"
                min="1"
                value={item.quantity}
                onChange={(event) =>
                  handleQuantityChange(item.productId, event.target.value)
                }
              />{" "}
              шт.{" "}
              <button
                type="button"
                onClick={() => handleRemove(item.productId)}
              >
                Удалить
              </button>
            </li>
          ))}
        </ul>
      )}

      <div>
        <label htmlFor="order-due-date">Желаемая дата</label>
        <input
          id="order-due-date"
          type="date"
          min={today}
          value={dueDate}
          onChange={(event) => setDueDate(event.target.value)}
        />
        {dueDate !== "" && !dateIsValid && <p>Дата не может быть в прошлом.</p>}
      </div>

      <div>
        <label htmlFor="order-comment">Комментарий</label>
        <textarea
          id="order-comment"
          rows={3}
          maxLength={200}
          value={comment}
          onChange={(event) => setComment(event.target.value)}
        />
      </div>

      {submitError && (
        <div role="alert">
          <strong>Ошибка:</strong> {submitError}
        </div>
      )}

      <button type="submit" disabled={!canSubmit}>
        {isSubmitting ? "Отправка…" : "Оформить заказ"}
      </button>
    </form>
  );
}

export default OrderForm;
