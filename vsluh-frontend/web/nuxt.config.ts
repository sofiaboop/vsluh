import { defineNuxtConfig } from 'nuxt/config'

const siteUrl = process.env.NUXT_PUBLIC_SITE_URL || 'https://vsluh.club'

export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',

  // devtools доустанавливает зависимости в рантайме и ломает npm ci
  devtools: { enabled: false },

  modules: ['@nuxtjs/tailwindcss', '@vueuse/nuxt', '@nuxtjs/google-fonts', '@nuxtjs/i18n'],

  // вложенности в CSS нет, а модуль тянет для неё отдельный пакет и ругается
  postcss: { plugins: { 'tailwindcss/nesting': false } },

  devServer: { host: '127.0.0.1', port: 3000 },

  // Один маршрут на все языки: сайт статический, язык выбирается в рантайме
  i18n: {
    baseUrl: siteUrl,
    defaultLocale: 'ru',
    strategy: 'no_prefix',
    locales: [
      { code: 'ru', language: 'ru-RU', name: 'Русский', file: 'ru.json' },
      { code: 'en', language: 'en-US', name: 'English', file: 'en.json' },
      { code: 'es', language: 'es-ES', name: 'Español', file: 'es.json' },
      { code: 'pt', language: 'pt-BR', name: 'Português', file: 'pt.json' },
      { code: 'fr', language: 'fr-FR', name: 'Français', file: 'fr.json' },
      { code: 'de', language: 'de-DE', name: 'Deutsch', file: 'de.json' },
      { code: 'it', language: 'it-IT', name: 'Italiano', file: 'it.json' },
      { code: 'pl', language: 'pl-PL', name: 'Polski', file: 'pl.json' },
      { code: 'tr', language: 'tr-TR', name: 'Türkçe', file: 'tr.json' },
      { code: 'ar', language: 'ar-SA', name: 'العربية', dir: 'rtl', file: 'ar.json' },
    ],
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'vsluh_locale',
      cookieSecure: true,
      alwaysRedirect: false,
      fallbackLocale: 'en',
      redirectOn: 'root',
    },
  },

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
      // вход обкатывается на dev раньше, чем включается на проде
      featureAuth: process.env.NUXT_PUBLIC_FEATURE_AUTH === 'true',
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
        { name: 'theme-color', content: '#f3f6fa' },
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'ВСЛУХ' },
        { property: 'og:url', content: siteUrl },
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
