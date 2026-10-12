<script setup lang="ts">
const { locale, locales, setLocale, t } = useI18n()

const root = ref<HTMLElement | null>(null)
const open = ref(false)
const current = computed(() => locales.value.find((item) => item.code === locale.value))

function choose(code: string) {
  setLocale(code as typeof locale.value)
  open.value = false
}

onClickOutside(root, () => (open.value = false))
useEventListener(window, 'keydown', (e: KeyboardEvent) => {
  if (e.key === 'Escape') open.value = false
})
</script>

<template>
  <div ref="root" class="lang">
    <button
      type="button"
      class="lang-button"
      :aria-label="t('language.switch')"
      :aria-expanded="open"
      @click="open = !open"
    >
      <span class="lang-globe" aria-hidden="true" />
      <span class="lang-code">{{ current?.code.toUpperCase() }}</span>
    </button>

    <ul v-if="open" class="lang-menu" role="listbox">
      <li v-for="item in locales" :key="item.code">
        <button
          type="button"
          class="lang-option"
          :class="{ 'lang-option--active': item.code === locale }"
          role="option"
          :aria-selected="item.code === locale"
          @click="choose(item.code)"
        >
          {{ item.name }}
        </button>
      </li>
    </ul>
  </div>
</template>
