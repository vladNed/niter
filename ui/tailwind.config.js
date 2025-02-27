/** @type {import('tailwindcss').Config} */
import colors from 'tailwindcss/colors';
import defaultTheme from 'tailwindcss/defaultTheme';

module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Inter"', 'sans-serif'],
        outfit: ['"Outfit"', ...defaultTheme.fontFamily.sans],
        kanit: ['"Kanit"', ...defaultTheme.fontFamily.sans],
      },
      colors: {
        primary: colors.blue,
        secondary: colors.sky,
      },
    },
    screens: {
      xs: '300px',
      ...defaultTheme.screens,
    },
    backgroundImage: {
      // eslint-disable-next-line quotes
      'mvx-white': "url('../multiversx-white.svg')"
    },
  },
  plugins: []
};
