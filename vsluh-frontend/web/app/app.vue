<script setup lang="ts">
const { t } = useI18n()
const head = useLocaleHead()

useHead(() => ({
  htmlAttrs: head.value.htmlAttrs,
  title: t('meta.title'),
  meta: [
    { name: 'description', content: t('meta.description') },
    { property: 'og:title', content: t('meta.title') },
    { property: 'og:description', content: t('meta.ogDescription') },
    ...(head.value.meta ?? []),
  ],
  link: head.value.link,
}))
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
