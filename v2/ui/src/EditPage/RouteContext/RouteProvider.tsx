import { ReactNode, useEffect, useState } from "react";
import { RouteContext } from "./RouteContext";
import { Result } from "../../result";
import * as api from "../../api";
import { apiErrorToString, Route } from "../../api";

function nameFrom(uri: string): string {
  const parts = uri.substring(1).split("/");
  return parts[1] ?? "";
}

export const RouteProvider = ({ children }: { children: ReactNode }) => {
  const [result, setResult] = useState<Result<Route>>(
    Result.of({ name: "", url: "" })
  );

  const name = nameFrom(location.pathname);

  useEffect(() => {
    if (name === "") {
      setResult(Result.of({ name: "", url: "" }));
      return;
    }

    Result.from(
      () => api.getRoute(name),
      { name: name, url: "" },
      apiErrorToString
    ).then(setResult);
  }, [setResult, name]);

  const updateRoute = async (name: string, url: string) =>
    setResult(
      await Result.from(
        () => api.postRoute(name, url),
        { name, url },
        apiErrorToString
      )
    );

  const deleteRoute = async (name: string) =>
    setResult(
      await Result.from(
        () => api.deleteRoute(name),
        { name, url: "" },
        apiErrorToString
      )
    );

  return (
    <RouteContext.Provider value={{ result, updateRoute, deleteRoute }}>
      {children}
    </RouteContext.Provider>
  );
};
