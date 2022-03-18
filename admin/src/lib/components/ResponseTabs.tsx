import { TabContext, TabList, TabPanel } from "@mui/lab";
import { Box, Paper, Tab, Typography } from "@mui/material";
import React, { useState } from "react";

import { KeyValuePairTable, KeyValuePair, KeyValuePairTableProps } from "./KeyValuePair";

import Editor from "lib/components/Editor";

interface ResponseTabsProps {
  headers: KeyValuePair[];
  onHeaderChange?: KeyValuePairTableProps["onChange"];
  onHeaderDelete?: KeyValuePairTableProps["onDelete"];
  body?: string | null;
  onBodyChange?: (value: string) => void;
  hasResponse: boolean;
}

enum TabValue {
  Body = "body",
  Headers = "headers",
}

const reqNotSent = (
  <Paper variant="centered">
    <Typography>Response not received yet.</Typography>
  </Paper>
);

function ResponseTabs(props: ResponseTabsProps): JSX.Element {
  const { headers, onHeaderChange, onHeaderDelete, body, onBodyChange, hasResponse } = props;
  const [tabValue, setTabValue] = useState(TabValue.Body);

  const contentType = headers.find((header) => header.key.toLowerCase() === "content-type")?.value;

  const tabSx = {
    textTransform: "none",
  };

  const headersLength = onHeaderChange ? headers.length - 1 : headers.length;

  return (
    <Box height="100%" sx={{ display: "flex", flexDirection: "column" }}>
      <TabContext value={tabValue}>
        <Box sx={{ borderBottom: 1, borderColor: "divider", mb: 1 }}>
          <TabList onChange={(_, value) => setTabValue(value)}>
            <Tab
              value={TabValue.Body}
              label={"Body" + (body?.length ? ` (${body.length} byte` + (body.length > 1 ? "s" : "") + ")" : "")}
              sx={tabSx}
            />
            <Tab value={TabValue.Headers} label={"Headers" + (headersLength ? ` (${headersLength})` : "")} sx={tabSx} />
          </TabList>
        </Box>
        <Box flex="1 auto" overflow="hidden">
          <TabPanel value={TabValue.Body} sx={{ p: 0, height: "100%" }}>
            {hasResponse && (
              <Editor
                content={body || ""}
                onChange={(value) => {
                  onBodyChange && onBodyChange(value || "");
                }}
                monacoOptions={{ readOnly: onBodyChange === undefined }}
                contentType={contentType}
              />
            )}
            {!hasResponse && reqNotSent}
          </TabPanel>
          <TabPanel value={TabValue.Headers} sx={{ p: 0, height: "100%", overflow: "scroll" }}>
            {hasResponse && <KeyValuePairTable items={headers} onChange={onHeaderChange} onDelete={onHeaderDelete} />}
            {!hasResponse && reqNotSent}
          </TabPanel>
        </Box>
      </TabContext>
    </Box>
  );
}

export default ResponseTabs;
