/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ['class'],
  content: [
    './pages/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
    './app/**/*.{ts,tsx}',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
    extend: {
      textColor: {
        default: '#000',
      },
      colors: {
        white: "#fff",
        "slate-light-11": "#687076",
        "slate-light-7": "#d7dbdf",
        "slate-light-12": "#11181c",
        "slate-light-10": "#7e868c",
        "slate-light-3": "#f1f3f5",
      },
      spacing: {},
      fontFamily: {
        sans: [
          '"Open Sans", sans-serif',
        ],
      },
      borderRadius: {
        "181xl": "200px",
      },
      keyframes: {
        'accordion-down': {
          from: { height: 0 },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: 0 },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
      },
      backgroundImage: {
          'djb': "url('/background.png')",
        }
    },
  },
  plugins: [
    require('tailwindcss-animate'),
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ],
};