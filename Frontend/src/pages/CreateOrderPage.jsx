import { orders, products } from "../mocks/mockData";

function CreateOrderPage() {
  const order = orders[0];

  return (
    <div>
      <div className="page-title">
        <h2>Новый заказ</h2>
        <p>Добавьте изделия в заказ, укажите количество и желаемую дату</p>
      </div>

      <div className="card">
        <h3>1. Добавить изделие</h3>
        <div className="add-row">
          <select defaultValue="">
            <option value="" disabled>
              Выберите изделие
            </option>
            {products.map((product) => (
              <option key={product.id} value={product.id}>
                {product.title}
              </option>
            ))}
          </select>
          <button type="button">Добавить</button>
        </div>
      </div>

      <div className="card table-card">
        <h3 className="table-title">2. Состав заказа</h3>
        <table>
          <thead>
            <tr>
              <th>Изделие</th>
              <th>Количество, шт.</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {order.items.map((item) => (
              <tr key={item.id}>
                <td>{item.productTitle}</td>
                <td>
                  <input type="number" min="1" defaultValue={item.quantity} />
                </td>
                <td className="actions">
                  <button type="button" className="delete-btn">
                    Удалить
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="card">
        <h3>3. Параметры заказа</h3>
        <div className="form">
          <label>
            Желаемая дата
            <input type="date" defaultValue={order.dueDate} />
          </label>

          <label>
            Комментарий
            <textarea rows="3" defaultValue={order.comment} />
          </label>
        </div>
      </div>

      <div className="actions">
        <button type="button" className="submit-btn">
          Оформить заказ
        </button>
      </div>
    </div>
  );
}

export default CreateOrderPage;
