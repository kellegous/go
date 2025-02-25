import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { EditView } from "./EditView";

import "./edit.main.scss";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<EditView />
	</StrictMode>
);