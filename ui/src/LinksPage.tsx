import { LinksView } from './LinksView';
import { RoutesProvider } from './RoutesContext';

export const LinksPage = () => {
	return (
		<RoutesProvider>
			<h1>Go Links</h1>
			<LinksView />
		</RoutesProvider>
	);
};