/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        snisid: {
          black: '#0a0a0f',
          dark: '#12121a',
          primary: '#00d4ff',
          secondary: '#7c3aed',
          danger: '#ef4444',
          warning: '#f59e0b',
          success: '#10b981',
          info: '#3b82f6',
          nsa: {
            blue: '#003366',
            gold: '#c5a059',
            red: '#8b0000',
          },
          china: {
            red: '#de2910',
            yellow: '#ffde00',
          },
          fbi: {
            blue: '#002868',
            gold: '#bf9b30',
          },
        },
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'scan': 'scan 2s linear infinite',
        'alert': 'alert 0.5s ease-in-out infinite alternate',
      },
      keyframes: {
        scan: {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100%)' },
        },
        alert: {
          '0%': { opacity: '0.5', transform: 'scale(1)' },
          '100%': { opacity: '1', transform: 'scale(1.05)' },
        },
      },
      backgroundImage: {
        'grid-pattern': "linear-gradient(to right, rgba(0, 212, 255, 0.05) 1px, transparent 1px), linear-gradient(to bottom, rgba(0, 212, 255, 0.05) 1px, transparent 1px)",
        'radar-gradient': 'radial-gradient(circle, rgba(0, 212, 255, 0.1) 0%, transparent 70%)',
      },
    },
  },
  plugins: [],
}
