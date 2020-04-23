namespace dom {
	export var q = (s: string) => {
		return <HTMLElement>document.querySelector(s);
	};

	export var qa = (s: string) => {
		return document.querySelectorAll(s);
	};

	export var c = (n: string) => {
		return document.createElement(n);
	};

	export var css = (el: HTMLElement, p: string, v: any) => {
		el.style.setProperty(p, v, '');
	};
}
