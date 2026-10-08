/** @type {import('tailwindcss').Config} */
export default {
  content: ['./app/**/*.{vue,js,ts}', './nuxt.config.ts'],
  theme: {
    extend: {
      // Значения макета. Компонентные стили лежат в app/assets/css/tailwind.css
      colors: {
        canvas: '#f3f6fa',
        ink: '#0a1020',
        brand: '#5ea9ff',
      },
      fontFamily: {
        sans: ['Manrope', 'SF Pro Display', '-apple-system', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
