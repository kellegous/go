import { createContext } from "react";
import { Config } from "../api";

export const ConfigContext = createContext<Config | null>(null);
