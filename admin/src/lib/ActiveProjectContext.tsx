import React, { createContext, useContext } from "react";

import { Project, useProjectsQuery } from "./graphql/generated";

const ActiveProjectContext = createContext<Project | null>(null);

interface Props {
  children?: React.ReactNode | undefined;
}

export function ActiveProjectProvider({ children }: Props): JSX.Element {
  const { data } = useProjectsQuery();
  const project = data?.projects.find((project) => project.isActive) || null;

  return <ActiveProjectContext.Provider value={project}>{children}</ActiveProjectContext.Provider>;
}

export function useActiveProject() {
  return useContext(ActiveProjectContext);
}
