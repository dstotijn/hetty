import { gql, useMutation } from "@apollo/client";
import {
  Box,
  Button,
  CircularProgress,
  createStyles,
  makeStyles,
  TextField,
  Theme,
  Typography,
} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import React, { useState } from "react";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    projectName: {
      marginTop: -6,
      marginRight: theme.spacing(2),
    },
    button: {
      marginRight: theme.spacing(2),
    },
  })
);

const CREATE_PROJECT = gql`
  mutation CreateProject($name: String!) {
    createProject(name: $name) {
      id
      name
    }
  }
`;

const OPEN_PROJECT = gql`
  mutation OpenProject($id: ID!) {
    openProject(id: $id) {
      id
      name
      isActive
    }
  }
`;

function NewProject(): JSX.Element {
  const classes = useStyles();
  const [input, setInput] = useState(null);

  const [createProject, { error: createProjErr, loading: createProjLoading }] = useMutation(CREATE_PROJECT, {
    onError: () => { },
    onCompleted(data) {
      input.value = "";
      openProject({ variables: { id: data.createProject.id } });
    },
  });
  const [openProject, { error: openProjErr, loading: openProjLoading }] = useMutation(OPEN_PROJECT, {
    onError: () => { },
    onCompleted() {
      input.value = "";
    },
    update(cache, { data: { openProject } }) {
      cache.modify({
        fields: {
          activeProject() {
            const activeProjRef = cache.writeFragment({
              id: openProject.id,
              data: openProject,
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
              id: openProject.id,
              data: openProject,
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
        },
      });
    },
  });

  const handleCreateAndOpenProjectForm = (e: React.SyntheticEvent) => {
    e.preventDefault();
    createProject({ variables: { name: input.value } });
  };

  return (
    <div>
      <Box mb={3}>
        <Typography variant="h6">New project</Typography>
      </Box>
      <form onSubmit={handleCreateAndOpenProjectForm} autoComplete="off">
        <TextField
          className={classes.projectName}
          color="secondary"
          inputProps={{
            id: "projectName",
            ref: (node) => {
              setInput(node);
            },
          }}
          label="Project name"
          placeholder="Project name…"
          error={Boolean(createProjErr || openProjErr)}
          helperText={createProjErr && createProjErr.message || openProjErr && openProjErr.message}
        />
        <Button
          className={classes.button}
          type="submit"
          variant="contained"
          color="secondary"
          size="large"
          disabled={createProjLoading || openProjLoading}
          startIcon={createProjLoading || openProjLoading ? <CircularProgress size={22} /> : <AddIcon />}
        >
          Create & open project
        </Button>
      </form>
    </div>
  );
}

export default NewProject;
