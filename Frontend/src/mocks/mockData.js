export const products = [
  {
    id: 1,
    title: "Вал приводной",
    description: "Стальной вал для передачи крутящего момента от двигателя к редуктору",
  },
  {
    id: 2,
    title: "Втулка",
    description: "Бронзовая втулка для подшипникового узла",
  },
  {
    id: 3,
    title: "Шестерня",
    description: "Цилиндрическая прямозубая шестерня, 32 зуба",
  },
  {
    id: 4,
    title: "Кронштейн",
    description: "Крепёжный кронштейн для установки узла на раму",
  },
];

export const orders = [
  {
    id: 1,
    dueDate: "2026-11-20",
    comment: "Срочно для сборки редуктора",
    priority: "HIGH",
    status: "NEW",
    createdAt: "2026-11-02T10:15:00",
    items: [
      { id: 1, productId: 1, productTitle: "Вал приводной", quantity: 10 },
      { id: 2, productId: 3, productTitle: "Шестерня", quantity: 20 },
      { id: 3, productId: 2, productTitle: "Втулка", quantity: 40 },
    ],
  },
  {
    id: 2,
    dueDate: "2026-12-05",
    comment: "",
    priority: "MEDIUM",
    status: "PLANNED",
    createdAt: "2026-11-04T14:40:00",
    items: [
      { id: 4, productId: 4, productTitle: "Кронштейн", quantity: 15 },
      { id: 5, productId: 2, productTitle: "Втулка", quantity: 30 },
    ],
  },
];
