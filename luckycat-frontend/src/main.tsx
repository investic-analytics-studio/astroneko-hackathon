import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { QueryProvider } from "./providers/QueryProvider";
import { Toaster } from "./components/ui/sonner";
import "./config/i18n";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryProvider>
      <App />
      <Toaster position="top-center" />
    </QueryProvider>
  </StrictMode>
);
