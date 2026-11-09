import { Link } from "react-router-dom";
import { products } from "../mocks/mockData";

function ProductsPage() {
  return (
    <div>
      <div className="page-title">
        <h2>Каталог изделий</h2>
        <p>Выберите изделия, которые нужно добавить в заказ</p>
      </div>

      <div className="products">
        {products.map((product) => (
          <div className="product-card" key={product.id}>
            <div className="product-icon">{product.title[0]}</div>
            <h3>{product.title}</h3>
            <p>{product.description}</p>
            <Link to="/orders/new" className="button-link">
              Добавить в заказ
            </Link>
          </div>
        ))}
      </div>
    </div>
  );
}

export default ProductsPage;
