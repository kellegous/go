import { createContext } from 'react';
import { Result } from '../result';
import { Route } from '../api';

export const RoutesContext = createContext<Result<Route[]> | null>(
	null
);