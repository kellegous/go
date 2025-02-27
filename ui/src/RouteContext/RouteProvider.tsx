import { ReactNode, useEffect, useState } from "react";
import * as api from "../api";
import { RouteContext, RouteInfo } from "./RouteContext";

function nameFrom(uri: string): string {
	const parts = uri.substring(1).split('/');
	return parts[1] ?? '';
}

export const RouteProvider = ({ children }: { children: ReactNode }) => {
	const [info, setRoute] = useState<RouteInfo>({
		route: { name: '', url: '' },
		error: '',
	});

	const name = nameFrom(location.pathname);

	useEffect(() => {
		if (name === '') {
			setRoute({
				route: { name: '', url: '' },
				error: ''
			});
			return;
		}

		api.getRoute(name)
			.then(route => setRoute({ route, error: '' }))
			.catch((error) => console.error(error));
	}, [setRoute, name]);

	const updateRoute = async (name: string, url: string) => {
		try {
			setRoute({ route: await api.postRoute(name, url), error: '' });
		} catch (e) {
			const message = (e instanceof Error) ? e.message : 'Failed to update route';
			setRoute({ route: { name, url }, error: message });
		}
	};

	const deleteRoute = async (name: string) => {
		try {
			await api.deleteRoute(name);
			setRoute({ route: { name, url: '' }, error: '' });
		} catch (e) {
			const message = (e instanceof Error) ? e.message : 'Failed to update route';
			setRoute({ route: { name, url: '' }, error: message });
		}
	};

	return (
		<RouteContext.Provider value={{ info, updateRoute, deleteRoute }} >
			{children}
		</RouteContext.Provider >
	);
};