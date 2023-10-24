import './edit_view.scss';

class Input {
	#label: HTMLLabelElement;

	private constructor(
		public readonly el: HTMLElement,
		label: HTMLLabelElement,
		public readonly value: HTMLInputElement,
	) {
		this.#label = label;
	}

	get label(): string {
		return this.#label.textContent ?? '';
	}

	set label(v: string) {
		this.#label.textContent = v;
	}

	appendTo(el: HTMLElement) {
		el.appendChild(this.el);
	}

	get classList(): DOMTokenList {
		return this.el.classList;
	}

	static create(id: string): Input {
		const el = document.createElement('div');
		el.classList.add('input');

		const label = document.createElement('label');
		label.classList.add('label');
		label.setAttribute('for', id);
		el.appendChild(label);

		const value = document.createElement('input');
		value.classList.add('value');
		value.setAttribute('id', id);
		el.appendChild(value);

		return new Input(el, label, value);
	}
}

export class EditView {
	private constructor(
		form: HTMLFormElement,
		name: Input,
		url: Input,
	) {
		console.log(form, name, url);
	}

	static createIn(el: HTMLElement): EditView {
		const form = document.createElement('form');
		form.autocomplete = 'off';
		form.classList.add('edit-view');

		const name = Input.create('edit-name');
		name.label = 'short name';
		name.classList.add('name');
		name.appendTo(form);

		const url = Input.create('edit-url');
		url.label = 'url';
		url.classList.add('url');
		url.appendTo(form);

		el.appendChild(form);

		return new this(form, name, url);
	}
}