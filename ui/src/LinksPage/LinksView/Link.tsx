import css from "./Link.module.scss";

export interface LinkProps {
  name: string;
  short_url: string;
  url: string;
}

export const Link = ({ short_url, url, name }: LinkProps) => {
  return (
    <div className={css.link}>
      <div className={css.info}>
        <div className={css["short-url"]}>
          <a href={short_url}>{short_url}</a>
        </div>
        <div className={css.url}>{url}</div>
      </div>
      <div className={css.edit}>
        <a href={`/edit/${name}`} title="edit">
          <span className="material-symbols-outlined">settings</span>
        </a>
      </div>
    </div>
  );
};
