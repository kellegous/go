import { ConfigProvider } from "../ConfigContext";
import { RouteProvider } from "./RouteContext";
import { UrlInput } from "./UrlInput";

export const EditPage = () => {
  return (
    <ConfigProvider>
      <RouteProvider>
        <UrlInput />
      </RouteProvider>
    </ConfigProvider>
  );
};
