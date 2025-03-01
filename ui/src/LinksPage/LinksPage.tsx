import { ConfigProvider } from "../ConfigContext";
import { LinksView } from "./LinksView";
import { RoutesProvider } from "./RoutesContext";

export const LinksPage = () => {
  return (
    <ConfigProvider>
      <RoutesProvider>
        <h1>Go Links</h1>
        <LinksView />
      </RoutesProvider>
    </ConfigProvider>
  );
};
