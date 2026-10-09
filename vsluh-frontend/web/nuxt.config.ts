import { defineNuxtConfig } from 'nuxt/config'

const siteUrl = process.env.NUXT_PUBLIC_SITE_URL || 'https://vsluh.club'

export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',

  // devtools доустанавливает зависимости в рантайме и ломает npm ci
  devtools: { enabled: false },

  modules: ['@nuxtjs/tailwindcss', '@vueuse/nuxt', '@nuxtjs/google-fonts'],

  devServer: { host: '127.0.0.1', port: 3000 },

  components: [
    { path: '~/components', pathPrefix: false },
    { path: '~/widgets', pathPrefix: false },
    { path: '~/features', pathPrefix: false, pattern: '**/ui/**/*.vue' },
    { path: '~/entities', pathPrefix: false, pattern: '**/ui/**/*.vue' },
  ],

  // Набор начертаний — ровно как в макете
  googleFonts: {
    families: { Manrope: [400, 500, 600, 700] },
    display: 'swap',
    preload: true,
  },

  css: ['~/assets/css/tailwind.css'],

  runtimeConfig: {
    public: {
      siteUrl,
      apiUrl: process.env.NUXT_PUBLIC_API_URL || '',
      botName: process.env.NUXT_PUBLIC_BOT_NAME || '',
      appVersion: process.env.NUXT_PUBLIC_APP_VERSION || 'local',
    },
  },

  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      title: 'ВСЛУХ | Сервис развития речи',
      meta: [
        { charset: 'UTF-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1.0' },
        {
          name: 'description',
          content:
            'ВСЛУХ — тренажёр спонтанной речи. 30 секунд на подготовку, минута, чтобы говорить. Выбирай темы и практикуйся каждый день.',
        },
        { name: 'theme-color', content: '#f3f6fa' },
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'ВСЛУХ' },
        { property: 'og:title', content: 'ВСЛУХ | Сервис развития речи' },
        {
          property: 'og:description',
          content: '30 секунд на подготовку, минута, чтобы говорить. Практика спонтанной речи.',
        },
        { property: 'og:url', content: siteUrl },
        { property: 'og:locale', content: 'ru_RU' },
        { name: 'twitter:card', content: 'summary_large_image' },
      ],
      link: [{ rel: 'icon', href: '/favicon.svg', type: 'image/svg+xml' }],
    },
  },

  nitro: {
    preset: 'static',
    prerender: { crawlLinks: true, routes: ['/'] },
  },

  typescript: { strict: true, typeCheck: false },
})
