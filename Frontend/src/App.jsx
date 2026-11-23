import { BrowserRouter, Routes, Route, Navigate, NavLink, useLocation, useNavigate } from "react-router-dom";
import ProductsPage from "./pages/ProductsPage";
import CreateOrderPage from "./pages/CreateOrderPage";
import OrdersPage from "./pages/OrdersPage";
import OrderDetailsPage from "./pages/OrderDetailsPage";
import PlannerOrderPage from "./pages/PlannerOrderPage";
import PlannerOrdersPage from "./pages/PlannerOrdersPage";
import "./App.css";
import { OrdersProvider } from "./context/OrdersContext";

const rolePaths = {
  client: "/orders",
  planner: "/planner/orders",
};

const roleMenus = {
  client: "Мои заказы",
  planner: "Заказы",
};

function Header() {
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const section = pathname.split("/")[1];
  const role = Object.hasOwn(rolePaths, section) ? section : "client";

  return (
    <header className="header workflow-header">
      <h1>Мастер-план</h1>
      <nav>
        <NavLink to={rolePaths[role]}>
          {roleMenus[role]}
        </NavLink>
        {role === "client" && <NavLink to="/products">Изделия</NavLink>}
      </nav>
      <label className="role-picker">
        Тестовая роль
        <select value={role} onChange={(event) => navigate(rolePaths[event.target.value])}>
          <option value="client">Клиент</option>
          <option value="planner">Планировщик</option>
        </select>
      </label>
    </header>
  );
}

function App() {
  return (
    <BrowserRouter>
      <OrdersProvider>
      <Header />

      <main className="content">
        <Routes>
          <Route path="/" element={<Navigate to="/orders" />} />
          <Route path="/orders" element={<OrdersPage />} />
          <Route path="/orders/new" element={<CreateOrderPage />} />
          <Route path="/orders/:orderId" element={<OrderDetailsPage />} />
          <Route path="/planner/orders" element={<PlannerOrdersPage />} />
          <Route path="/planner/orders/:orderId" element={<PlannerOrderPage />} />
          <Route path="/products" element={<ProductsPage />} />
        </Routes>
      </main>
      </OrdersProvider>
    </BrowserRouter>
  );
}

export default App;
