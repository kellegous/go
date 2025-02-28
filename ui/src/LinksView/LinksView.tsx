import { useRoutes } from '../RoutesContext';
import css from './LinksView.module.scss';

export const LinksView = () => {
	// TODO(kellegous): handle the error.
	const { value } = useRoutes();

	return (
		<div className={css.links}>
			{value.map(({ name, url }) => (
				<div key={name}>
					<span>{name}</span>
					<span>{url}</span>
				</div>
			))}
		</div>
	);
};