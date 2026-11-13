import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
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
      <div className="page-title">
        <h2>Каталог изделий</h2>
        <p>Выберите изделия, которые нужно добавить в заказ</p>
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

      {!loading && !error && products.length === 0 && <p>Каталог пуст.</p>}

      {!loading && !error && products.length > 0 && (
        <div className="products">
          {products.map((product) => (
            <div className="product-card" key={product.id}>
              <div className="product-icon">
                {(product.title || "?").charAt(0)}
              </div>
              <h3>{product.title}</h3>
              <p>ID: {product.id}</p>
              <p>{product.description || "Описание не указано"}</p>
              <Link to="/orders/new" className="button-link">
                Перейти к оформлению
              </Link>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default ProductsPage;
