import { createRef, useEffect } from 'react';
import css from './CenterForm.module.scss';

export interface CenterFormProps {
	children: React.ReactNode;
}

export const CenterForm = ({ children }: CenterFormProps) => {
	const ref = createRef<HTMLFormElement>();

	useEffect(() => {
		const { current } = ref;
		if (!current) {
			return;
		}

		const didResize = () => {
			const { height } = current.getBoundingClientRect();
			current.style.marginTop = (window.innerHeight / 3 - height / 2) + 'px';
		};

		window.addEventListener('resize', didResize);
		didResize();

		return () => window.removeEventListener('resize', didResize);
	}, [ref]);

	return (
		<form className={css.form} ref={ref}>
			{children}
		</form>
	);
}