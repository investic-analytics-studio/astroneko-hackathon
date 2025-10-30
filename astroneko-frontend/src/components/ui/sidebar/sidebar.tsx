import * as React from "react";
import { cn } from "@/lib/utils";
import { useSidebar } from "./context";

export function Sidebar({
  side = "left",
  variant = "sidebar",
  collapsible = "offcanvas",
  className,
  children,
  ...props
}: React.ComponentProps<"div"> & {
  side?: "left" | "right";
  variant?: "sidebar" | "floating" | "inset";
  collapsible?: "offcanvas" | "icon" | "none";
}) {
  const { isMobile, state, openMobile } = useSidebar();

  if (collapsible === "none") {
    return (
      <div
        className={cn(
          "flex h-full w-[--sidebar-width] flex-col bg-sidebar text-sidebar-foreground",
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }

  if (isMobile) {
    return (
      <div
        data-state={openMobile ? "open" : "closed"}
        className={cn(
          "fixed inset-y-0 z-50 flex h-full w-[--sidebar-width] flex-col bg-sidebar text-sidebar-foreground md:hidden",
          side === "left" ? "left-0" : "right-0",
          variant === "floating" ? "m-2 rounded-lg border" : "border-r",
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }

  return (
    <div
      data-state={state}
      data-collapsible={state === "collapsed" ? collapsible : ""}
      className={cn(
        "group peer hidden md:block",
        "data-[collapsible=offcanvas]:translate-x-0 data-[collapsible=offcanvas]:shadow-lg",
        "data-[collapsible=icon]:translate-x-0 data-[collapsible=icon]:shadow-lg",
        "[&[data-collapsible=offcanvas]_div]:w-[--sidebar-width]",
        "[&[data-collapsible=icon]_div]:w-[--sidebar-width-icon]",
        "[&[data-collapsible=offcanvas][data-state=collapsed]]:translate-x-[calc(-1*var(--sidebar-width))]",
        "[&[data-collapsible=icon][data-state=collapsed]]:translate-x-[calc(-1*var(--sidebar-width-icon))]",
        side === "right" && "[&[data-collapsible=offcanvas][data-state=collapsed]]:translate-x-[calc(var(--sidebar-width))]",
        side === "right" && "[&[data-collapsible=icon][data-state=collapsed]]:translate-x-[calc(var(--sidebar-width-icon))]",
        "duration-200 ease-linear",
        "flex h-full w-[--sidebar-width] flex-col bg-sidebar text-sidebar-foreground",
        side === "left" ? "border-r" : "border-l",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}