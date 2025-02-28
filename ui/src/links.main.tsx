import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { LinksView } from "./LinksView";

import "./edit.main.scss";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<LinksView />
	</StrictMode>
);
