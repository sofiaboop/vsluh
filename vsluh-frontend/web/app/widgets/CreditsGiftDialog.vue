<script setup lang="ts">
const { t } = useI18n()
const { giftCredits, closeGift } = useAuth()
</script>

<template>
  <Teleport to="body">
    <div class="gift-overlay" @click.self="closeGift">
      <div aria-modal="true" class="gift-dialog" role="dialog">
        <button :aria-label="t('auth.close')" class="gift-close" type="button" @click="closeGift" />
        <span aria-hidden="true" class="gift-emoji">🎉</span>
        <h2>{{ t('credits.gift.title', { count: giftCredits }) }}</h2>
        <p>{{ t('credits.gift.text', { count: giftCredits }) }}</p>
        <button class="gift-cta" type="button" @click="closeGift">{{ t('credits.gift.cta') }}</button>
      </div>
    </div>
  </Teleport>
</template>
