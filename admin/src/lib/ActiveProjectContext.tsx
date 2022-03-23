import React, { createContext, useContext } from "react";

import { Project, useActiveProjectQuery } from "./graphql/generated";

const ActiveProjectContext = createContext<Project | null>(null);

interface Props {
  children?: React.ReactNode | undefined;
}

export function ActiveProjectProvider({ children }: Props): JSX.Element {
  const { data } = useActiveProjectQuery();
  const project = data?.activeProject || null;

  return <ActiveProjectContext.Provider value={project}>{children}</ActiveProjectContext.Provider>;
}

export function useActiveProject() {
  return useContext(ActiveProjectContext);
}
