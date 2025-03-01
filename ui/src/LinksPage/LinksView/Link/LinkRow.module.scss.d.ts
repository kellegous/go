export type Styles = {
  linkrow: string;
  upper: string;
  lower: string;
  url: string;
  controls: string;
  time: string;
};

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;
