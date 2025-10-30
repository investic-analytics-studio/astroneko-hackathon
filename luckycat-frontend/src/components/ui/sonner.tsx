import { useTheme } from "next-themes";
import { Toaster as Sonner, ToasterProps } from "sonner";

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = "system" } = useTheme();

  return (
    <Sonner
      theme={theme as ToasterProps["theme"]}
      className="toaster group"
      richColors
      style={
        {
          "--normal-bg": "var(--popover)",
          "--normal-text": "var(--popover-foreground)",
          "--normal-border": "var(--border)",
          "--error-bg": "#ef4444",
          "--error-text": "#ffffff",
          "--error-border": "#f87171",
          "--warning-bg": "#f97316",
          "--warning-text": "#ffffff",
          "--warning-border": "#fb923c",
          "--success-bg": "#ffffff",
          "--success-text": "#000000",
          "--success-border": "#d1d5db",
          "--info-bg": "#ffffff",
          "--info-text": "#000000",
          "--info-border": "#d1d5db",
        } as React.CSSProperties
      }
      {...props}
    />
  );
};

export { Toaster };
