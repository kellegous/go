import { Copy } from "./Copy";
import css from "./Link.module.scss";

export interface LinkProps {
  name: string;
}

export const Link = ({ name }: LinkProps) => {
  const url = `${location.origin}/${name}`;

  return (
    <div className={css.root}>
      <div className={css.link}>
        <a href={url}>{url}</a>
      </div>
      <Copy text={url} />
    </div>
  );
};
