<script setup lang="ts">
const { t } = useI18n()
</script>

<template>
  <main class="app-shell">
    <div class="ambient ambient--one" />
    <div class="ambient ambient--two" />

    <TopBar />

    <slot />

    <footer class="footer">
      <span>{{ t('footer') }}</span>
    </footer>
  </main>
</template>
