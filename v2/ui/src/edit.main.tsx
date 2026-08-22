import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { EditPage } from "./EditPage";

import "./edit.main.scss";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <EditPage />
  </StrictMode>
);
