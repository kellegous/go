import { useEffect, useState } from "react";
import css from "./UrlInput.module.scss";
import { useRoute } from "../RouteContext";
import { CenterForm } from "../CenterForm";
import { Drawer, DrawerStyle } from "../Drawer";
import { Link } from "../Link";
import { useConfig } from "../../ConfigContext";

export const UrlInput = () => {
  const { host } = useConfig();

  const { result, updateRoute, deleteRoute } = useRoute();
  const { value: route, error } = result;

  const [url, setUrl] = useState(route.url);

  useEffect(() => {
    setUrl(route.url);
  }, [setUrl, route]);

  const urlDidChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setUrl(event.target.value);
  };

  const formDidSubmit = () => {
    updateRoute(route.name, url);
  };

  const clearButtonDidClick = () => {
    console.log("clearButtonDidClick");
    deleteRoute(route.name);
  };

  const routeHasUrl = route.url !== "";
  const hasError = error !== "";

  const clearButtonClass =
    url.length > 0
      ? `${css["clear-button"]} ${css.visible}`
      : css["clear-button"];

  return (
    <CenterForm onSubmit={formDidSubmit}>
      <div className={css.bar}>
        <div className={clearButtonClass} onClick={clearButtonDidClick}></div>
        <input
          id="url"
          className={css.url}
          type="text"
          placeholder="Enter the url to shorten"
          autoComplete="off"
          value={url}
          onChange={urlDidChange}
        />
      </div>
      <Drawer
        visible={routeHasUrl || hasError}
        style={hasError ? DrawerStyle.Error : DrawerStyle.Normal}
      >
        {hasError && <div>{String(error)}</div>}
        {routeHasUrl && !hasError && <Link name={route.name} host={host} />}
      </Drawer>
    </CenterForm>
  );
};
