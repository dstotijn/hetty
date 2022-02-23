import { Box } from "@mui/system";
import { AllotmentProps } from "allotment";
import { PaneProps } from "allotment/dist/types/src/allotment";
import { ComponentType, useEffect, useRef, useState } from "react";

import { Layout, Page } from "features/Layout";
import EditRequest from "features/sender/components/EditRequest";
import History from "features/sender/components/History";

function Index(): JSX.Element {
  const isMountedRef = useRef(false);
  const [Allotment, setAllotment] = useState<
    (ComponentType<AllotmentProps> & { Pane: ComponentType<PaneProps> }) | null
  >(null);
  useEffect(() => {
    isMountedRef.current = true;
    import("allotment")
      .then((mod) => {
        if (!isMountedRef.current) {
          return;
        }
        setAllotment(mod.Allotment);
      })
      .catch((err) => console.error(err, `could not import allotment ${err.message}`));
    return () => {
      isMountedRef.current = false;
    };
  }, []);
  if (!Allotment) {
    return <div>Loading...</div>;
  }

  return (
    <Layout page={Page.Sender} title="Sender">
      <Allotment vertical={true} defaultSizes={[70, 30]}>
        <Box sx={{ pt: 0.75, height: "100%" }}>
          <EditRequest />
        </Box>
        <Box sx={{ height: "100%", py: 2, overflow: "hidden" }}>
          <Box sx={{ height: "100%", overflow: "scroll" }}>
            <History />
          </Box>
        </Box>
      </Allotment>
    </Layout>
  );
}

export default Index;
