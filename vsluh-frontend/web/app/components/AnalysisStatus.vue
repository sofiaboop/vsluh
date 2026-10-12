<script setup lang="ts">
const { variant = 'enabled' } = defineProps<{ variant?: 'enabled' | 'planned' | 'recording' }>()

const { t } = useI18n()
const label = computed(() => t(`analysis.${variant}`))
</script>

<template>
  <div class="analysis-status" :class="`analysis-status--${variant}`" role="status">
    <span aria-hidden="true" />{{ label }}
  </div>
</template>
