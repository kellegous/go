export type Styles = {
	button: string;
	icon: string;
	feedback: string;
	active: string;
}

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;