/// <reference types="astro/client" />

<<<<<<< HEAD
=======
// PostHog global type declaration
declare global {
  interface Window {
    posthog?: {
      capture: (event: string, properties?: Record<string, unknown>) => void;
      identify: (
        distinctId: string,
        properties?: Record<string, unknown>,
      ) => void;
      reset: () => void;
      captureException: (error: unknown) => void;
      get_session_id?: () => string | null;
    };
  }
}

>>>>>>> b1623401 (feat: init muziboo site with signup and posthog)
// Extend Astro's built-in types
declare module "astro" {
  interface AstroBuiltinProps {
    class?: string;
    className?: string;
  }
}

// Component prop types
export interface ComponentProps {
  class?: string;
  className?: string;
  id?: string;
  "aria-label"?: string;
  "aria-labelledby"?: string;
  "aria-describedby"?: string;
}

// Layout props
export interface LayoutProps extends ComponentProps {
  title?: string;
  description?: string;
  lang?: string;
}

// Form data types
export interface ContactFormData {
  name: string;
  email: string;
  message: string;
}

export interface APIResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}
