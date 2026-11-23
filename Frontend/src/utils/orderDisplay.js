const statusNames = {
  NEW: "Новый",
  READY_FOR_PLANNING: "Готов к планированию",
  PLANNED: "Запланирован",
  IN_PRODUCTION: "В производстве",
  COMPLETED: "Выполнен",
  CANCELLED: "Отменён",
};

export function getOrderStatus(order) {
  if (order.status === "NEW") {
    if (order.proposal?.status === "PENDING") return "Ожидает согласования даты";
    if (order.proposal?.status === "ACCEPTED") return "Дата согласована, план не утверждён";
    if (order.proposal?.status === "REJECTED") return "Дата отклонена";
  }
  return statusNames[order.status] ?? "Неизвестное состояние";
}

export function formatAgreedDate(date) {
  if (!date) return "Не согласована";
  const [year, month, day] = date.split("-");
  return `${day}.${month}.${year}`;
}
