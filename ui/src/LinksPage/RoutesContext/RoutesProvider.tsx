import { ReactNode, useEffect, useMemo, useState } from "react";
import { RoutesContext } from "./RoutesContext";
import { Result } from "../../result";
import { Route, apiErrorToString, getAllRoutes } from "../../api";
import { restartableAsync } from "../../restartable";

export const RoutesProvider = ({ children }: { children: ReactNode }) => {
  const [result, setResult] = useState<Result<Route[]>>(Result.of([]));

  const res = useMemo(() => restartableAsync(getAllRoutes()), []);

  useEffect(() => {
    const fetchAllRoutes = async () => {
      const allRoutes = [];
      for await (const routes of res()) {
        allRoutes.push(...routes);
        setResult(Result.of(allRoutes));
      }
    };
    fetchAllRoutes().catch((e) =>
      setResult(Result.error([], apiErrorToString(e)))
    );
  }, [res, setResult]);

  return (
    <RoutesContext.Provider value={result}>{children}</RoutesContext.Provider>
  );
};
