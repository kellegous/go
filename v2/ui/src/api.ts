interface RawRoute {
  name: string;
  url: string;
  source_host: string;
  time: string;
}

export interface Route {
  name: string;
  url: string;
  time?: Date;
}

export interface Config {
  host: string;
}

function toRoute(route: RawRoute): Route {
  const { name, url, time } = route;
  return { name, url, time: new Date(time) };
}

interface RouteResponse {
  ok: boolean;
  error?: string;
  route?: RawRoute;
}

interface RoutesResponse {
  ok: boolean;
  error?: string;
  routes?: RawRoute[];
  next: string;
}

async function fromResponse<T extends { ok: boolean; error?: string }, V>(
  res: Response,
  getValue: (json: T) => V | null
): Promise<V | null> {
  if (res.status == 404) {
    return null;
  }

  const data = (await res.json()) as T;
  const { ok, error } = data;
  const value = getValue(data);

  if (!ok || !value) {
    throw new ApiError(error ?? "Oof. Something went sideways.");
  }

  return value;
}

export class ApiError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ApiError";
  }
}

export async function getRoute(name: string): Promise<Route> {
  const route = await fromResponse(
    await fetch(`/api/url/${name}`),
    (data: RouteResponse) => (data.route ? toRoute(data.route) : null)
  );
  return route ?? { name, url: "" };
}

export async function getConfig(): Promise<Config> {
  const { host } = await fetch("/api/config").then((res) => res.json());
  return host === "" ? { host: location.origin } : { host };
}

export async function getRoutes(
  next: string,
  limit: number = 1000
): Promise<[Route[], string]> {
  const value = await fromResponse(
    await fetch(`/api/urls/?cursor=${next}&limit=${limit}`),
    (data: RoutesResponse) =>
      [data.routes?.map(toRoute) ?? [], data.next] as [Route[], string]
  );
  return value ?? [[], next];
}

export async function* getAllRoutes(
  pageSize: number = 1000
): AsyncGenerator<Route[]> {
  let cursor = "";
  do {
    const [routes, next] = await getRoutes(cursor, pageSize);
    yield routes;
    cursor = next;
  } while (cursor !== "");
}

export async function postRoute(name: string, url: string): Promise<Route> {
  const route = await fromResponse(
    await fetch(`/api/url/${name}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ url }),
    }),
    (data: RouteResponse) => (data.route ? toRoute(data.route) : null)
  );
  return route ?? { name, url };
}

export async function deleteRoute(name: string): Promise<Route> {
  const route = await fromResponse(
    await fetch(`/api/url/${name}`, {
      method: "DELETE",
    }),
    () => ({ name, url: "" })
  );
  return route ?? { name, url: "" };
}

export function apiErrorToString(e: unknown): string {
  if (e instanceof ApiError) {
    return e.message;
  }
  return "Oops! Something went sideways!";
}
