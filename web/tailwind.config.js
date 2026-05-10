/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        deep: '#06080f',
        base: '#0b0f1a',
        surface: '#111827',
        elevated: '#1a2236',
        accent: '#8b5cf6',
        'accent-hover': '#a78bfa',
        cyan: '#06b6d4',
      },
      fontFamily: {
        sans: ['Noto Sans SC', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
    },
  },
  plugins: [],
}
