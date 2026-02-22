/**
 * Error handling utilities
 */

export class AppError extends Error {
  constructor(
    message: string,
    public statusCode: number = 500,
<<<<<<< HEAD
    public code?: string
=======
    public code?: string,
>>>>>>> b1623401 (feat: init muziboo site with signup and posthog)
  ) {
    super(message);
    this.name = "AppError";
    Object.setPrototypeOf(this, AppError.prototype);
  }
}

export class ValidationError extends AppError {
<<<<<<< HEAD
  constructor(message: string, public field?: string) {
=======
  constructor(
    message: string,
    public field?: string,
  ) {
>>>>>>> b1623401 (feat: init muziboo site with signup and posthog)
    super(message, 400, "VALIDATION_ERROR");
    this.name = "ValidationError";
    Object.setPrototypeOf(this, ValidationError.prototype);
  }
}

/**
 * Handles errors and returns a user-friendly message
 */
export function handleError(error: unknown): string {
  if (error instanceof ValidationError) {
    return error.message;
  }

  if (error instanceof AppError) {
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "An unexpected error occurred. Please try again.";
}

/**
 * Safe error logger (only logs in development)
 */
export function logError(error: unknown, context?: string): void {
  if (import.meta.env.DEV) {
    console.error(context ? `[${context}]` : "[Error]", error);
  }
}
