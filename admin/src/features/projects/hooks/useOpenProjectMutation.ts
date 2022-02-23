import { gql } from "@apollo/client";

import { useOpenProjectMutation as _useOpenProjectMutation } from "lib/graphql/generated";

export default function useOpenProjectMutation() {
  return _useOpenProjectMutation({
    errorPolicy: "all",
    update(cache, { data }) {
      cache.modify({
        fields: {
          activeProject() {
            const activeProjRef = cache.writeFragment({
              data: data?.openProject,
              fragment: gql`
                fragment ActiveProject on Project {
                  id
                  name
                  isActive
                  type
                }
              `,
            });
            return activeProjRef;
          },
          projects(_, { DELETE }) {
            cache.writeFragment({
              id: data?.openProject?.id,
              data: data?.openProject,
              fragment: gql`
                fragment OpenProject on Project {
                  id
                  name
                  isActive
                  type
                }
              `,
            });
            return DELETE;
          },
          httpRequestLogFilter(_, { DELETE }) {
            return DELETE;
          },
        },
      });
    },
  });
}
