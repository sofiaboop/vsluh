<script setup lang="ts">
const { t } = useI18n()
const { dialogMode, openAuth, closeAuth, submitAuth } = useAuth()

const isLogin = computed(() => dialogMode.value === 'login')

function onSubmit(event: Event) {
  const data = new FormData(event.target as HTMLFormElement)
  submitAuth({
    name: String(data.get('name') ?? ''),
    email: String(data.get('email') ?? ''),
  })
}

useEventListener(window, 'keydown', (e: KeyboardEvent) => {
  if (e.key === 'Escape') closeAuth()
})
</script>

<template>
  <Teleport to="body">
    <div
      v-if="dialogMode"
      class="auth-overlay"
      role="dialog"
      aria-modal="true"
      :aria-label="isLogin ? t('auth.login') : t('auth.register')"
      @mousedown.self="closeAuth"
    >
      <div class="auth-dialog">
        <button
          type="button"
          class="support-close"
          :aria-label="t('auth.close')"
          @click="closeAuth"
        />

        <div class="auth-tabs" role="tablist">
          <button
            type="button"
            role="tab"
            class="auth-tab"
            :class="{ 'auth-tab--active': isLogin }"
            :aria-selected="isLogin"
            @click="openAuth('login')"
          >
            {{ t('auth.login') }}
          </button>
          <button
            type="button"
            role="tab"
            class="auth-tab"
            :class="{ 'auth-tab--active': !isLogin }"
            :aria-selected="!isLogin"
            @click="openAuth('register')"
          >
            {{ t('auth.register') }}
          </button>
        </div>

        <div class="auth-heading">
          <span class="auth-kicker">{{ t('auth.kicker') }}</span>
          <h2>{{ isLogin ? t('auth.loginTitle') : t('auth.registerTitle') }}</h2>
          <p>{{ isLogin ? t('auth.loginLead') : t('auth.registerLead') }}</p>
        </div>

        <form class="auth-form" @submit.prevent="onSubmit">
          <label v-if="!isLogin" class="auth-field">
            <span>{{ t('auth.name') }}</span>
            <input name="name" type="text" autocomplete="name" required :placeholder="t('auth.namePlaceholder')" >
          </label>

          <label class="auth-field">
            <span>{{ t('auth.email') }}</span>
            <input name="email" type="email" autocomplete="email" required :placeholder="t('auth.emailPlaceholder')" >
          </label>

          <label class="auth-field">
            <span>{{ t('auth.password') }}</span>
            <input
              name="password"
              type="password"
              :autocomplete="isLogin ? 'current-password' : 'new-password'"
              minlength="6"
              required
              :placeholder="t('auth.passwordPlaceholder')"
            >
          </label>

          <button type="submit" class="auth-submit">
            {{ isLogin ? t('auth.submitLogin') : t('auth.submitRegister') }}
          </button>
        </form>

        <button type="button" class="auth-skip" @click="closeAuth">{{ t('auth.skip') }}</button>
      </div>
    </div>
  </Teleport>
</template>
