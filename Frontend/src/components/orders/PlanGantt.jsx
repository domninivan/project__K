function minutesOfDay(time) {
  const [hours, minutes] = time.split(":").map(Number);
  return hours * 60 + minutes;
}

export default function PlanGantt({ plan }) {
  if (!plan.operations.length) return null;
  const start = Math.min(...plan.operations.map((operation) => minutesOfDay(operation.start)));
  const end = Math.max(...plan.operations.map((operation) => minutesOfDay(operation.end)));
  const span = end - start;
  if (span <= 0) return null;
  return (
    <div className="table-wrap">
      <div className="plan-gantt">
        <p>График операций · {plan.finishDate.split("-").reverse().join(".")} (МСК)</p>
        {plan.operations.map((operation) => {
          const left = (minutesOfDay(operation.start) - start) / span * 100;
          const width = (minutesOfDay(operation.end) - minutesOfDay(operation.start)) / span * 100;
          return (
            <div className="gantt-row" key={operation.id}>
              <span>{operation.machine} · {operation.productTitle}</span>
              <div className="gantt-track">
                <span className="gantt-bar" style={{ left: `${left}%`, width: `${width}%` }} />
              </div>
              <span>{operation.start}–{operation.end}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
