export type Styles = {
  link: string;
  info: string;
  edit: string;
  ["short-url"]: string;
  url: string;
};

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;
