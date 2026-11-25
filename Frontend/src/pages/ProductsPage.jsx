import { products } from "../mocks/mockData";

function ProductsPage() {
  return (
    <div>
      <div className="page-header">
        <div>
          <h2>Изделия</h2>
          <p>Номенклатура выпускаемой продукции.</p>
        </div>
      </div>

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
                <td>{product.description}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default ProductsPage;
