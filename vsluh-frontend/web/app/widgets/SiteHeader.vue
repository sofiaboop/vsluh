<script setup lang="ts">
const supportUrl = 'https://t.me/vsluh_club'
</script>

<template>
  <header class="relative z-20 flex items-center justify-between px-5 py-6 sm:px-10 sm:py-8">
    <NuxtLink to="/" class="flex items-center gap-3" aria-label="ВСЛУХ — на главную">
      <LogoMark class="h-9 w-9 sm:h-10 sm:w-10" />
      <span class="text-[11px] font-bold tracking-logo text-ink-strong">ВСЛУХ</span>
    </NuxtLink>

    <a
      :href="supportUrl"
      target="_blank"
      rel="noopener noreferrer"
      class="rounded-pill border border-glass bg-glass px-[17px] py-2.5 text-[13px] font-semibold
             text-brand backdrop-blur-sm transition-all duration-200 hover:-translate-y-0.5
             hover:shadow-chipHover"
    >
      Поддержать
    </a>
  </header>
</template>
