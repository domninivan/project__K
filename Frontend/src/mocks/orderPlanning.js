const demoOperations = [
  { name: "Подготовка заготовок", machine: "CUT-001", machineName: "Отрезной станок" },
  { name: "Обработка", machine: "TURN-001", machineName: "Токарный станок" },
];

export function createMockPlan(order, today = new Date()) {
  const version = (order.plan?.version ?? 0) + 1;
  const date = new Date(today.getTime() + (6 + version) * 86_400_000);
  const finishDate = new Intl.DateTimeFormat("sv-SE", {
    timeZone: "Europe/Moscow", year: "numeric", month: "2-digit", day: "2-digit",
  }).format(date);
  function timeAt(minutes) {
    const hours = String(Math.floor(minutes / 60)).padStart(2, "0");
    return `${hours}:${String(minutes % 60).padStart(2, "0")}`;
  }

  // В демоверсии используем готовые окна, без расчёта загрузки станков.
  return {
    version,
    status: "PROPOSED",
    finishDate,
    operations: order.items.flatMap((item, index) =>
      demoOperations.map((operation, step) => {
        const start = 9 * 60 + index * 60 + step * 30;
        return {
          ...operation,
          id: `${item.id}:${step + 1}`,
          itemId: item.id,
          productTitle: item.productTitle,
          sequence: step + 1,
          quantity: item.quantity,
          start: timeAt(start),
          end: timeAt(start + 30),
        };
      }),
    ),
  };
}

export function getPlanningActions(order) {
  return {
    calculate: order.status === "NEW" &&
      (!order.plan || order.plan.status === "REJECTED"),
    propose: order.status === "NEW" && order.plan?.status === "PROPOSED" &&
      order.proposal?.version !== order.plan.version,
    respond: order.status === "NEW" && order.proposal?.status === "PENDING" &&
      order.proposal.version === order.plan?.version,
    approve: order.status === "NEW" && order.proposal?.status === "ACCEPTED" &&
      order.proposal.version === order.plan?.version && Boolean(order.agreedDate),
  };
}

export function updatePlanning(order, action, data = {}) {
  if (!getPlanningActions(order)[action]) return order;
  if (action !== "calculate" && data.version !== order.plan?.version) return order;

  switch (action) {
    case "calculate":
      return {
        ...order,
        plan: data.plan,
        planHistory: order.plan ? [...(order.planHistory ?? []), order.plan] : [],
      };
    case "propose":
      return {
        ...order,
        proposal: {
          version: order.plan.version,
          date: order.plan.finishDate,
          reason: data.reason,
          status: "PENDING",
        },
      };
    case "respond":
      if (!data.accepted) {
        return {
          ...order,
          plan: { ...order.plan, status: "REJECTED" },
          proposal: { ...order.proposal, status: "REJECTED", response: data.response },
        };
      }
      return {
        ...order,
        originalAgreedDate: order.originalAgreedDate ?? order.proposal.date,
        agreedDate: order.proposal.date,
        proposal: { ...order.proposal, status: "ACCEPTED" },
      };
    case "approve":
      return { ...order, status: "PLANNED", plan: { ...order.plan, status: "ACCEPTED" } };
    default:
      return order;
  }
}
