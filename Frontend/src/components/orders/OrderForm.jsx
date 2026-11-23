import { useState } from "react";

function OrderForm({ products, onSubmit, onCancel, isSubmitting = false, submitError = null }) {
  const [selectedId, setSelectedId] = useState("");
  const [items, setItems] = useState([]);
  const [comment, setComment] = useState("");
  const [errors, setErrors] = useState({});

  function addItem() {
    if (isSubmitting) return;
    const product = products.find((p) => p.id === Number(selectedId));
    if (!product) return;
    setItems((previous) => {
      if (previous.some((item) => item.productId === product.id)) {
        return previous.map((item) => item.productId === product.id
          ? { ...item, quantity: Number(item.quantity) + 1 } : item);
      }
      return [...previous, { productId: product.id, quantity: 1 }];
    });
    setErrors({});
  }

  function validate() {
    const result = {};
    if (items.length === 0) result.items = "Добавьте хотя бы одно изделие";
    for (const item of items) {
      if (!products.some((p) => p.id === item.productId)) {
        result["product-" + item.productId] = "Изделие отсутствует в справочнике";
      }
      const quantity = Number(item.quantity);
      if (!Number.isSafeInteger(quantity) || quantity <= 0) {
        result["quantity-" + item.productId] = "Введите целое число больше 0";
      }
    }
    if (comment.length > 200) result.comment = "Комментарий не должен быть длиннее 200 символов";
    return result;
  }

  function handleSubmit(event) {
    event.preventDefault();
    if (isSubmitting) return;
    const nextErrors = validate();
    setErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) return;
    onSubmit({ comment, items: items.map((item) => ({
      productId: products.find((p) => p.id === item.productId).id,
      quantity: Number(item.quantity),
    })) });
  }

  return (
    <form className="card form-card" onSubmit={handleSubmit} noValidate>
      <fieldset disabled={isSubmitting}>
        <h3>Состав заказа</h3>
        <div className="add-row">
          <label htmlFor="product">Изделие</label>
          <select id="product" value={selectedId} onChange={(e) => setSelectedId(e.target.value)}>
            <option value="">Выберите изделие</option>
            {products.map((p) => <option key={p.id} value={p.id}>{p.title}</option>)}
          </select>
          <button type="button" className="outline-btn" disabled={!selectedId} onClick={addItem}>Добавить позицию</button>
        </div>
        <div className="table-wrap">
          <table>
            <thead><tr><th>Изделие</th><th>Количество</th><th className="center">Действия</th></tr></thead>
            <tbody>
              {items.length === 0 && <tr><td colSpan="3" className="empty-state">Добавьте изделия в заказ</td></tr>}
              {items.map((item) => {
                const title = products.find((p) => p.id === item.productId)?.title || "Неизвестное изделие";
                return (
                  <tr key={item.productId}>
                    <td>{title}{errors["product-" + item.productId] && <p className="error" role="alert">{errors["product-" + item.productId]}</p>}</td>
                    <td>
                      <input type="number" min="1" step="1" aria-label={`Количество: ${title}`}
                        aria-invalid={Boolean(errors["quantity-" + item.productId])}
                        aria-describedby={errors["quantity-" + item.productId] ? `quantity-error-${item.productId}` : undefined}
                        value={item.quantity} onChange={(e) => setItems(items.map((row) => row.productId === item.productId ? { ...row, quantity: e.target.value } : row))} />
                      {errors["quantity-" + item.productId] && <p id={`quantity-error-${item.productId}`} className="error" role="alert">{errors["quantity-" + item.productId]}</p>}
                    </td>
                    <td className="center"><button type="button" className="delete-btn" onClick={() => { setItems(items.filter((row) => row.productId !== item.productId)); setErrors({}); }}>Удалить</button></td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        {errors.items && <p className="error" role="alert">{errors.items}</p>}
        <div className="form-fields">
          <label className="comment-field">Комментарий
            <textarea rows="2" value={comment} onChange={(e) => setComment(e.target.value)} aria-invalid={Boolean(errors.comment)} aria-describedby={errors.comment ? "comment-error" : undefined} />
            <span className="field-hint">{comment.length} / 200</span>
            {errors.comment && <span id="comment-error" className="error" role="alert">{errors.comment}</span>}
          </label>
        </div>
        {submitError && <p className="error" role="alert">{submitError}</p>}
        <div className="form-buttons">
          <button type="button" className="secondary" onClick={onCancel}>Отмена</button>
          <button type="submit">{isSubmitting ? "Создать…" : "Создать заказ"}</button>
        </div>
      </fieldset>
    </form>
  );
}

export default OrderForm;
