import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig(({ mode }) => {
  // Загружаем все переменные окружения (без фильтра по префиксу VITE_),
  // чтобы прочитать JAVA_API_TARGET.
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [react()],

    server: {
      host: "0.0.0.0",
      port: 5173,
      strictPort: true,
      proxy: {
        // Перенаправляем запросы /api на Java backend.
        // Префикс /api сохраняется: Java-контроллеры используют его в адресах.
        "/api": {
          target: env.JAVA_API_TARGET || "http://localhost:8080",
          changeOrigin: true,
        },
      },
    },

    preview: {
      host: "0.0.0.0",
      port: 4173,
      strictPort: true,
    },
  };
});
