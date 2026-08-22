import { Route } from "../../../api";
import { ShortLink } from "../ShortLink";
import { Controls } from "../Controls";
import css from "./LinkRow.module.scss";

const dateFormat = new Intl.DateTimeFormat("en-US", {
  month: "short",
  day: "numeric",
  year: "numeric",
});

export interface LinkRowProps {
  route: Route;
  short_url: string;
}

export const LinkRow = ({ short_url, route }: LinkRowProps) => {
  const { name, url, time } = route;
  return (
    <div className={css.linkrow}>
      <div className={css.upper}>
        <ShortLink url={short_url} />
        <div className={css.time}>{dateFormat.format(time)}</div>
      </div>
      <div className={css.lower}>
        <div className={css.url}>{url}</div>
        <div className={css.controls}>
          <Controls name={name} />
        </div>
      </div>
    </div>
  );
};
