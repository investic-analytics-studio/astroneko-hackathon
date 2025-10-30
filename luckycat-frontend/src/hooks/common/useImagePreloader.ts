import { useEffect, useState } from "react";
import { preloadImages } from "@/lib/utils";
import { logger } from "@/lib/logger";

interface UseImagePreloaderOptions {
  images: string[];
  priority?: boolean;
}

export const useImagePreloader = ({
  images,
  priority = false,
}: UseImagePreloaderOptions) => {
  const [isPreloaded, setIsPreloaded] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (priority) {
      setIsPreloaded(true);
      return;
    }

    const preloadImagesAsync = async () => {
      try {
        await preloadImages(images);
        setIsPreloaded(true);
      } catch (err) {
        logger.warn("Failed to preload some images:", err);
        setError(
          err instanceof Error ? err.message : "Failed to preload images"
        );
        // Still mark as preloaded to not block the UI
        setIsPreloaded(true);
      }
    };

    preloadImagesAsync();
  }, [images, priority]);

  return { isPreloaded, error };
};
