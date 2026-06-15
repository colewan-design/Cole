/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        j: {
          bg:      '#080d1a',
          surface: '#0d1526',
          card:    '#111c33',
          border:  '#1a3050',
          accent:  '#2563eb',
          glow:    '#3b82f6',
          text:    '#e2e8f0',
          muted:   '#64748b',
          success: '#22c55e',
          warn:    '#f59e0b',
          error:   '#ef4444',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'Consolas', 'monospace'],
      },
      boxShadow: {
        glow: '0 0 20px rgba(37,99,235,0.3)',
        'glow-sm': '0 0 8px rgba(37,99,235,0.2)',
      },
    },
  },
  plugins: [],
}
