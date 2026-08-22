import { useContext } from "react";
import { RouteContext } from "./RouteContext";

export const useRoute = () => {
	const context = useContext(RouteContext);
	if (!context) {
		throw new Error("useRoute must be used within a RouteProvider");
	}
	return context;
}