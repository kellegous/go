import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { LinksPage } from "./LinksPage";

import "./links.main.scss";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<LinksPage />
	</StrictMode>
);
