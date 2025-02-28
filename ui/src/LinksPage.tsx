import { LinksView } from './LinksView';
import { RoutesProvider } from './RoutesContext';

export const LinksPage = () => {
	return (
		<RoutesProvider>
			<LinksView />
		</RoutesProvider>
	);
};