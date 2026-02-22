<wizard-report>
# PostHog post-wizard report

The wizard has completed a deep integration of PostHog analytics into the Muziboo Astro (hybrid) project. Here's a summary of all changes made:

## What was set up

- **Client-side analytics** via the PostHog JS snippet, initialised in `src/components/posthog.astro` and loaded on every page through `src/layouts/Layout.astro`
- **Server-side analytics** via `posthog-node`, using a singleton client in `src/lib/posthog-server.ts`, called from the contact API route
- **User identification** on both signup and contact form submission — PostHog's `identify()` is called with the user's email as the distinct ID
- **Session correlation** — the contact API route reads `X-PostHog-Session-Id` and `X-PostHog-Distinct-Id` headers to stitch client and server events together
- **Error tracking** — `captureException` is called on network/form errors client-side, and `$exception` events are captured server-side in the API route's catch block
- **Environment variables** — PostHog keys are stored in `.env` (gitignored) and referenced via `PUBLIC_POSTHOG_KEY` / `PUBLIC_POSTHOG_HOST`

## Files created

| File                           | Purpose                                         |
| ------------------------------ | ----------------------------------------------- |
| `src/components/posthog.astro` | PostHog JS snippet component (client-side)      |
| `src/lib/posthog-server.ts`    | posthog-node singleton for server-side tracking |

## Files edited

| File                                 | Changes                                                                                                                  |
| ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------ |
| `src/layouts/Layout.astro`           | Imports and renders `<PostHog />` in `<head>`                                                                            |
| `src/pages/api/contact.ts`           | Added `export const prerender = false`, server-side `contact_form_submitted_server` capture, and error tracking          |
| `src/components/MuzibooHero.astro`   | Added `id="hero-cta-btn"` to CTA link and `signup_cta_clicked` capture script                                            |
| `src/components/MuzibooSignup.astro` | Added `signup_form_submitted` capture + `identify()` on submit                                                           |
| `src/components/ContactForm.astro`   | Added `contact_form_submitted` + `identify()` on success, `contact_form_submission_failed` + `captureException` on error |
| `src/types/astro.d.ts`               | Added `declare global { interface Window { posthog? } }` type declaration                                                |

## Events instrumented

| Event name                       | Description                                                         | File                                 |
| -------------------------------- | ------------------------------------------------------------------- | ------------------------------------ |
| `signup_cta_clicked`             | User clicked the "Join Early Access" CTA button in the hero section | `src/components/MuzibooHero.astro`   |
| `signup_form_submitted`          | User submitted the early access signup form (with identify)         | `src/components/MuzibooSignup.astro` |
| `contact_form_submitted`         | Contact form successfully submitted via Web3Forms (with identify)   | `src/components/ContactForm.astro`   |
| `contact_form_submission_failed` | Contact form submission failed (network error or Web3Forms error)   | `src/components/ContactForm.astro`   |
| `contact_form_submitted_server`  | Server-side: contact form submitted via the custom API route        | `src/pages/api/contact.ts`           |

## Next steps

We've built some insights and a dashboard for you to keep an eye on user behaviour, based on the events we just instrumented:

- 📊 **Dashboard — Analytics basics**: https://eu.posthog.com/project/130233/dashboard/535818
- 🔀 **Signup Conversion Funnel** (CTA click → form submission): https://eu.posthog.com/project/130233/insights/G18JCtrj
- 📈 **Signups & Contact Form Submissions (Daily)**: https://eu.posthog.com/project/130233/insights/ejE5khm7
- 📊 **Hero CTA Clicks (Daily)**: https://eu.posthog.com/project/130233/insights/JjLeaZkW
- ⚠️ **Contact Form Failure Rate**: https://eu.posthog.com/project/130233/insights/Q47qTCJp
- 🔀 **Contact Form Success vs Failure**: https://eu.posthog.com/project/130233/insights/FJ7Wah6c

### Agent skill

We've left an agent skill folder in your project at `.claude/skills/posthog-integration-astro-hybrid/`. You can use this context for further agent development when using Claude Code. This will help ensure the model provides the most up-to-date approaches for integrating PostHog.

</wizard-report>
