import './edit.scss';

import { html, El } from './lib/dom';
import { EditView } from './lib/edit_view';

function nameFrom(uri: string): string {
	const parts = uri.substring(1).split('/');
	return parts[1] ?? '';
}

class App {
	private constructor(
		formEl: El<HTMLFormElement>,
		copyEl: El,
		clearEl: El,
		urlEl: El<HTMLInputElement>,
	) {
		formEl.on('submit', (e) => this.formDidSubmit(e));
	}

	private formDidSubmit(e: Event): void {
		e.preventDefault();
		console.log(e);
	}

	static async load(): Promise<App> {
		const name = nameFrom(location.pathname);

		if (name !== '') {
			const data = await fetch(`/api/url/${name}`);
			console.log(data);
		}

		return new App(
			html.q<HTMLFormElement>('form'),
			html.q('#cmp'),
			html.q('#cls'),
			html.q<HTMLInputElement>('#url'),
		);
	}
}

EditView.createIn(
	document.querySelector('#app')!
);
