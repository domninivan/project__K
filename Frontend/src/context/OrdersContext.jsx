import { createContext, useContext, useRef, useState } from "react";
import { orders as initialOrders, products } from "../mocks/mockData";
import { createMockPlan, getPlanningActions, updatePlanning } from "../mocks/orderPlanning";

const OrdersContext = createContext(null);

export function OrdersProvider({ children }) {
  const [orders, setOrders] = useState(() => initialOrders.map((order) => ({
    ...order, items: order.items.map((item) => ({ ...item })),
  })));
  const nextOrderId = useRef(Math.max(0, ...initialOrders.map((o) => o.id)) + 1);
  const nextItemId = useRef(Math.max(0, ...initialOrders.flatMap((o) => o.items.map((i) => i.id))) + 1);

  function createOrder(payload) {
    const items = payload.items.map((item) => {
      const product = products.find((p) => p.id === item.productId);
      if (!product) throw new Error("Изделие отсутствует в справочнике");
      return { id: nextItemId.current++, productId: product.id,
        productTitle: product.title, quantity: item.quantity };
    });
    const order = {
      id: nextOrderId.current++, agreedDate: null, comment: payload.comment,
      status: "NEW", priority: "MEDIUM", createdAt: new Date().toISOString(), items,
    };
    setOrders((previous) => [order, ...previous]);
    return order;
  }

  function changePlanning(orderId, action, data = {}) {
    setOrders((previous) => previous.map((order) => {
      if (order.id !== orderId) return order;
      return updatePlanning(order, action, data);
    }));
  }

  async function calculatePlan(orderId) {
    const order = orders.find((item) => item.id === orderId);
    if (!order || !getPlanningActions(order).calculate) {
      throw new Error("Расчёт сейчас недоступен. Откройте карточку заказа.");
    }
    if (!order.items.length) throw new Error("В заказе нет изделий для расчёта.");
    changePlanning(orderId, "calculate", { plan: createMockPlan(order) });
  }

  return (
    <OrdersContext.Provider value={{ orders, createOrder, changePlanning, calculatePlan }}>
      {children}
    </OrdersContext.Provider>
  );
}

export function useOrders() {
  const context = useContext(OrdersContext);
  if (!context) throw new Error("useOrders используется внутри OrdersProvider");
  return context;
}
