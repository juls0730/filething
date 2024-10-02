/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./app.vue",
    "./error.vue",
  ],
  darkMode: "selector",
  theme: {
    extend: {
      colors: {
        base: "rgb(var(--color-base))",
        surface: "rgb(var(--color-surface))",
        overlay: "rgb(var(--color-overlay))",
        muted: "rgb(var(--color-muted))",
        subtle: "rgb(var(--color-subtle))",
        text: "rgb(var(--color-text))",
        foam: "rgb(var(--color-foam))",
        love: "rgb(var(--color-love))",
        pine: "rgb(var(--color-pine))",
        accent: "rgb(var(--color-accent))",
        "highlight-low": "rgb(var(--highlight-low))",
        "highlight-med": "rgb(var(--highlight-med))",
        "highlight-high": "rgb(var(--highlight-high))",
      },
      transitionProperty: {
        bg: "background-color",
      }
    },
  },
  future: {
    hoverOnlyWhenSupported: true,
  },
  plugins: [],
};
