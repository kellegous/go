export type Styles = {
	bar: string;
	url: string;
	["clear-button"]: string;
	visible: string;
}

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;