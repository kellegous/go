import { useConfig } from "../../ConfigContext";
import { useRoutes } from "../RoutesContext";
import { Link } from "./Link";
import css from "./LinksView.module.scss";

export const LinksView = () => {
  const { host } = useConfig();
  const { value } = useRoutes();

  return (
    <div className={css.links}>
      {value.map(({ name, url }, i) => {
        return (
          <>
            {i !== 0 && <hr key={name + "-div"} />}
            <Link
              key={name}
              name={name}
              url={url}
              short_url={`${host}/${name}`}
            />
          </>
        );
      })}
    </div>
  );
};
