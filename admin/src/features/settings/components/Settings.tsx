import { useApolloClient } from "@apollo/client";
import { TabContext, TabPanel } from "@mui/lab";
import TabList from "@mui/lab/TabList";
import {
  Box,
  Button,
  CircularProgress,
  FormControl,
  FormControlLabel,
  FormHelperText,
  Switch,
  Tab,
  TextField,
  Typography,
} from "@mui/material";
import { SwitchBaseProps } from "@mui/material/internal/SwitchBase";
import { useEffect, useState } from "react";

import { useActiveProject } from "lib/ActiveProjectContext";
import Link from "lib/components/Link";
import { ActiveProjectDocument, useUpdateInterceptSettingsMutation } from "lib/graphql/generated";
import { withoutTypename } from "lib/graphql/omitTypename";

enum TabValue {
  Intercept = "intercept",
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
    },
  });

  const [interceptReqFilter, setInterceptReqFilter] = useState("");

  useEffect(() => {
    setInterceptReqFilter(activeProject?.settings.intercept.requestFilter || "");
  }, [activeProject?.settings.intercept.requestFilter]);

  const handleInterceptEnabled: SwitchBaseProps["onChange"] = (e, checked) => {
    if (!activeProject) {
      e.preventDefault();
      return;
    }

    updateInterceptSettings({
      variables: {
        input: {
          ...withoutTypename(activeProject.settings.intercept),
          enabled: checked,
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

  const [tabValue, setTabValue] = useState(TabValue.Intercept);

  const tabSx = {
    textTransform: "none",
  };

  return (
    <Box p={4}>
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
              <FormControl>
                <FormControlLabel
                  control={
                    <Switch
                      disabled={updateIntercepSettingsResult.loading}
                      onChange={handleInterceptEnabled}
                      checked={activeProject.settings.intercept.enabled}
                    />
                  }
                  label="Enable proxy interception"
                  labelPlacement="start"
                  sx={{ display: "inline-block", m: 0 }}
                />
                <FormHelperText>
                  When enabled, incoming HTTP requests to the proxy are stalled for{" "}
                  <Link href="/proxy/intercept">manual review</Link>.
                </FormHelperText>
              </FormControl>
              <Typography variant="h6" sx={{ mt: 3 }}>
                Rules
              </Typography>
              <form>
                <FormControl sx={{ width: "50%" }}>
                  <TextField
                    label="Request filter"
                    placeholder={`method = "GET" OR url =~ "/foobar"`}
                    color="primary"
                    variant="outlined"
                    value={interceptReqFilter}
                    onChange={(e) => setInterceptReqFilter(e.target.value)}
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
                  />
                  <FormHelperText>
                    Filter expression to match incoming requests on. When set, only matching requests are intercepted.
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
            </TabPanel>
          </TabContext>
        </>
      )}
    </Box>
  );
}
