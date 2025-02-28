import { Route } from "../api";
import { createContext } from "react";
import { Result } from "../result";

export interface RouteContextState {
	result: Result<Route>;
	updateRoute: (name: string, url: string) => void;
	deleteRoute: (name: string) => void;
}

export const RouteContext = createContext<RouteContextState | null>(null);