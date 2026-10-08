<script setup lang="ts">
import { SUPPORT } from '~/shared/constants/practice'

const supportOpen = ref(false)

function close() {
  supportOpen.value = false
}

onMounted(() => {
  const onKey = (e: KeyboardEvent) => {
    if (e.key === 'Escape') close()
  }
  window.addEventListener('keydown', onKey)
  onUnmounted(() => window.removeEventListener('keydown', onKey))
})
</script>

<template>
  <header class="topbar">
    <button aria-label="На главную" class="brand" type="button">
      <span class="brand-mark"><span /></span>
      <span>ВСЛУХ</span>
    </button>

    <button class="support-button" type="button" @click="supportOpen = true">Поддержать</button>
  </header>

  <Teleport to="body">
    <div v-if="supportOpen" class="support-overlay" @click.self="close">
      <div class="support-dialog" role="dialog" aria-modal="true">
        <button aria-label="Закрыть" class="support-close" type="button" @click="close" />
        <img alt="Автор проекта" class="support-portrait" src="/support-portrait.png" />
        <h2>{{ SUPPORT.title }}</h2>
        <p>{{ SUPPORT.text }}</p>
        <a class="support-link" :href="SUPPORT.payUrl" rel="noreferrer" target="_blank">
          Поддержать
        </a>
        <div class="support-dialog-footer">
          Сделано Соней ·
          <a class="support-inline-link" :href="SUPPORT.authorUrl" rel="noreferrer" target="_blank">
            {{ SUPPORT.authorHandle }}
          </a>
        </div>
      </div>
    </div>
  </Teleport>
</template>
