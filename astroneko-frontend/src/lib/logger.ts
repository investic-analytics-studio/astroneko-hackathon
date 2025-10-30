/**
 * Logger Utility
 *
 * Centralized logging utility for the application
 * Provides consistent logging with environment-based control
 */

/**
 * Log levels
 */
export enum LogLevel {
  DEBUG = 'debug',
  INFO = 'info',
  WARN = 'warn',
  ERROR = 'error',
}

/**
 * Logger interface
 */
interface Logger {
  debug: (message: string, ...args: unknown[]) => void;
  info: (message: string, ...args: unknown[]) => void;
  warn: (message: string, ...args: unknown[]) => void;
  error: (message: string, ...args: unknown[]) => void;
}

/**
 * Format log message with timestamp
 */
const formatMessage = (level: LogLevel, message: string): string => {
  const timestamp = new Date().toISOString();
  return `[${timestamp}] [${level.toUpperCase()}] ${message}`;
};

/**
 * Create logger instance
 */
const createLogger = (): Logger => {
  const shouldLog = import.meta.env.DEV;

  return {
    debug: (message: string, ...args: unknown[]) => {
      if (shouldLog) {
        console.debug(formatMessage(LogLevel.DEBUG, message), ...args);
      }
    },

    info: (message: string, ...args: unknown[]) => {
      if (shouldLog) {
        console.info(formatMessage(LogLevel.INFO, message), ...args);
      }
    },

    warn: (message: string, ...args: unknown[]) => {
      if (shouldLog) {
        console.warn(formatMessage(LogLevel.WARN, message), ...args);
      }
    },

    error: (message: string, ...args: unknown[]) => {
      // Always log errors, even in production
      console.error(formatMessage(LogLevel.ERROR, message), ...args);
    },
  };
};

/**
 * Default logger instance
 * Use this instead of direct console.log statements
 *
 * @example
 * ```ts
 * import { logger } from '@/lib/logger';
 *
 * // Development only
 * logger.debug('User clicked button', { userId: '123' });
 * logger.info('API request successful');
 *
 * // Always logged
 * logger.error('Failed to fetch data', error);
 * ```
 */
export const logger = createLogger();

/**
 * Create a namespaced logger for a specific module
 *
 * @example
 * ```ts
 * const authLogger = createNamespacedLogger('Auth');
 * authLogger.info('User logged in'); // [INFO] [Auth] User logged in
 * ```
 */
export const createNamespacedLogger = (namespace: string): Logger => {
  return {
    debug: (message: string, ...args: unknown[]) =>
      logger.debug(`[${namespace}] ${message}`, ...args),
    info: (message: string, ...args: unknown[]) =>
      logger.info(`[${namespace}] ${message}`, ...args),
    warn: (message: string, ...args: unknown[]) =>
      logger.warn(`[${namespace}] ${message}`, ...args),
    error: (message: string, ...args: unknown[]) =>
      logger.error(`[${namespace}] ${message}`, ...args),
  };
};