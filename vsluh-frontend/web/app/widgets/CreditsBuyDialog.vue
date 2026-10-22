<script setup lang="ts">
import { CREDIT_PACKS } from '~/shared/constants/credits'

const { t } = useI18n()
const { closeBuy } = useAuth()

// Оплата ещё не подключена: вместо платёжки показываем честное «скоро».
const soonPack = ref<number | null>(null)

useEventListener(window, 'keydown', (e: KeyboardEvent) => {
  if (e.key === 'Escape') closeBuy()
})
</script>

<template>
  <Teleport to="body">
    <div
      :aria-label="t('credits.buy.title')"
      aria-modal="true"
      class="buy-overlay"
      role="dialog"
      @mousedown.self="closeBuy"
    >
      <div class="buy-dialog">
        <button :aria-label="t('auth.close')" class="support-close" type="button" @click="closeBuy" />
        <h2>{{ t('credits.buy.title') }}</h2>
        <p class="buy-lead">{{ t('credits.buy.lead') }}</p>

        <ul class="buy-packs">
          <li v-for="pack in CREDIT_PACKS" :key="pack.credits" class="buy-pack">
            <div class="buy-pack-info">
              <strong>{{ t('credits.buy.pack', { count: pack.credits }) }}</strong>
              <span>{{ t('credits.buy.price', { price: pack.price }) }}</span>
            </div>
            <button class="buy-pack-cta" type="button" @click="soonPack = pack.credits">
              {{ t('credits.buy.cta') }}
            </button>
          </li>
        </ul>

        <p v-if="soonPack !== null" class="buy-soon" role="status">{{ t('credits.buy.soon') }}</p>
      </div>
    </div>
  </Teleport>
</template>
