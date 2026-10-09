<script setup lang="ts">
import { WHEEL } from '~/shared/constants/practice'
import { playTick } from '~/shared/lib/sound'

const {
  rows,
  spinning,
  settled,
  selectedIndex,
  mood,
  chosen,
  startPrepare,
  spin,
  finishSpin,
  highlight,
} = usePractice()

const listRef = ref<HTMLElement | null>(null)
let raf = 0
let passed = WHEEL.startIndex

/** Пока идёт CSS-анимация, читаем её transform и щёлкаем на каждой проехавшей строке. */
function track() {
  const el = listRef.value
  if (!el) return

  const matched = getComputedStyle(el).transform.match(/^matrix(3d)?\((.+)\)$/)
  if (matched) {
    const parts = matched[2]!.split(',').map(Number)
    const translateY = matched[1] ? parts[13]! : parts[5]!
    const rowHeight = el.querySelector('.wheel-row')?.getBoundingClientRect().height ?? WHEEL.rowHeight
    const centered = Math.floor((-translateY - rowHeight / 2) / rowHeight + 0.02)

    if (centered > passed) {
      for (let i = passed + 1; i <= centered; i += 1) playTick()
      passed = centered
      highlight(centered)
    }
  }

  raf = requestAnimationFrame(track)
}

watch(spinning, (active) => {
  cancelAnimationFrame(raf)
  if (!active) return
  passed = WHEEL.startIndex
  raf = requestAnimationFrame(track)
})

onBeforeUnmount(() => cancelAnimationFrame(raf))

function onAnimationEnd(event: AnimationEvent) {
  if (event.animationName === 'wheel-spin') finishSpin()
}

function distance(index: number) {
  return Math.min(Math.abs(index - selectedIndex.value), 4)
}
</script>

<template>
  <section class="content wheel-content">
    <div class="step-label">02 / ВЫБОР ТЕМЫ</div>

    <div class="wheel-wrap">
      <div
        ref="listRef"
        class="wheel-list"
        :class="{ 'wheel-list--spinning': spinning, 'wheel-list--settled': settled }"
        @animationend="onAnimationEnd"
      >
        <div
          v-for="(row, index) in rows"
          :key="index"
          class="wheel-row"
          :class="[
            index === selectedIndex ? 'wheel-row--selected' : '',
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
