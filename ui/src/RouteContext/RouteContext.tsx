import { Route } from "../api";
import { createContext } from "react";

export interface RouteInfo {
	route: Route;
	error: string;
}

export interface RouteContextState {
	info: RouteInfo;
	updateRoute: (name: string, url: string) => void;
	deleteRoute: (name: string) => void;
}

export const RouteContext = createContext<RouteContextState | null>(null);