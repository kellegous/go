import { useConfig } from "../../ConfigContext";
import { useRoutes } from "../RoutesContext";
import { LinkRow } from "./Link/LinkRow";
import css from "./LinksView.module.scss";

export const LinksView = () => {
  const { host } = useConfig();
  const { value } = useRoutes();

  return (
    <div className={css.links}>
      {value.map((route) => {
        return (
          <LinkRow
            key={route.name}
            route={route}
            short_url={`${host}/${route.name}`}
          />
        );
      })}
    </div>
  );
};
