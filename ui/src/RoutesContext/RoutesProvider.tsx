import { ReactNode, useEffect, useState } from "react";
import { RoutesContext } from "./RoutesContext";
import { Result } from "../result";
import { Route, apiErrorToString, getRoutes } from "../api";

export const RoutesProvider = ({ children }: { children: ReactNode }) => {
	const [result, setResult] = useState<Result<Route[]>>(Result.of([]));

	useEffect(() => {
		Result.from(
			() => getRoutes(),
			[],
			apiErrorToString,
		).then(setResult);
	}, [setResult]);

	return (
		<RoutesContext.Provider value={result}>
			{children}
		</RoutesContext.Provider>
	)
};