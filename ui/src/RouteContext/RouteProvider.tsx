import { ReactNode, useEffect, useState } from "react";
import * as api from "../api";
import { RouteContext } from "./RouteContext";
import { Result } from "../result";

function nameFrom(uri: string): string {
	const parts = uri.substring(1).split('/');
	return parts[1] ?? '';
}

function errorToString(e: unknown): string {
	if (typeof e === 'string') {
		return e;
	} else if (e instanceof Error) {
		return e.message;
	}
	return 'An unknown error occurred';
}

export const RouteProvider = ({ children }: { children: ReactNode }) => {
	const [result, setResult] = useState<Result<api.Route>>(
		{ value: { name: '', url: '' }, error: '' },
	);

	const name = nameFrom(location.pathname);

	useEffect(() => {
		if (name === '') {
			setResult({ value: { name: '', url: '' }, error: '' });
			return;
		}

		api.getRoute(name)
			.then(route => setResult({ value: route, error: '' }))
			.catch(error => setResult({ value: { name, url: '' }, error: errorToString(error) }));
	}, [setResult, name]);

	const updateRoute = async (name: string, url: string) => {
		try {
			setResult({
				value: await api.postRoute(name, url),
				error: '',
			});
		} catch (e) {
			setResult({
				value: { name, url },
				error: errorToString(e),
			});
		}
	};

	const deleteRoute = async (name: string) => {
		try {
			await api.deleteRoute(name);
			setResult({ value: { name, url: '' }, error: '' });
		} catch (e) {
			setResult({
				value: { name, url: '' },
				error: errorToString(e),
			});
		}
	};

	return (
		<RouteContext.Provider value={{ result, updateRoute, deleteRoute }} >
			{children}
		</RouteContext.Provider >
	);
};