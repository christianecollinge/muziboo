import typography from "@tailwindcss/typography";

/** @type {import('tailwindcss').Config} */
export default {
  plugins: [typography],
  theme: {
    extend: {
      colors: {
        muziboo: {
          bg: "var(--muziboo-bg)",
          "bg-softer": "var(--muziboo-bg-softer)",
          gold: "var(--muziboo-gold)",
          "gold-muted": "var(--muziboo-gold-muted)",
          text: "var(--muziboo-text)",
          "text-muted": "var(--muziboo-text-muted)",
          red: "var(--muziboo-red)",
          orange: "var(--muziboo-orange)",
          teal: "var(--muziboo-teal)",
          "teal-light": "var(--muziboo-teal-light)",
          bar: {
            1: "var(--muziboo-bar-1)",
            2: "var(--muziboo-bar-2)",
            3: "var(--muziboo-bar-3)",
            4: "var(--muziboo-bar-4)",
            5: "var(--muziboo-bar-5)",
            6: "var(--muziboo-bar-6)",
          },
          border: "var(--muziboo-border)",
        },
      },
    },
  },
};
