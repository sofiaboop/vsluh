/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './app/**/*.{vue,js,ts}',
    './app/app.vue',
    './nuxt.config.ts',
  ],
  theme: {
    extend: {
      colors: {
        // Токены сняты с макета Figma (froze-chomp-91118995.figma.site)
        canvas: '#F3F6FA',
        ink: {
          DEFAULT: '#0A1020',
          strong: '#0E1424',
          muted: '#344157',
          soft: '#5B6880',
        },
        brand: {
          DEFAULT: '#408FDC',
          strong: '#327FC8',
          border: '#7DBBFF',
          glow: '#5EA9FF',
        },
      },
      backgroundColor: {
        chip: 'rgba(255, 255, 255, 0.62)',
        'chip-active': 'rgba(255, 255, 255, 0.32)',
        glass: 'rgba(255, 255, 255, 0.46)',
      },
      borderColor: {
        chip: 'rgba(10, 16, 32, 0.07)',
        glass: 'rgba(94, 169, 255, 0.40)',
      },
      boxShadow: {
        chip: '0 6px 17px 0 rgba(31, 54, 83, 0.035)',
        chipHover: '0 10px 24px 0 rgba(31, 54, 83, 0.07)',
        cta: '0 12px 30px 0 rgba(31, 54, 83, 0.09)',
      },
      borderRadius: {
        pill: '999px',
      },
      maxWidth: {
        stage: '740px',
      },
      fontFamily: {
        sans: ['Manrope', 'SF Pro Display', '-apple-system', 'system-ui', 'sans-serif'],
      },
      letterSpacing: {
        logo: '0.19em',
        display: '-0.055em',
      },
      fontSize: {
        // 3.75vw — заголовок макета: 54px при ширине окна 1440px
        display: ['clamp(2rem, 3.75vw, 4.5rem)', { lineHeight: '1.05' }],
      },
      keyframes: {
        'fade-up': {
          '0%': { opacity: '0', transform: 'translateY(12px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        drift: {
          '0%, 100%': { transform: 'translate3d(0, 0, 0) scale(1)' },
          '50%': { transform: 'translate3d(2%, -3%, 0) scale(1.05)' },
        },
      },
      animation: {
        'fade-up': 'fade-up 0.6s cubic-bezier(0.16, 1, 0.3, 1) both',
        drift: 'drift 18s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
