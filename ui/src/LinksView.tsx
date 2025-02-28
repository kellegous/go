import { useRoutes } from './RoutesContext';

export const LinksView = () => {
	const { value, error } = useRoutes();
	console.log(value, error);

	return (
		<div>{value.length}</div>
	);
};