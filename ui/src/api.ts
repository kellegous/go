
interface RawRoute {
	name: string;
	url: string;
	source_host: string;
	time: string;
}

export interface Route {
	name: string;
	url: string;
	source_host?: string;
	time?: Date;
}

function toRoute(route: RawRoute): Route {
	return {
		name: route.name,
		url: route.url,
		source_host: route.source_host,
		time: new Date(route.time),
	};
}

interface RouteResponse {
	ok: boolean;
	error?: string;
	route?: RawRoute;
}

async function fromResponse(res: Response): Promise<Route | null> {
	if (res.status == 404) {
		return null;
	}

	const { ok, error, route } = await res.json() as RouteResponse;
	if (!ok || !route) {
		throw new ApiError(error ?? 'Oof. Something went sideways.');
	}

	return toRoute(route);
}

export class ApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'ApiError';
	}
}

export async function getRoute(name: string): Promise<Route> {
	const route = await fromResponse(await fetch(`/api/url/${name}`));
	return route ?? { name, url: "" };
}

export async function postRoute(name: string, url: string): Promise<Route> {
	const route = await fromResponse(await fetch(`/api/url/${name}`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({ url }),
	}));
	return route ?? { name, url };
}

export async function deleteRoute(name: string): Promise<Route> {
	const res = await fetch(`/api/url/${name}`, {
		method: 'DELETE',
	});
	if (!res.ok) {
		throw new ApiError('Failed to delete route');
	}
	return { name, url: '' };
}