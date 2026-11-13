// Общий HTTP-клиент.
//
// Выполняет fetch по относительному адресу и приводит результат к единому
// формату: разобранные данные (или null) при успехе, либо ошибка с полями
// message и status.
//
// Клиент не показывает alert, не меняет состояние React и не выполняет
// переходы между страницами. Повторные POST-запросы не выполняются.

const CONNECTION_ERROR_MESSAGE =
  "Не удалось связаться с сервером. Проверьте, запущен ли Java backend.";

function hasContentType(headers) {
  return Object.keys(headers || {}).some(
    (name) => name.toLowerCase() === "content-type"
  );
}

async function buildHttpError(response) {
  // Общее сообщение, если тело не является JSON или не содержит message.
  let message = `Не удалось выполнить запрос. Код ответа: ${response.status}.`;

  const contentType = response.headers.get("content-type") || "";
  if (contentType.includes("application/json")) {
    try {
      const data = await response.json();
      if (data && typeof data.message === "string" && data.message.trim() !== "") {
        message = data.message;
      }
    } catch {
      // Оставляем общее сообщение.
    }
  }

  const error = new Error(message);
  error.status = response.status;
  return error;
}

export async function requestJson(path, options = {}) {
  const headers = { ...(options.headers || {}) };

  // При отправке тела считаем данные JSON и добавляем заголовок,
  // если он не задан явно.
  if (options.body != null && !hasContentType(headers)) {
    headers["Content-Type"] = "application/json";
  }

  let response;
  try {
    response = await fetch(path, { ...options, headers });
  } catch {
    const error = new Error(CONNECTION_ERROR_MESSAGE);
    error.status = null;
    throw error;
  }

  if (!response.ok) {
    throw await buildHttpError(response);
  }

  const text = await response.text();
  if (text === "") {
    return null;
  }

  return JSON.parse(text);
}
