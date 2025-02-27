import css from './Link.module.scss';

export interface LinkProps {
	name: string;
}

export const Link = ({ name }: LinkProps) => {
	const url = `${location.origin}/${name}`;

	return (
		<div className={css.root}>
			<span className={css.link}>
				<a href={url}>{url}</a>
			</span>
			<span className={`${css.copy} material-symbols-outlined`}>
				content_copy
			</span>
		</div>
	);
};