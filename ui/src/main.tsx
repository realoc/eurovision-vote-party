import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import "./index.css";
import {Theme} from "@radix-ui/themes";
import App from "./App.tsx";

const rootElement = document.getElementById("root");
if (!rootElement) {
  throw new Error("Root element not found");
}

createRoot(rootElement).render(
  <StrictMode>
    <Theme>
      <App />
    </Theme>
  </StrictMode>,
);
