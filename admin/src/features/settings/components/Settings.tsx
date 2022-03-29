import { useApolloClient } from "@apollo/client";
import { TabContext, TabPanel } from "@mui/lab";
import TabList from "@mui/lab/TabList";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  FormControl,
  FormControlLabel,
  FormHelperText,
  Snackbar,
  Switch,
  Tab,
  TextField,
  TextFieldProps,
  Typography,
} from "@mui/material";
import MaterialLink from "@mui/material/Link";
import { SwitchBaseProps } from "@mui/material/internal/SwitchBase";
import { useEffect, useState } from "react";

import { useActiveProject } from "lib/ActiveProjectContext";
import Link from "lib/components/Link";
import { ActiveProjectDocument, useUpdateInterceptSettingsMutation } from "lib/graphql/generated";
import { withoutTypename } from "lib/graphql/omitTypename";

enum TabValue {
  Intercept = "intercept",
}

function FilterTextField(props: TextFieldProps): JSX.Element {
  return (
    <TextField
      color="primary"
      variant="outlined"
      InputProps={{
        sx: { fontFamily: "'JetBrains Mono', monospace" },
        autoCorrect: "false",
        spellCheck: "false",
      }}
      InputLabelProps={{
        shrink: true,
      }}
      margin="normal"
      sx={{ mr: 1 }}
      {...props}
    />
  );
}

export default function Settings(): JSX.Element {
  const client = useApolloClient();
  const activeProject = useActiveProject();
  const [updateInterceptSettings, updateIntercepSettingsResult] = useUpdateInterceptSettingsMutation({
    onCompleted(data) {
      client.cache.updateQuery({ query: ActiveProjectDocument }, (cachedData) => ({
        activeProject: {
          ...cachedData.activeProject,
          settings: {
            ...cachedData.activeProject.settings,
            intercept: data.updateInterceptSettings,
          },
        },
      }));

      setInterceptReqFilter(data.updateInterceptSettings.requestFilter || "");
      setInterceptResFilter(data.updateInterceptSettings.responseFilter || "");
      setSettingsUpdatedOpen(true);
    },
  });

  const [interceptReqFilter, setInterceptReqFilter] = useState("");
  const [interceptResFilter, setInterceptResFilter] = useState("");

  useEffect(() => {
    setInterceptReqFilter(activeProject?.settings.intercept.requestFilter || "");
  }, [activeProject?.settings.intercept.requestFilter]);

  useEffect(() => {
    setInterceptResFilter(activeProject?.settings.intercept.responseFilter || "");
  }, [activeProject?.settings.intercept.responseFilter]);

  const handleReqInterceptEnabled: SwitchBaseProps["onChange"] = (e, checked) => {
    if (!activeProject) {
      e.preventDefault();
      return;
    }

    updateInterceptSettings({
      variables: {
        input: {
          ...withoutTypename(activeProject.settings.intercept),
          requestsEnabled: checked,
        },
      },
    });
  };

  const handleResInterceptEnabled: SwitchBaseProps["onChange"] = (e, checked) => {
    if (!activeProject) {
      e.preventDefault();
      return;
    }

    updateInterceptSettings({
      variables: {
        input: {
          ...withoutTypename(activeProject.settings.intercept),
          responsesEnabled: checked,
        },
      },
    });
  };

  const handleInterceptReqFilter = () => {
    if (!activeProject) {
      return;
    }

    updateInterceptSettings({
      variables: {
        input: {
          ...withoutTypename(activeProject.settings.intercept),
          requestFilter: interceptReqFilter,
        },
      },
    });
  };

  const handleInterceptResFilter = () => {
    if (!activeProject) {
      return;
    }

    updateInterceptSettings({
      variables: {
        input: {
          ...withoutTypename(activeProject.settings.intercept),
          responseFilter: interceptResFilter,
        },
      },
    });
  };

  const [tabValue, setTabValue] = useState(TabValue.Intercept);
  const [settingsUpdatedOpen, setSettingsUpdatedOpen] = useState(false);

  const handleSettingsUpdatedClose = (_: Event | React.SyntheticEvent, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }

    setSettingsUpdatedOpen(false);
  };

  const tabSx = {
    textTransform: "none",
  };

  return (
    <Box p={4}>
      <Snackbar open={settingsUpdatedOpen} autoHideDuration={3000} onClose={handleSettingsUpdatedClose}>
        <Alert onClose={handleSettingsUpdatedClose} severity="info">
          Intercept settings have been updated.
        </Alert>
      </Snackbar>

      <Typography variant="h4" sx={{ mb: 2 }}>
        Settings
      </Typography>
      <Typography paragraph sx={{ mb: 4 }}>
        Settings allow you to tweak the behaviour of Hetty’s features.
      </Typography>
      <Typography variant="h5" sx={{ mb: 2 }}>
        Project settings
      </Typography>
      {!activeProject && (
        <Typography paragraph>
          There is no project active. To configure project settings, first <Link href="/projects">open a project</Link>.
        </Typography>
      )}
      {activeProject && (
        <>
          <TabContext value={tabValue}>
            <TabList onChange={(_, value) => setTabValue(value)} sx={{ borderBottom: 1, borderColor: "divider" }}>
              <Tab value={TabValue.Intercept} label="Intercept" sx={tabSx} />
            </TabList>

            <TabPanel value={TabValue.Intercept} sx={{ px: 0 }}>
              <Typography variant="h6" sx={{ mt: 3, mb: 1 }}>
                Requests
              </Typography>
              <FormControl sx={{ mb: 2 }}>
                <FormControlLabel
                  control={
                    <Switch
                      disabled={updateIntercepSettingsResult.loading}
                      onChange={handleReqInterceptEnabled}
                      checked={activeProject.settings.intercept.requestsEnabled}
                    />
                  }
                  label="Enable request interception"
                  labelPlacement="start"
                  sx={{ display: "inline-block", m: 0 }}
                />
                <FormHelperText>
                  When enabled, incoming HTTP requests to the proxy are stalled for{" "}
                  <Link href="/proxy/intercept">manual review</Link>.
                </FormHelperText>
              </FormControl>
              <form>
                <FormControl sx={{ width: "50%" }}>
                  <FilterTextField
                    label="Request filter"
                    placeholder={`Example: method = "GET" OR url =~ "/foobar"`}
                    value={interceptReqFilter}
                    onChange={(e) => setInterceptReqFilter(e.target.value)}
                  />
                  <FormHelperText>
                    Filter expression to match incoming requests on. When set, only matching requests are intercepted.{" "}
                    <MaterialLink
                      href="https://hetty.xyz/docs/guides/intercept?utm_source=hettyapp#request-filter"
                      target="_blank"
                    >
                      Read docs.
                    </MaterialLink>
                  </FormHelperText>
                </FormControl>
                <Button
                  type="submit"
                  variant="text"
                  color="primary"
                  size="large"
                  sx={{
                    mt: 2,
                    py: 1.8,
                  }}
                  onClick={handleInterceptReqFilter}
                  disabled={updateIntercepSettingsResult.loading}
                  startIcon={updateIntercepSettingsResult.loading ? <CircularProgress size={22} /> : undefined}
                >
                  Update
                </Button>
              </form>
              <Typography variant="h6" sx={{ mt: 3 }}>
                Responses
              </Typography>
              <FormControl sx={{ mb: 2 }}>
                <FormControlLabel
                  control={
                    <Switch
                      disabled={updateIntercepSettingsResult.loading}
                      onChange={handleResInterceptEnabled}
                      checked={activeProject.settings.intercept.responsesEnabled}
                    />
                  }
                  label="Enable response interception"
                  labelPlacement="start"
                  sx={{ display: "inline-block", m: 0 }}
                />
                <FormHelperText>
                  When enabled, HTTP responses received by the proxy are stalled for{" "}
                  <Link href="/proxy/intercept">manual review</Link>.
                </FormHelperText>
              </FormControl>
              <form>
                <FormControl sx={{ width: "50%" }}>
                  <FilterTextField
                    label="Response filter"
                    placeholder={`Example: statusCode =~ "^2" OR body =~ "foobar"`}
                    value={interceptResFilter}
                    onChange={(e) => setInterceptResFilter(e.target.value)}
                  />
                  <FormHelperText>
                    Filter expression to match received responses on. When set, only matching responses are intercepted.{" "}
                    <MaterialLink
                      href="https://hetty.xyz/docs/guides/intercept/?utm_source=hettyapp#response-filter"
                      target="_blank"
                    >
                      Read docs.
                    </MaterialLink>
                  </FormHelperText>
                </FormControl>
                <Button
                  type="submit"
                  variant="text"
                  color="primary"
                  size="large"
                  sx={{
                    mt: 2,
                    py: 1.8,
                  }}
                  onClick={handleInterceptResFilter}
                  disabled={updateIntercepSettingsResult.loading}
                  startIcon={updateIntercepSettingsResult.loading ? <CircularProgress size={22} /> : undefined}
                >
                  Update
                </Button>
              </form>
            </TabPanel>
          </TabContext>
        </>
      )}
    </Box>
  );
}
