import { ReactNode, useEffect, useState } from "react";
import * as api from "../api";
import { RouteContext } from "./RouteContext";
import { Result } from "../result";

function nameFrom(uri: string): string {
	const parts = uri.substring(1).split('/');
	return parts[1] ?? '';
}

function errorToString(e: unknown): string {
	if (e instanceof api.ApiError) {
		return e.message;
	}
	return 'Oops! Something went sideways!';
}

export const RouteProvider = ({ children }: { children: ReactNode }) => {
	const [result, setResult] = useState<Result<api.Route>>(
		Result.of({ name: '', url: '' })
	);

	const name = nameFrom(location.pathname);

	useEffect(() => {
		if (name === '') {
			setResult(Result.of({ name: '', url: '' }));
			return;
		}

		Result.from(
			() => api.getRoute(name),
			{ name: name, url: '' },
			errorToString
		).then(setResult);
	}, [setResult, name]);

	const updateRoute = async (name: string, url: string) =>
		setResult(await Result.from(
			() => api.postRoute(name, url),
			{ name, url },
			errorToString
		));

	const deleteRoute = async (name: string) =>
		setResult(await Result.from(
			() => api.deleteRoute(name),
			{ name, url: '' },
			errorToString
		));

	return (
		<RouteContext.Provider value={{ result, updateRoute, deleteRoute }} >
			{children}
		</RouteContext.Provider >
	);
};