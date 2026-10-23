<script setup lang="ts">
import { SUPPORT } from '~/shared/constants/practice'

const { t } = useI18n()
const { backToTopics } = usePractice()
const { isAuthenticated, dialogMode, giftCredits, buyOpen, credits, userName, openAuth, logout } =
  useAuth()
const features = useFeatures()
const supportOpen = ref(false)
const menuOpen = ref(false)
const menuRef = ref<HTMLElement | null>(null)

/** Первая буква имени для кружка-аватара: у старых аккаунтов имя может быть пустым. */
const initial = computed(() => (userName.value.trim()[0] ?? '?').toUpperCase())

function close() {
  supportOpen.value = false
}

function closeMenu() {
  menuOpen.value = false
}

onMounted(() => {
  const onKey = (e: KeyboardEvent) => {
    if (e.key !== 'Escape') return
    close()
    closeMenu()
  }
  // На тач-экранах наведения нет, меню открывается нажатием — значит нужен и выход по промаху.
  const onPointer = (e: PointerEvent) => {
    if (menuOpen.value && !menuRef.value?.contains(e.target as Node)) closeMenu()
  }
  window.addEventListener('keydown', onKey)
  window.addEventListener('pointerdown', onPointer)
  onUnmounted(() => {
    window.removeEventListener('keydown', onKey)
    window.removeEventListener('pointerdown', onPointer)
  })
})

watch(isAuthenticated, closeMenu)
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
          <div ref="menuRef" class="user-menu" @mouseenter="menuOpen = true" @mouseleave="closeMenu">
            <button
              :aria-expanded="menuOpen"
              aria-haspopup="menu"
              class="user-menu-button"
              type="button"
              @click="menuOpen = !menuOpen"
            >
              <span aria-hidden="true" class="user-menu-avatar">{{ initial }}</span>
              <span class="user-menu-name">{{ userName || t('nav.account') }}</span>
            </button>

            <div v-show="menuOpen" class="user-menu-drop" role="menu">
              <NuxtLink class="user-menu-item" role="menuitem" to="/history" @click="closeMenu">
                {{ t('nav.history') }}
              </NuxtLink>
              <button class="user-menu-item" role="menuitem" type="button" @click="logout">
                {{ t('auth.logout') }}
              </button>
            </div>
          </div>
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
  <CreditsBuyDialog v-if="features.auth && buyOpen" />

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
