import { ReactNode, useEffect, useMemo, useState } from "react";
import { RoutesContext } from "./RoutesContext";
import { Result } from "../../result";
import { Route, apiErrorToString, getRoutes } from "../../api";

export const RoutesProvider = ({ children }: { children: ReactNode }) => {
  const [result, setResult] = useState<Result<Route[]>>(Result.of([]));

  const res = useMemo(() => {
    return Result.from(() => getRoutes(), [], apiErrorToString);
  }, []);

  useEffect(() => {
    res.then(setResult);
  }, [res, setResult]);

  return (
    <RoutesContext.Provider value={result}>{children}</RoutesContext.Provider>
  );
};
