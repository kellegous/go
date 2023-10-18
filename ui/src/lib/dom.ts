type NativeEl = HTMLElement | SVGElement;

function toNative<T extends NativeEl>(el: El<T> | T): T {
	return el instanceof El ? el.el : el;
}

export class El<T extends NativeEl = HTMLElement> {
	constructor(
		public readonly el: T,
	) {
	}

	withAttrs(attrs: { [k: string]: any }): this {
		const { el } = this;
		for (const [k, v] of Object.entries(attrs)) {
			el.setAttribute(k, v);
		}
		return this;
	}

	withCSS(props: { [k: string]: any }): this {
		const { el } = this,
			{ style } = el;
		for (const [k, v] of Object.entries(props)) {
			style.setProperty(k, v, '');
		}
		return this;
	}

	withClass(c: string): this {
		this.el.classList.add(c);
		return this;
	}

	withText(t: string): this {
		this.el.textContent = t;
		return this;
	}

	appendTo<V extends NativeEl>(e: El<V> | V): this {
		toNative(e).appendChild(this.el);
		return this;
	}

	append<V extends NativeEl>(e: El<V> | V): this {
		this.el.appendChild(toNative(e));
		return this;
	}

	prependTo<V extends NativeEl>(e: El<V> | V): this {
		const n = toNative(e),
			{ el } = this;
		n.insertBefore(el, n.firstChild);
		return this;
	}

	prepend<V extends NativeEl>(e: El<V> | V): this {
		const { el } = this;
		el.insertBefore(toNative(e), el.firstChild);
		return this;
	}

	on(
		name: string,
		fn: EventListenerOrEventListenerObject,
		capture?: boolean | undefined,
	): this {
		this.el.addEventListener(name, fn, capture);
		return this;
	}

	static from<T extends NativeEl>(el: T): El<T> {
		return new this(el);
	}

	static of<T extends NativeEl>(name: string, ns: string | null = null): El<T> {
		const el = (ns === null)
			? document.createElement(name)
			: document.createAttributeNS(ns, name);
		return new this(el as T);
	}
}

export namespace html {
	export function of<T extends HTMLElement = HTMLElement>(name: string): El<T> {
		return El.of(name);
	}

	export function from<T extends HTMLElement = HTMLElement>(el: T): El<T> {
		return El.from(el);
	}

	export function q<T extends HTMLElement = HTMLElement>(s: string): El<T> {
		return El.from(document.querySelector(s)!);
	}
}

export namespace svg {
	export const NS = 'http://www.w3.org/2000/svg';

	export function of<T extends SVGElement = SVGElement>(name: string): El<T> {
		return El.of(name, NS);
	}

	export function from<T extends SVGElement = SVGElement>(el: T): El<T> {
		return El.from(el);
	}
}