import { useContext } from "react";
import { ConfigContext } from "./ConfigContext";

export const useConfig = () => {
  const config = useContext(ConfigContext);
  if (!config) {
    throw new Error("useConfig must be used within a ConfigProvider");
  }
  return config;
};
