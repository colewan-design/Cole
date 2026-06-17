/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        j: {
          bg:      'var(--color-bg)',
          surface: 'var(--color-surface)',
          card:    'var(--color-surface-2)',
          border:  'var(--color-border)',
          accent:  '#007AFF',
          glow:    '#007AFF',
          text:    'var(--color-text)',
          muted:   '#8E8E93',
          success: '#34C759',
          warn:    '#FF9500',
          error:   '#FF3B30',
        },
      },
      fontFamily: {
        sans: ['-apple-system', 'SF Pro Display', 'Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'Consolas', 'monospace'],
      },
      boxShadow: {
        glow:    '0 2px 12px rgba(0,122,255,0.3)',
        'glow-sm': '0 1px 6px rgba(0,122,255,0.2)',
        card:    'var(--shadow-card)',
      },
    },
  },
  plugins: [],
}
