<script setup lang="ts">
import { SUPPORT } from '~/shared/constants/practice'
import { FEATURES } from '~/shared/constants/features'

const { t } = useI18n()
const { backToTopics } = usePractice()
const { isAuthenticated, dialogMode, openAuth, logout } = useAuth()
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
    <button :aria-label="t('nav.home')" class="brand" type="button" @click="backToTopics">
      <span class="brand-mark"><span /></span>
      <span>ВСЛУХ</span>
    </button>

    <div class="topbar-actions">
      <LanguagePicker />

      <template v-if="FEATURES.auth">
        <button v-if="isAuthenticated" class="auth-nav-button" type="button" @click="logout">
          {{ t('auth.logout') }}
        </button>
        <div v-else class="auth-nav">
          <button class="auth-nav-button" type="button" @click="openAuth('login')">
            {{ t('auth.login') }}
          </button>
          <button
            class="auth-nav-button auth-nav-button--register"
            type="button"
            @click="openAuth('register')"
          >
            {{ t('auth.register') }}
          </button>
        </div>
      </template>

      <button
        v-if="FEATURES.support"
        class="support-button"
        type="button"
        @click="supportOpen = true"
      >
        {{ t('nav.support') }}
      </button>
    </div>
  </header>

  <AuthDialog v-if="FEATURES.auth && dialogMode" />

  <Teleport to="body">
    <div v-if="supportOpen" class="support-overlay" @click.self="close">
      <div class="support-dialog" role="dialog" aria-modal="true">
        <button :aria-label="t('support.close')" class="support-close" type="button" @click="close" />
        <img :alt="t('support.portraitAlt')" class="support-portrait" src="/support-portrait.png" />
        <h2>{{ t('support.title') }}</h2>
        <p>{{ t('support.text') }}</p>
        <a class="support-link" :href="SUPPORT.payUrl" rel="noreferrer" target="_blank">
          {{ t('support.cta') }}
        </a>
        <div class="support-dialog-footer">
          {{ t('support.madeBy') }}
          <a class="support-inline-link" :href="SUPPORT.authorUrl" rel="noreferrer" target="_blank">
            {{ SUPPORT.authorHandle }}
          </a>
        </div>
      </div>
    </div>
  </Teleport>
</template>
