import CloseIcon from "@mui/icons-material/Close";
import DeleteIcon from "@mui/icons-material/Delete";
import DescriptionIcon from "@mui/icons-material/Description";
import LaunchIcon from "@mui/icons-material/Launch";
import SettingsIcon from "@mui/icons-material/Settings";
import { Alert } from "@mui/lab";
import {
  Avatar,
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  ListItemSecondaryAction,
  ListItemText,
  Paper,
  Snackbar,
  Tooltip,
  Typography,
  useTheme,
} from "@mui/material";
import React, { useState } from "react";

import useOpenProjectMutation from "../hooks/useOpenProjectMutation";

import Link, { NextLinkComposed } from "lib/components/Link";
import {
  ProjectsQuery,
  useCloseProjectMutation,
  useDeleteProjectMutation,
  useProjectsQuery,
} from "lib/graphql/generated";

function ProjectList(): JSX.Element {
  const theme = useTheme();
  const projResult = useProjectsQuery({ fetchPolicy: "network-only" });
  const [openProject, openProjResult] = useOpenProjectMutation();
  const [closeProject, closeProjResult] = useCloseProjectMutation({
    errorPolicy: "all",
    onCompleted() {
      closeProjResult.client.resetStore();
    },
    update(cache) {
      cache.modify({
        fields: {
          activeProject() {
            return null;
          },
          projects(_, { DELETE }) {
            return DELETE;
          },
          httpRequestLogFilter(_, { DELETE }) {
            return DELETE;
          },
        },
      });
    },
  });
  const [deleteProject, deleteProjResult] = useDeleteProjectMutation({
    errorPolicy: "all",
    update(cache) {
      cache.modify({
        fields: {
          projects(_, { DELETE }) {
            return DELETE;
          },
        },
      });
      setDeleteDiagOpen(false);
      setDeleteNotifOpen(true);
    },
  });

  const [deleteProj, setDeleteProj] = useState<ProjectsQuery["projects"][number]>();
  const [deleteDiagOpen, setDeleteDiagOpen] = useState(false);
  const handleDeleteButtonClick = (project: ProjectsQuery["projects"][number]) => {
    setDeleteProj(project);
    setDeleteDiagOpen(true);
  };
  const handleDeleteConfirm = () => {
    if (deleteProj) {
      deleteProject({ variables: { id: deleteProj.id } });
    }
  };
  const handleDeleteCancel = () => {
    setDeleteDiagOpen(false);
  };

  const [deleteNotifOpen, setDeleteNotifOpen] = useState(false);
  const handleCloseDeleteNotif = (_: Event | React.SyntheticEvent, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }
    setDeleteNotifOpen(false);
  };

  return (
    <div>
      <Dialog open={deleteDiagOpen} onClose={handleDeleteCancel}>
        <DialogTitle>
          Delete project “<strong>{deleteProj?.name}</strong>”?
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Deleting a project permanently removes all its data from the database. This action is irreversible.
          </DialogContentText>
          {deleteProjResult.error && (
            <Alert severity="error">Error closing project: {deleteProjResult.error.message}</Alert>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel} autoFocus color="secondary" variant="contained">
            Cancel
          </Button>
          <Button
            sx={{
              color: "white",
              backgroundColor: "error.main",
              "&:hover": {
                backgroundColor: "error.dark",
              },
            }}
            onClick={handleDeleteConfirm}
            disabled={deleteProjResult.loading}
            variant="contained"
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={deleteNotifOpen}
        autoHideDuration={3000}
        onClose={handleCloseDeleteNotif}
        anchorOrigin={{ horizontal: "center", vertical: "bottom" }}
      >
        <Alert onClose={handleCloseDeleteNotif} severity="info">
          Project <strong>{deleteProj?.name}</strong> was deleted.
        </Alert>
      </Snackbar>

      <Box mb={3}>
        <Typography variant="h6">Manage projects</Typography>
      </Box>

      <Box mb={4}>
        {projResult.loading && <CircularProgress />}
        {projResult.error && <Alert severity="error">Error fetching projects: {projResult.error.message}</Alert>}
        {openProjResult.error && <Alert severity="error">Error opening project: {openProjResult.error.message}</Alert>}
        {closeProjResult.error && (
          <Alert severity="error">Error closing project: {closeProjResult.error.message}</Alert>
        )}
      </Box>

      {projResult.data && projResult.data.projects.length > 0 && (
        <Paper>
          <List>
            {projResult.data.projects.map((project) => (
              <ListItem key={project.id}>
                <ListItemAvatar>
                  <Avatar
                    sx={{
                      ...(project.isActive && {
                        color: theme.palette.secondary.dark,
                        backgroundColor: theme.palette.primary.main,
                      }),
                    }}
                  >
                    <DescriptionIcon />
                  </Avatar>
                </ListItemAvatar>
                <ListItemText>
                  {project.name} {project.isActive && <em>(Active)</em>}
                </ListItemText>
                <ListItemSecondaryAction>
                  <Tooltip title="Project settings">
                    <IconButton LinkComponent={Link} href="/settings" disabled={!project.isActive}>
                      <SettingsIcon />
                    </IconButton>
                  </Tooltip>
                  {project.isActive && (
                    <Tooltip title="Close project">
                      <IconButton onClick={() => closeProject()}>
                        <CloseIcon />
                      </IconButton>
                    </Tooltip>
                  )}
                  {!project.isActive && (
                    <Tooltip title="Open project">
                      <span>
                        <IconButton
                          disabled={openProjResult.loading || projResult.loading}
                          onClick={() =>
                            openProject({
                              variables: { id: project.id },
                            })
                          }
                        >
                          <LaunchIcon />
                        </IconButton>
                      </span>
                    </Tooltip>
                  )}
                  <Tooltip title="Delete project">
                    <span>
                      <IconButton onClick={() => handleDeleteButtonClick(project)} disabled={project.isActive}>
                        <DeleteIcon />
                      </IconButton>
                    </span>
                  </Tooltip>
                </ListItemSecondaryAction>
              </ListItem>
            ))}
          </List>
        </Paper>
      )}
      {projResult.data?.projects.length === 0 && (
        <Alert severity="info">There are no projects. Create one to get started.</Alert>
      )}
    </div>
  );
}

export default ProjectList;
