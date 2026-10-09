<script setup lang="ts">
import { WHEEL } from '~/shared/constants/practice'

const { rows, spinning, settled, mood, chosen, startPrepare, spin } = usePractice()

function distance(index: number) {
  // пока крутится, подсвечена стартовая строка; после остановки — финальная
  const anchor = settled.value ? WHEEL.selectedIndex : WHEEL.startIndex
  return Math.min(Math.abs(index - anchor), 4)
}
</script>

<template>
  <section class="content wheel-content">
    <div class="step-label">02 / ВЫБОР ТЕМЫ</div>

    <div class="wheel-wrap">
      <div
        class="wheel-list"
        :class="{ 'wheel-list--spinning': spinning, 'wheel-list--settled': settled }"
      >
        <div
          v-for="(row, index) in rows"
          :key="index"
          class="wheel-row"
          :class="[
            distance(index) === 0 ? 'wheel-row--selected' : '',
            `wheel-row--distance-${distance(index)}`,
          ]"
        >
          <span class="wheel-category">{{ row.topic }}</span>
          <span class="wheel-topic">{{ row.question }}</span>
        </div>
      </div>
    </div>

    <div v-if="chosen" class="chosen-mobile">
      <span>{{ chosen.topic }}</span>{{ chosen.question }}
    </div>

    <div v-if="mood" class="mood-pill"><span class="mood-dot" />{{ mood }}</div>

    <div class="wheel-actions">
      <button
        type="button"
        class="action-button action-button--primary"
        :disabled="spinning"
        @click="startPrepare"
      >
        Начать
      </button>
      <button
        type="button"
        class="action-button action-button--secondary"
        :disabled="spinning"
        @click="spin"
      >
        <span class="spin-icon" :class="{ 'spin-icon--active': spinning }">↻</span>
        {{ spinning ? 'Крутим…' : 'Крутить ещё' }}
      </button>
    </div>
  </section>
</template>
