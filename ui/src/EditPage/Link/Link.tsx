import { Copy } from "./Copy";
import css from "./Link.module.scss";

export interface LinkProps {
  name: string;
  host: string;
}

export const Link = ({ name, host }: LinkProps) => {
  const url = `${host}/${name}`;

  return (
    <div className={css.root}>
      <div className={css.link}>
        <a href={url}>{url}</a>
      </div>
      <Copy text={url} />
    </div>
  );
};
