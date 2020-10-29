import React from "react";
import {
  makeStyles,
  Theme,
  createStyles,
  useTheme,
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  Drawer,
  Divider,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Tooltip,
} from "@material-ui/core";
import Link from "next/link";
import MenuIcon from "@material-ui/icons/Menu";
import HomeIcon from "@material-ui/icons/Home";
import SettingsEthernetIcon from "@material-ui/icons/SettingsEthernet";
import SendIcon from "@material-ui/icons/Send";
import FolderIcon from "@material-ui/icons/Folder";
import LocationSearchingIcon from "@material-ui/icons/LocationSearching";
import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import clsx from "clsx";

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

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      display: "flex",
      width: "100%",
    },
    appBar: {
      zIndex: theme.zIndex.drawer + 1,
      transition: theme.transitions.create(["width", "margin"], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
      }),
    },
    appBarShift: {
      marginLeft: drawerWidth,
      width: `calc(100% - ${drawerWidth}px)`,
      transition: theme.transitions.create(["width", "margin"], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen,
      }),
    },
    menuButton: {
      marginRight: 28,
    },
    hide: {
      display: "none",
    },
    drawer: {
      width: drawerWidth,
      flexShrink: 0,
      whiteSpace: "nowrap",
    },
    drawerOpen: {
      width: drawerWidth,
      transition: theme.transitions.create("width", {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen,
      }),
    },
    drawerClose: {
      transition: theme.transitions.create("width", {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
      }),
      overflowX: "hidden",
      width: theme.spacing(7) + 1,
      [theme.breakpoints.up("sm")]: {
        width: theme.spacing(7) + 8,
      },
    },
    toolbar: {
      display: "flex",
      alignItems: "center",
      justifyContent: "flex-end",
      padding: theme.spacing(0, 1),
      // necessary for content to be below app bar
      ...theme.mixins.toolbar,
    },
    content: {
      flexGrow: 1,
      padding: theme.spacing(3),
    },
    listItem: {
      paddingLeft: 16,
      paddingRight: 16,
      [theme.breakpoints.up("sm")]: {
        paddingLeft: 20,
        paddingRight: 20,
      },
    },
    listItemIcon: {
      minWidth: 42,
    },
    titleHighlight: {
      color: theme.palette.secondary.main,
      marginRight: 4,
    },
  })
);

interface Props {
  children: React.ReactNode;
  title: string;
  page: Page;
}

export function Layout({ title, page, children }: Props): JSX.Element {
  const classes = useStyles();
  const theme = useTheme();
  const [open, setOpen] = React.useState(false);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  return (
    <div className={classes.root}>
      <AppBar
        position="fixed"
        className={clsx(classes.appBar, {
          [classes.appBarShift]: open,
        })}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            onClick={handleDrawerOpen}
            edge="start"
            className={clsx(classes.menuButton, {
              [classes.hide]: open,
            })}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h5" noWrap>
            <span className={title !== "" ? classes.titleHighlight : ""}>
              Hetty://
            </span>
            {title}
          </Typography>
        </Toolbar>
      </AppBar>
      <Drawer
        variant="permanent"
        className={clsx(classes.drawer, {
          [classes.drawerOpen]: open,
          [classes.drawerClose]: !open,
        })}
        classes={{
          paper: clsx({
            [classes.drawerOpen]: open,
            [classes.drawerClose]: !open,
          }),
        }}
      >
        <div className={classes.toolbar}>
          <IconButton onClick={handleDrawerClose}>
            {theme.direction === "rtl" ? (
              <ChevronRightIcon />
            ) : (
              <ChevronLeftIcon />
            )}
          </IconButton>
        </div>
        <Divider />
        <List>
          <Link href="/" passHref>
            <ListItem
              button
              component="a"
              key="home"
              selected={page === Page.Home}
              className={classes.listItem}
            >
              <Tooltip title="Home">
                <ListItemIcon className={classes.listItemIcon}>
                  <HomeIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Home" />
            </ListItem>
          </Link>
          <Link href="/proxy/logs" passHref>
            <ListItem
              button
              component="a"
              key="proxyLogs"
              selected={page === Page.ProxyLogs}
              className={classes.listItem}
            >
              <Tooltip title="Proxy">
                <ListItemIcon className={classes.listItemIcon}>
                  <SettingsEthernetIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Proxy" />
            </ListItem>
          </Link>
          <Link href="/sender" passHref>
            <ListItem
              button
              component="a"
              key="sender"
              selected={page === Page.Sender}
              className={classes.listItem}
            >
              <Tooltip title="Sender">
                <ListItemIcon className={classes.listItemIcon}>
                  <SendIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Sender" />
            </ListItem>
          </Link>
          <Link href="/scope" passHref>
            <ListItem
              button
              component="a"
              key="scope"
              selected={page === Page.Scope}
              className={classes.listItem}
            >
              <Tooltip title="Scope">
                <ListItemIcon className={classes.listItemIcon}>
                  <LocationSearchingIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Scope" />
            </ListItem>
          </Link>
          <Link href="/projects" passHref>
            <ListItem
              button
              component="a"
              key="projects"
              selected={page === Page.Projects}
              className={classes.listItem}
            >
              <Tooltip title="Projects">
                <ListItemIcon className={classes.listItemIcon}>
                  <FolderIcon />
                </ListItemIcon>
              </Tooltip>
              <ListItemText primary="Projects" />
            </ListItem>
          </Link>
        </List>
      </Drawer>
      <main className={classes.content}>
        <div className={classes.toolbar} />
        {children}
      </main>
    </div>
  );
}

export default Layout;
