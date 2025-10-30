import * as React from "react";
import { useState, useEffect, useRef } from "react";
import { getImagePath } from "@/lib/utils";

export interface OptimizedImageProps {
  /**
   * Image source path (relative or absolute)
   */
  src: string;
  /**
   * Alternative text for accessibility
   */
  alt: string;
  /**
   * Additional CSS classes
   */
  className?: string;
  /**
   * Loading strategy
   * @default "lazy"
   */
  loading?: "lazy" | "eager";
  /**
   * Priority loading (bypasses lazy loading)
   * @default false
   */
  priority?: boolean;
  /**
   * Callback when image loads successfully
   */
  onLoad?: () => void;
  /**
   * Callback when image fails to load
   */
  onError?: () => void;
  /**
   * Fallback image source if main source fails
   */
  fallbackSrc?: string;
}

export const OptimizedImage = React.forwardRef<
  HTMLDivElement,
  OptimizedImageProps
>(
  (
    {
      src,
      alt,
      className = "",
      loading = "lazy",
      priority = false,
      onLoad,
      onError,
      fallbackSrc,
    },
    ref
  ) => {
    const [isLoaded, setIsLoaded] = useState(false);
    const [hasError, setHasError] = useState(false);
    const [isInView, setIsInView] = useState(priority);
    const [currentSrc, setCurrentSrc] = useState<string>("");
    const imgRef = useRef<HTMLImageElement>(null);
    const observerRef = useRef<IntersectionObserver | null>(null);

    // Validate src to prevent empty string
    const validSrc =
      src && src.trim() ? src : fallbackSrc || "/cards/THE_FOOL.webp";

    // Get the correct image path for production
    const imageSrc = getImagePath(validSrc);
    const fallbackImageSrc = fallbackSrc ? getImagePath(fallbackSrc) : null;

    useEffect(() => {
      setCurrentSrc(imageSrc);
    }, [imageSrc]);

    useEffect(() => {
      if (priority) {
        setIsInView(true);
        return;
      }

      if (!imgRef.current) return;

      observerRef.current = new IntersectionObserver(
        ([entry]) => {
          if (entry.isIntersecting) {
            setIsInView(true);
            observerRef.current?.disconnect();
          }
        },
        {
          rootMargin: "50px",
          threshold: 0.1,
        }
      );

      observerRef.current.observe(imgRef.current);

      return () => {
        observerRef.current?.disconnect();
      };
    }, [priority]);

    const handleLoad = () => {
      setIsLoaded(true);
      setHasError(false);
      onLoad?.();
    };

    const handleError = () => {
      if (fallbackImageSrc && currentSrc !== fallbackImageSrc) {
        // Try fallback image
        setCurrentSrc(fallbackImageSrc);
        setHasError(false);
        setIsLoaded(false);
      } else {
        // No fallback or fallback also failed
        setHasError(true);
        onError?.();
      }
    };

    return (
      <div className={`relative ${className}`} ref={ref || imgRef}>
        {!isLoaded && !hasError && (
          <div className="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-pink-500/5 rounded-xl border-2 border-purple-500/20 animate-pulse">
            <div className="absolute inset-0 rounded-xl animate-pulse-glow" />
          </div>
        )}

        {isInView && currentSrc && (
          <img
            src={currentSrc}
            alt={alt}
            className={`${className} transition-all duration-300 ease-out ${
              isLoaded ? "opacity-100 scale-100" : "opacity-0 scale-95"
            }`}
            loading={loading}
            decoding="async"
            onLoad={handleLoad}
            onError={handleError}
            style={{
              willChange: "opacity, transform",
            }}
          />
        )}

        {hasError && (
          <div className="absolute inset-0 flex items-center justify-center bg-gray-800/50 rounded-xl">
            <span className="text-gray-400 text-sm">Failed to load image</span>
          </div>
        )}
      </div>
    );
  }
);

OptimizedImage.displayName = "OptimizedImage";
