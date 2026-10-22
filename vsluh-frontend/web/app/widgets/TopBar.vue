<script setup lang="ts">
import { SUPPORT } from '~/shared/constants/practice'

const { t } = useI18n()
const { backToTopics } = usePractice()
const { isAuthenticated, dialogMode, giftCredits, credits, openAuth, logout } = useAuth()
const features = useFeatures()
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
      <span>{{ t('nav.brand') }}</span>
    </button>

    <div class="topbar-actions">
      <LanguagePicker />

      <template v-if="features.auth">
        <div v-if="isAuthenticated" class="auth-nav">
          <span class="credits-badge" :title="t('credits.hint')">
            <span aria-hidden="true" class="credits-badge-dot" />
            {{ t('credits.badge', { count: credits }) }}
          </span>
          <button class="auth-nav-button" type="button" @click="logout">
            {{ t('auth.logout') }}
          </button>
        </div>
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
        v-if="features.support"
        class="support-button"
        type="button"
        @click="supportOpen = true"
      >
        {{ t('nav.support') }}
      </button>
    </div>
  </header>

  <AuthDialog v-if="features.auth && dialogMode" />
  <CreditsGiftDialog v-if="features.auth && giftCredits !== null" />

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
