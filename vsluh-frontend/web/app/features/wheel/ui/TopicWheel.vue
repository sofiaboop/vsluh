<script setup lang="ts">
import { WHEEL } from '~/shared/constants/practice'
import { playTick } from '~/shared/lib/sound'

const { t, tm, rt } = useI18n()
const {
  rows,
  spinning,
  settled,
  selectedIndex,
  moodIndex,
  analysisEnabled,
  chosen,
  startPrepare,
  spin,
  finishSpin,
  highlight,
} = usePractice()

const mood = computed(() => {
  if (moodIndex.value === null) return ''
  const options = tm('mood.options') as unknown[]
  const option = options[moodIndex.value]
  return option ? rt(option as string) : ''
})

const listRef = ref<HTMLElement | null>(null)
const winnerRef = ref<HTMLElement | null>(null)
let raf = 0
let passed = WHEEL.startIndex

useFitText(winnerRef, () => [settled.value, chosen.value?.question])

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

function setWinnerRef(el: Element | ComponentPublicInstance | null) {
  winnerRef.value = el as HTMLElement | null
}
</script>

<template>
  <section class="content wheel-content">
    <div class="step-label">{{ t('steps.wheel') }}</div>

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
          <span class="wheel-category">{{ t(`categories.${row.topic}`) }}</span>
          <span
            :ref="index === WHEEL.selectedIndex ? setWinnerRef : undefined"
            class="wheel-topic"
          >{{ row.question }}</span>
        </div>
      </div>
    </div>

    <div v-if="chosen" class="chosen-mobile">
      <span>{{ t(`categories.${chosen.topic}`) }}</span>{{ chosen.question }}
    </div>

    <div class="wheel-statuses">
      <div v-if="mood" class="mood-pill"><span class="mood-dot" />{{ mood }}</div>
      <AnalysisStatus v-if="analysisEnabled" variant="planned" />
    </div>

    <div class="wheel-actions">
      <button
        type="button"
        class="action-button action-button--primary"
        :disabled="spinning"
        @click="startPrepare"
      >
        {{ t('wheel.start') }}
      </button>
      <button
        type="button"
        class="action-button action-button--secondary"
        :disabled="spinning"
        @click="spin"
      >
        <span class="spin-icon" :class="{ 'spin-icon--active': spinning }">↻</span>
        {{ spinning ? t('wheel.spinning') : t('wheel.spinAgain') }}
      </button>
    </div>
  </section>
</template>
