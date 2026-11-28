import { createContext, useContext, useEffect, useRef, useState } from "react";
import { orders as initialOrders, products } from "../mocks/mockData";
import { createMockPlan, getPlanningActions, updatePlanning } from "../mocks/orderPlanning";

import { isApi } from "../config/dataMode";
import { getOrders, createOrder as createApiOrder } from "../api/ordersApi";

const OrdersContext = createContext(null);

export function OrdersProvider({ children }) {
  const [orders, setOrders] = useState(() => (isApi ? [] : initialOrders).map((order) => ({
    ...order, items: order.items.map((item) => ({ ...item })),
  })));
  const nextOrderId = useRef(Math.max(0, ...initialOrders.map((o) => o.id)) + 1);
  const nextItemId = useRef(Math.max(0, ...initialOrders.flatMap((o) => o.items.map((i) => i.id))) + 1);

  const [loading, setLoading] = useState(isApi);
  const [error, setError] = useState(null);

  async function loadOrders() {
    if (!isApi) return;
    setLoading(true);
    setError(null);
    try { setOrders(await getOrders()); }
    catch (failure) { setError(failure.message); }
    finally { setLoading(false); }
  }

  useEffect(() => { loadOrders(); }, []);

  async function createOrder(payload) {
    if (isApi) {
      const order = await createApiOrder(payload);
      setOrders((previous) => [order, ...previous]);
      return order;
    }
    const items = payload.items.map((item) => {
      const product = products.find((p) => p.id === item.productId);
      if (!product) throw new Error("Изделие отсутствует в справочнике");
      return { id: nextItemId.current++, productId: product.id,
        productTitle: product.title, quantity: item.quantity };
    });
    const order = {
      id: nextOrderId.current++, agreedDate: null, dueDate: payload.dueDate, comment: payload.comment,
      status: "NEW", priority: "MEDIUM", createdAt: new Date().toISOString(), items,
    };
    setOrders((previous) => [order, ...previous]);
    return order;
  }

  function changePlanning(orderId, action, data = {}) {
    if (isApi) return;
    setOrders((previous) => previous.map((order) => {
      if (order.id !== orderId) return order;
      return updatePlanning(order, action, data);
    }));
  }

  async function calculatePlan(orderId) {
    if (isApi) throw new Error("Сценарий согласования доступен только в демонстрационном режиме.");
    const order = orders.find((item) => item.id === orderId);
    if (!order || !getPlanningActions(order).calculate) {
      throw new Error("Расчёт сейчас недоступен. Откройте карточку заказа.");
    }
    if (!order.items.length) throw new Error("В заказе нет изделий для расчёта.");
    changePlanning(orderId, "calculate", { plan: createMockPlan(order) });
  }

  return (
    <OrdersContext.Provider value={{ orders, setOrders, loading, error, loadOrders, createOrder, changePlanning, calculatePlan }}>
      {children}
    </OrdersContext.Provider>
  );
}

export function useOrders() {
  const context = useContext(OrdersContext);
  if (!context) throw new Error("useOrders используется внутри OrdersProvider");
  return context;
}
