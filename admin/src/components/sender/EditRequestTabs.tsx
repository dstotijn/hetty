import { TabContext, TabList, TabPanel } from "@mui/lab";
import { Box, Tab } from "@mui/material";
import { useState } from "react";
import Editor from "../common/Editor";

import KeyValuePairTable, { KeyValuePair, KeyValuePairTableProps } from "./KeyValuePair";

enum TabValue {
  QueryParams = "queryParams",
  Headers = "headers",
  Body = "body",
}

export type EditRequestTabsProps = {
  queryParams: KeyValuePair[];
  headers: KeyValuePair[];
  onQueryParamChange: KeyValuePairTableProps["onChange"];
  onQueryParamDelete: KeyValuePairTableProps["onDelete"];
  onHeaderChange: KeyValuePairTableProps["onChange"];
  onHeaderDelete: KeyValuePairTableProps["onDelete"];
  body: string;
  onBodyChange: (value: string) => void;
};

function EditRequestTabs(props: EditRequestTabsProps): JSX.Element {
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

  return (
    <Box height="100%" sx={{ display: "flex", flexDirection: "column" }}>
      <TabContext value={tabValue}>
        <Box sx={{ borderBottom: 1, borderColor: "divider", mb: 2 }}>
          <TabList onChange={(_, value) => setTabValue(value)}>
            <Tab
              value={TabValue.QueryParams}
              label={"Query Params" + (queryParams.length - 1 ? ` (${queryParams.length - 1})` : "")}
              sx={tabSx}
            />
            <Tab
              value={TabValue.Headers}
              label={"Headers" + (headers.length - 1 ? ` (${headers.length - 1})` : "")}
              sx={tabSx}
            />
            <Tab
              value={TabValue.Body}
              label={"Body" + (body.length ? ` (${body.length} byte` + (body.length > 1 ? "s" : "") + ")" : "")}
              sx={tabSx}
            />
          </TabList>
        </Box>
        <Box flex="1 auto" overflow="hidden">
          <TabPanel value={TabValue.QueryParams} sx={{ p: 0, height: "100%", overflow: "scroll" }}>
            <Box>
              <KeyValuePairTable items={queryParams} onChange={onQueryParamChange} onDelete={onQueryParamDelete} />
            </Box>
          </TabPanel>
          <TabPanel value={TabValue.Headers} sx={{ p: 0, height: "100%", overflow: "scroll" }}>
            <Box>
              <KeyValuePairTable items={headers} onChange={onHeaderChange} onDelete={onHeaderDelete} />
            </Box>
          </TabPanel>
          <TabPanel value={TabValue.Body} sx={{ p: 0, height: "100%" }}>
            <Editor
              content={body}
              onChange={(value) => {
                onBodyChange(value || "");
              }}
              monacoOptions={{ readOnly: false }}
              contentType={headers.find(({ key }) => key.toLowerCase() === "content-type")?.value}
            />
          </TabPanel>
        </Box>
      </TabContext>
    </Box>
  );
}

export default EditRequestTabs;
