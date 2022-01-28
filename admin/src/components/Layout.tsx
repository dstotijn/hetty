import React from "react";
import {
  Theme,
  useTheme,
  Toolbar,
  IconButton,
  Typography,
  Divider,
  List,
  Tooltip,
  styled,
  CSSObject,
  Box,
  ListItemText,
} from "@mui/material";
import MuiAppBar, { AppBarProps as MuiAppBarProps } from "@mui/material/AppBar";
import MuiDrawer from "@mui/material/Drawer";
import MuiListItemButton, { ListItemButtonProps } from "@mui/material/ListItemButton";
import MuiListItemIcon, { ListItemIconProps } from "@mui/material/ListItemIcon";
import Link from "next/link";
import MenuIcon from "@mui/icons-material/Menu";
import HomeIcon from "@mui/icons-material/Home";
import SettingsEthernetIcon from "@mui/icons-material/SettingsEthernet";
import SendIcon from "@mui/icons-material/Send";
import FolderIcon from "@mui/icons-material/Folder";
import LocationSearchingIcon from "@mui/icons-material/LocationSearching";
import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";

export enum Page {
  Home,
  GetStarted,
  Projects,
  ProxySetup,
  ProxyLogs,
  Sender,
  Scope,
}

const drawerWidth = 240;

const openedMixin = (theme: Theme): CSSObject => ({
  width: drawerWidth,
  transition: theme.transitions.create("width", {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen,
  }),
  overflowX: "hidden",
});

const closedMixin = (theme: Theme): CSSObject => ({
  transition: theme.transitions.create("width", {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  overflowX: "hidden",
  width: 56,
});

const DrawerHeader = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  justifyContent: "flex-start",
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
}));

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== "open",
})<AppBarProps>(({ theme, open }) => ({
  backgroundColor: theme.palette.secondary.dark,
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(["width", "margin"], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(["width", "margin"], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== "open" })(({ theme, open }) => ({
  width: drawerWidth,
  flexShrink: 0,
  whiteSpace: "nowrap",
  boxSizing: "border-box",
  ...(open && {
    ...openedMixin(theme),
    "& .MuiDrawer-paper": openedMixin(theme),
  }),
  ...(!open && {
    ...closedMixin(theme),
    "& .MuiDrawer-paper": closedMixin(theme),
  }),
}));

const ListItemButton = styled(MuiListItemButton)<ListItemButtonProps>(({ theme }) => ({
  [theme.breakpoints.up("sm")]: {
    px: 1,
  },
  "&.MuiListItemButton-root": {
    "&.Mui-selected": {
      backgroundColor: theme.palette.primary.main,
      "& .MuiListItemIcon-root": {
        color: theme.palette.secondary.dark,
      },
      "& .MuiListItemText-root": {
        color: theme.palette.secondary.dark,
      },
    },
  },
}));

const ListItemIcon = styled(MuiListItemIcon)<ListItemIconProps>(() => ({
  minWidth: 42,
}));

interface Props {
  children: React.ReactNode;
  title: string;
  page: Page;
}

export function Layout({ title, page, children }: Props): JSX.Element {
  const theme = useTheme();
  const [open, setOpen] = React.useState(false);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  const SiteTitle = styled("span")({
    ...(title !== "" && {
      color: theme.palette.primary.main,
      marginRight: 4,
    }),
  });

  return (
    <Box sx={{ display: "flex" }}>
      <AppBar position="fixed" open={open}>
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="Open drawer"
            onClick={handleDrawerOpen}
            edge="start"
            sx={{
              mr: 5,
              ...(open && { display: "none" }),
            }}
          >
            <MenuIcon />
          </IconButton>
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-around",
              width: "100%",
            }}
          >
            <Typography variant="h5" noWrap sx={{ width: "100%" }}>
              <SiteTitle>Hetty://</SiteTitle>
              {title}
            </Typography>
            <Box sx={{ flexShrink: 0, pt: 0.75 }}>v{process.env.NEXT_PUBLIC_VERSION || "0.0"}</Box>
          </Box>
        </Toolbar>
      </AppBar>
      <Drawer variant="permanent" open={open}>
        <DrawerHeader>
          <IconButton onClick={handleDrawerClose}>
            {theme.direction === "rtl" ? <ChevronRightIcon /> : <ChevronLeftIcon />}
          </IconButton>
        </DrawerHeader>
        <Divider />
        <List sx={{ p: 0 }}>
          <Link href="/" passHref>
            <ListItemButton key="home" selected={page === Page.Home}>
              <Tooltip title="Home">
                <ListItemIcon>
                  <HomeIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Home" />
            </ListItemButton>
          </Link>
          <Link href="/proxy/logs" passHref>
            <ListItemButton key="proxyLogs" selected={page === Page.ProxyLogs}>
              <Tooltip title="Proxy">
                <ListItemIcon>
                  <SettingsEthernetIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Proxy" />
            </ListItemButton>
          </Link>
          <Link href="/sender" passHref>
            <ListItemButton key="sender" selected={page === Page.Sender}>
              <Tooltip title="Sender">
                <ListItemIcon>
                  <SendIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Sender" />
            </ListItemButton>
          </Link>
          <Link href="/scope" passHref>
            <ListItemButton key="scope" selected={page === Page.Scope}>
              <Tooltip title="Scope">
                <ListItemIcon>
                  <LocationSearchingIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Scope" />
            </ListItemButton>
          </Link>
          <Link href="/projects" passHref>
            <ListItemButton key="projects" selected={page === Page.Projects}>
              <Tooltip title="Projects">
                <ListItemIcon>
                  <FolderIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Projects" />
            </ListItemButton>
          </Link>
        </List>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <DrawerHeader />
        {children}
      </Box>
    </Box>
  );
}

export default Layout;
