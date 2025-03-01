export type Styles = {
	root: string;
	link: string;
}

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;