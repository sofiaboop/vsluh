<script setup lang="ts">
import { FEATURES } from '~/shared/constants/features'

const { t, tm, rt } = useI18n()
const {
  stage,
  mode,
  chosen,
  moodIndex,
  analysisEnabled,
  remaining,
  paused,
  timerAngle,
  startSpeak,
  finish,
  togglePause,
  backToTopics,
} = usePractice()

const isSpeak = computed(() => stage.value === 'speak')
const isDeep = computed(() => mode.value === 'deep')
const showAnalysis = computed(() => FEATURES.analysis && analysisEnabled.value)

const mood = computed(() => {
  if (moodIndex.value === null) return ''
  const option = (tm('mood.options') as unknown[])[moodIndex.value]
  return option ? rt(option as string) : ''
})

const topicRef = ref<HTMLElement | null>(null)
useFitText(topicRef, () => [chosen.value?.question, stage.value], { fillScreen: true })

const clock = computed(() => {
  const total = Math.max(0, remaining.value)
  const mm = String(Math.floor(total / 60)).padStart(2, '0')
  const ss = String(total % 60).padStart(2, '0')
  return `${mm}:${ss}`
})
</script>

<template>
  <section class="content timer-screen" :class="{ 'timer-screen--deep': isDeep }">
    <div class="step-label">{{ isSpeak ? t('steps.speak') : t('steps.prepare') }}</div>

    <p ref="topicRef" class="prepare-topic prepare-topic--deep">{{ chosen?.question }}</p>

    <p v-if="isDeep && !isSpeak" class="deep-prepare-description">
      {{ t('modes.deep.prepareNote') }}
    </p>

    <div class="timer-panel" :class="isSpeak ? 'timer-panel--speak' : 'timer-panel--prepare'">
      <div class="timer-content">
        <div class="timer-display">
          <div
            aria-hidden="true"
            class="timer-aura timer-aura--animated"
            :style="{ '--timer-angle': `${timerAngle}deg` }"
          />
          <div class="timer-display-content">
            <div class="timer-value" :aria-label="t('timer.remaining', { time: clock })">
              {{ clock }}
            </div>
            <div v-if="!isDeep || isSpeak" class="timer-caption">
              {{ isSpeak ? t('timer.speakCaption') : t('timer.prepareCaption') }}
            </div>
          </div>
        </div>
        <div v-if="!isDeep && (mood || showAnalysis)" class="timer-indicators">
          <div v-if="mood" class="mood-pill"><span class="mood-dot" />{{ mood }}</div>
          <AnalysisStatus v-if="showAnalysis" :variant="isSpeak ? 'recording' : 'planned'" />
        </div>
      </div>
    </div>

    <div v-if="isDeep && (mood || showAnalysis)" class="timer-footer-indicators">
      <div v-if="mood" class="mood-pill"><span class="mood-dot" />{{ mood }}</div>
      <AnalysisStatus v-if="showAnalysis" :variant="isSpeak ? 'recording' : 'planned'" />
    </div>

    <div class="timer-actions">
      <button
        type="button"
        class="action-button action-button--primary"
        @click="isSpeak ? finish() : startSpeak()"
      >
        {{ isSpeak ? t('timer.finish') : t('timer.speakNow') }}
      </button>
      <button type="button" class="action-button action-button--secondary" @click="togglePause">
        {{ paused ? t('timer.resume') : t('timer.pause') }}
      </button>
      <button type="button" class="action-button action-button--tertiary" @click="backToTopics">
        {{ t('timer.close') }}
      </button>
    </div>
  </section>
</template>
