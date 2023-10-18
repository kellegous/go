import './edit.scss';

import { html, El } from './lib/dom';

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
		console.log(formEl, copyEl, clearEl, urlEl);
	}

	static async load(): Promise<App> {
		const name = nameFrom(location.pathname);

		console.log(`name = "${name}"`);
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

App.load();