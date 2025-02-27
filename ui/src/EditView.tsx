import { RouteProvider } from "./RouteContext"
import { UrlInput } from "./UrlInput"

export const EditView = () => {
	return (
		<RouteProvider>
			<UrlInput />
		</RouteProvider>
	);
};