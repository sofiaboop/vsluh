import { defineNuxtConfig } from 'nuxt/config'

export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },

  modules: ['@nuxtjs/tailwindcss', '@vueuse/nuxt', '@nuxtjs/google-fonts'],

  devServer: { host: '127.0.0.1', port: 3001 },

  components: [
    { path: '~/components', pathPrefix: false },
    { path: '~/widgets', pathPrefix: false },
    { path: '~/features', pathPrefix: false, pattern: '**/ui/**/*.vue' },
    { path: '~/entities', pathPrefix: false, pattern: '**/ui/**/*.vue' },
  ],

  googleFonts: {
    families: { Manrope: [400, 500, 600, 700, 800] },
    display: 'swap',
    preload: true,
  },

  css: ['~/assets/css/tailwind.css'],

  runtimeConfig: {
    public: {
      apiUrl: process.env.NUXT_PUBLIC_API_URL || 'http://127.0.0.1:3000',
      botName: process.env.NUXT_PUBLIC_BOT_NAME || '',
      appVersion: process.env.NUXT_PUBLIC_APP_VERSION || 'local',
    },
  },

  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      title: 'ВСЛУХ',
      meta: [
        { charset: 'utf-8' },
        {
          name: 'viewport',
          content: 'width=device-width, initial-scale=1, viewport-fit=cover, user-scalable=no',
        },
        { name: 'robots', content: 'noindex, nofollow' },
      ],
      // Официальный SDK должен загрузиться до старта приложения
      script: [{ src: 'https://telegram.org/js/telegram-web-app.js' }],
      link: [{ rel: 'icon', href: '/favicon.svg', type: 'image/svg+xml' }],
    },
  },

  nitro: {
    preset: 'static',
    prerender: { routes: ['/'] },
  },

  typescript: { strict: true, typeCheck: false },
})
