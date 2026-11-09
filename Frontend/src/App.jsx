import { BrowserRouter, Routes, Route, Navigate, NavLink } from "react-router-dom";
import ProductsPage from "./pages/ProductsPage";
import CreateOrderPage from "./pages/CreateOrderPage";
import OrdersPage from "./pages/OrdersPage";
import "./App.css";

function App() {
  return (
    <BrowserRouter>
      <header className="header">
        <h1>Мастер-план</h1>
        <nav>
          <NavLink to="/products">Каталог</NavLink>
          <NavLink to="/orders/new">Новый заказ</NavLink>
          <NavLink to="/orders" end>
            Заказы
          </NavLink>
        </nav>
      </header>

      <main className="content">
        <Routes>
          <Route path="/" element={<Navigate to="/products" />} />
          <Route path="/products" element={<ProductsPage />} />
          <Route path="/orders/new" element={<CreateOrderPage />} />
          <Route path="/orders" element={<OrdersPage />} />
        </Routes>
      </main>
    </BrowserRouter>
  );
}

export default App;
