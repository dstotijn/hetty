import { Menu } from "@mui/material";
import React, { useState } from "react";

interface ContextMenuProps {
  children?: React.ReactNode;
}

export default function useContextMenu(): [
  (props: ContextMenuProps) => JSX.Element,
  (e: React.MouseEvent) => void,
  () => void
] {
  const [contextMenu, setContextMenu] = useState<{
    mouseX: number;
    mouseY: number;
  } | null>(null);

  const handleContextMenu = (event: React.MouseEvent) => {
    event.preventDefault();
    setContextMenu(
      contextMenu === null
        ? {
            mouseX: event.clientX - 2,
            mouseY: event.clientY - 4,
          }
        : // repeated contextmenu when it is already open closes it with Chrome 84 on Ubuntu
          // Other native context menus might behave different.
          // With this behavior we prevent contextmenu from the backdrop to re-locale existing context menus.
          null
    );
  };

  const handleClose = () => {
    setContextMenu(null);
  };

  const menu = ({ children }: ContextMenuProps): JSX.Element => (
    <Menu
      open={contextMenu !== null}
      onClose={handleClose}
      anchorReference="anchorPosition"
      anchorPosition={contextMenu !== null ? { top: contextMenu.mouseY, left: contextMenu.mouseX } : undefined}
    >
      {children}
    </Menu>
  );

  return [menu, handleContextMenu, handleClose];
}
