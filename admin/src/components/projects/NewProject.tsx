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

const OPEN_PROJECT = gql`
  mutation OpenProject($name: String!) {
    openProject(name: $name) {
      name
      isActive
    }
  }
`;

function NewProject(): JSX.Element {
  const classes = useStyles();
  const [input, setInput] = useState(null);

  const [openProject, { error, loading }] = useMutation(OPEN_PROJECT, {
    onError: () => {},
    onCompleted() {
      input.value = "";
    },
    update(cache, { data: { openProject } }) {
      cache.modify({
        fields: {
          activeProject() {
            const activeProjRef = cache.writeFragment({
              id: openProject.name,
              data: openProject,
              fragment: gql`
                fragment ActiveProject on Project {
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
              id: openProject.name,
              data: openProject,
              fragment: gql`
                fragment OpenProject on Project {
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

  const handleNewProjectForm = (e: React.SyntheticEvent) => {
    e.preventDefault();
    openProject({ variables: { name: input.value } });
  };

  return (
    <div>
      <Box mb={3}>
        <Typography variant="h6">New project</Typography>
      </Box>
      <form onSubmit={handleNewProjectForm} autoComplete="off">
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
          error={Boolean(error)}
          helperText={error && error.message}
        />
        <Button
          className={classes.button}
          type="submit"
          variant="contained"
          color="secondary"
          size="large"
          disabled={loading}
          startIcon={loading ? <CircularProgress size={22} /> : <AddIcon />}
        >
          Create & open project
        </Button>
      </form>
    </div>
  );
}

export default NewProject;
