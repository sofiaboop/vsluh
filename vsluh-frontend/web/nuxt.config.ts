import { defineNuxtConfig } from 'nuxt/config'

const siteUrl = process.env.NUXT_PUBLIC_SITE_URL || 'https://vsluh.club'

export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },

  modules: ['@nuxtjs/tailwindcss', '@vueuse/nuxt', '@nuxtjs/google-fonts'],

  devServer: { host: '127.0.0.1', port: 3000 },

  // FSD: ui-слои фич и сущностей регистрируются автоматически
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
      siteUrl,
      apiUrl: process.env.NUXT_PUBLIC_API_URL || '',
      miniAppUrl: process.env.NUXT_PUBLIC_MINIAPP_URL || '',
      appVersion: process.env.NUXT_PUBLIC_APP_VERSION || 'local',
    },
  },

  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      title: 'ВСЛУХ — сервис развития речи',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content:
            'ВСЛУХ — тренажёр спонтанной речи. 30 секунд на подготовку, минута, чтобы говорить. Выбирай темы и практикуйся каждый день.',
        },
        { name: 'theme-color', content: '#F3F6FA' },
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'ВСЛУХ' },
        { property: 'og:title', content: 'ВСЛУХ — сервис развития речи' },
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

  // Статика в S3/CloudFront
  nitro: {
    preset: 'static',
    prerender: { crawlLinks: true, routes: ['/'] },
  },

  typescript: { strict: true, typeCheck: false },
})
