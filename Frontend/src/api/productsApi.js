import { requestJson } from "./httpClient";

// Возвращает массив изделий каталога: [{ id, title, description }].
export async function getProducts() {
  return requestJson("/api/products");
}
