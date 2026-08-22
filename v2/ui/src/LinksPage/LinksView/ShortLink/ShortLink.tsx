import css from "./ShortLink.module.scss";

export interface ShortLinkProps {
  url: string;
}

export const ShortLink = ({ url }: ShortLinkProps) => {
  return (
    <div className={css.shortlink}>
      <a href={url}>{stripScheme(url)}</a>
    </div>
  );
};

function stripScheme(url: string): string {
  return url.replace(/^https?:\/\//, "");
}
