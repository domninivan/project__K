import { BrowserRouter, Routes, Route, Navigate, NavLink } from "react-router-dom";
import ProductsPage from "./pages/ProductsPage";
import CreateOrderPage from "./pages/CreateOrderPage";
import OrdersPage from "./pages/OrdersPage";
import "./App.css";
import { OrdersProvider } from "./context/OrdersContext";

function App() {
  return (
    <BrowserRouter>
      <OrdersProvider>
      <header className="header">
        <h1>Мастер-план</h1>
        <nav>
          <NavLink to="/orders">Заказы</NavLink>
          <NavLink to="/products">Изделия</NavLink>
        </nav>
      </header>

      <main className="content">
        <Routes>
          <Route path="/" element={<Navigate to="/orders" />} />
          <Route path="/orders" element={<OrdersPage />} />
          <Route path="/orders/new" element={<CreateOrderPage />} />
          <Route path="/products" element={<ProductsPage />} />
        </Routes>
      </main>
      </OrdersProvider>
    </BrowserRouter>
  );
}

export default App;
