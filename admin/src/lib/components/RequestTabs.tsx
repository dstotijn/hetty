import { TabContext, TabList, TabPanel } from "@mui/lab";
import { Box, Tab } from "@mui/material";
import React, { useState } from "react";

import { KeyValuePairTable, KeyValuePair, KeyValuePairTableProps } from "./KeyValuePair";

import Editor from "lib/components/Editor";

enum TabValue {
  QueryParams = "queryParams",
  Headers = "headers",
  Body = "body",
}

interface RequestTabsProps {
  queryParams: KeyValuePair[];
  headers: KeyValuePair[];
  onQueryParamChange?: KeyValuePairTableProps["onChange"];
  onQueryParamDelete?: KeyValuePairTableProps["onDelete"];
  onHeaderChange?: KeyValuePairTableProps["onChange"];
  onHeaderDelete?: KeyValuePairTableProps["onDelete"];
  body?: string | null;
  onBodyChange?: (value: string) => void;
}

function RequestTabs(props: RequestTabsProps): JSX.Element {
  const {
    queryParams,
    onQueryParamChange,
    onQueryParamDelete,
    headers,
    onHeaderChange,
    onHeaderDelete,
    body,
    onBodyChange,
  } = props;
  const [tabValue, setTabValue] = useState(TabValue.QueryParams);

  const tabSx = {
    textTransform: "none",
  };

  const queryParamsLength = onQueryParamChange ? queryParams.length - 1 : queryParams.length;
  const headersLength = onHeaderChange ? headers.length - 1 : headers.length;

  return (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <TabContext value={tabValue}>
        <Box sx={{ borderBottom: 1, borderColor: "divider", mb: 1 }}>
          <TabList onChange={(_, value) => setTabValue(value)}>
            <Tab
              value={TabValue.QueryParams}
              label={"Query Params" + (queryParamsLength ? ` (${queryParamsLength})` : "")}
              sx={tabSx}
            />
            <Tab value={TabValue.Headers} label={"Headers" + (headersLength ? ` (${headersLength})` : "")} sx={tabSx} />
            <Tab
              value={TabValue.Body}
              label={"Body" + (body?.length ? ` (${body.length} byte` + (body.length > 1 ? "s" : "") + ")" : "")}
              sx={tabSx}
            />
          </TabList>
        </Box>
        <Box flex="1 auto" overflow="scroll" height="100%">
          <TabPanel value={TabValue.QueryParams} sx={{ p: 0, height: "100%" }}>
            <Box>
              <KeyValuePairTable items={queryParams} onChange={onQueryParamChange} onDelete={onQueryParamDelete} />
            </Box>
          </TabPanel>
          <TabPanel value={TabValue.Headers} sx={{ p: 0, height: "100%" }}>
            <Box>
              <KeyValuePairTable items={headers} onChange={onHeaderChange} onDelete={onHeaderDelete} />
            </Box>
          </TabPanel>
          <TabPanel value={TabValue.Body} sx={{ p: 0, height: "100%" }}>
            <Editor
              content={body || ""}
              onChange={(value) => {
                onBodyChange && onBodyChange(value || "");
              }}
              monacoOptions={{ readOnly: onBodyChange === undefined }}
              contentType={headers.find(({ key }) => key.toLowerCase() === "content-type")?.value}
            />
          </TabPanel>
        </Box>
      </TabContext>
    </Box>
  );
}

export default RequestTabs;
