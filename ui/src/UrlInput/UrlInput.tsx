import { useState } from 'react';
import css from './UrlInput.module.scss';

export interface UrlInputProps {
}

export const UrlInput = ({ }: UrlInputProps) => {
	const [url, setUrl] = useState('');

	const urlDidChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		setUrl(event.target.value);
	};

	const clearButtonClass = url.length > 0
		? `${css['clear-button']} ${css.visible}`
		: css['clear-button'];

	return (
		<>
			<div className={css.bar}>
				<div className={clearButtonClass}></div>
				<input id="url"
					className={css.url}
					type="text"
					placeholder="Enter the url to shorten"
					autoComplete="off"
					onChange={urlDidChange} />
			</div>
			<div id="cmp"></div>
		</>
	);
}