import { requestJson } from "./httpClient";

// Возвращает массив заказов.
export async function getOrders() {
  return requestJson("/api/orders");
}

// Создаёт заказ. payload уже имеет вид:
// { dueDate, comment, items: [{ productId, quantity }] }.
export async function createOrder(payload) {
  return requestJson("/api/orders", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });
}

// Отменяет заказ. Возвращает заказ со статусом CANCELLED.
export async function cancelOrder(id) {
  return requestJson(`/api/orders/${id}`, {
    method: "DELETE",
  });
}
