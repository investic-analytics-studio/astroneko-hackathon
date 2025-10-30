import { ExternalLink } from 'lucide-react';
import { Button, type ButtonProps } from './button';
import { shouldShowExternalBrowserButton, getExternalBrowserUrl, isAndroid, isiOS } from '../../lib/webviewDetection';

interface ExternalBrowserButtonProps extends Omit<ButtonProps, 'onClick' | 'asChild'> {
  targetUrl?: string;
  children?: React.ReactNode;
}

export function ExternalBrowserButton({
  targetUrl,
  children,
  variant = "outline",
  className,
  ...props
}: ExternalBrowserButtonProps) {
  if (!shouldShowExternalBrowserButton()) {
    return null;
  }

  const handleClick = () => {
    const url = getExternalBrowserUrl(targetUrl);
    window.open(url, '_blank');
  };

  const getPlatformText = () => {
    if (isAndroid()) return 'Open in Browser';
    if (isiOS()) return 'Open in Safari';
    return 'Open in Browser';
  };

  return (
    <Button
      variant={variant}
      onClick={handleClick}
      className={className}
      {...props}
    >
      <ExternalLink className="h-4 w-4" />
      {children || getPlatformText()}
    </Button>
  );
}