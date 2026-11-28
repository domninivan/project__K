// Централизованный выбор режима данных.
//
// VITE_DATA_MODE:
//   "api"  — данные загружаются из Java backend через HTTP-клиент;
//   иначе  — демонстрационный режим на mock-данных (значение по умолчанию).
//
// Значение по умолчанию — "mock": без .env приложение не выполняет HTTP-запросы.

const mode = (import.meta.env.VITE_DATA_MODE || "mock").toLowerCase();

export const isApi = mode === "api";
export const isMock = !isApi;
export const dataMode = mode;
