import { useEffect, useState } from "react";
import { products as mockProducts } from "../mocks/mockData";
import { isApi } from "../config/dataMode";
import { getProducts } from "../api/productsApi";

function ProductsPage() {
  const [products, setProducts] = useState(isApi ? [] : mockProducts);
  const [loading, setLoading] = useState(isApi);
  const [error, setError] = useState(null);

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

  return (
    <div>
      <div className="page-header">
        <div>
          <h2>Изделия</h2>
          <p>Номенклатура выпускаемой продукции.</p>
        </div>
      </div>

      {isApi && loading && <p>Загрузка каталога…</p>}

      {isApi && error && !loading && (
        <div role="alert">
          <p className="error">Не удалось загрузить каталог: {error}</p>
          <button type="button" className="outline-btn" onClick={loadProducts}>
            Повторить
          </button>
        </div>
      )}

      {!loading && !error && products.length === 0 && (
        <p className="empty-state">Каталог пуст.</p>
      )}

      {!loading && !error && products.length > 0 && (
        <div className="card table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Наименование</th>
                <th>Описание</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <tr key={product.id}>
                  <td>{product.id}</td>
                  <td>{product.title}</td>
                  <td>{product.description || "Описание не указано"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default ProductsPage;
