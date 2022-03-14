import { useApolloClient } from "@apollo/client";
import { TabContext, TabPanel } from "@mui/lab";
import TabList from "@mui/lab/TabList";
import { Box, FormControl, FormControlLabel, FormHelperText, Switch, Tab, Typography } from "@mui/material";
import { SwitchBaseProps } from "@mui/material/internal/SwitchBase";
import { useState } from "react";

import { useActiveProject } from "lib/ActiveProjectContext";
import Link from "lib/components/Link";
import { ActiveProjectDocument, useUpdateInterceptSettingsMutation } from "lib/graphql/generated";

enum TabValue {
  Intercept = "intercept",
}

export default function Settings(): JSX.Element {
  const client = useApolloClient();
  const activeProject = useActiveProject();
  const [updateInterceptSettings, updateIntercepSettingsResult] = useUpdateInterceptSettingsMutation();

  const handleInterceptEnabled: SwitchBaseProps["onChange"] = (_, checked) => {
    updateInterceptSettings({
      variables: {
        input: {
          enabled: checked,
        },
      },
      onCompleted(data) {
        client.cache.updateQuery({ query: ActiveProjectDocument }, (cachedData) => ({
          activeProject: {
            ...cachedData.activeProject,
            settings: {
              intercept: data.updateInterceptSettings,
            },
          },
        }));
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
      <Typography variant="h6" sx={{ mb: 2 }}>
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
            </TabPanel>
          </TabContext>
        </>
      )}
    </Box>
  );
}
