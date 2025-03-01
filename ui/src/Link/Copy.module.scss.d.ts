export type Styles = {
	button: string;
	icon: string;
	feedback: string;
}

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;