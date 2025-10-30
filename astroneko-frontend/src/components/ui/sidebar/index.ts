/**
 * Sidebar components
 * Modularized for better maintainability
 */

// Core components
export { SidebarProvider } from "./provider";
export { Sidebar } from "./sidebar";
export { useSidebar } from "./context";

// UI Components
export {
  SidebarTrigger,
  SidebarRail,
  SidebarInset,
  SidebarInput,
  SidebarHeader,
  SidebarFooter,
  SidebarSeparator,
  SidebarContent,
} from "./components";

// Menu components
export {
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarTooltipProvider,
} from "./menu";

// Constants
export * from "./constants";