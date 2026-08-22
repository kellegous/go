import { ReactNode, useEffect, useMemo, useState } from "react";
import { ConfigContext } from "./ConfigContext";
import { Config, getConfig } from "../api";

export const ConfigProvider = ({ children }: { children: ReactNode }) => {
  const [config, setConfig] = useState<Config>({ host: "" });

  const response = useMemo(() => getConfig(), []);

  useEffect(() => {
    response.then(setConfig);
  }, [response, setConfig]);

  return (
    <ConfigContext.Provider value={config}>{children}</ConfigContext.Provider>
  );
};
