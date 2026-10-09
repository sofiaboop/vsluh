<script setup lang="ts">
const {
  stage,
  chosen,
  mood,
  remaining,
  paused,
  timerAngle,
  startSpeak,
  finish,
  togglePause,
  backToTopics,
} = usePractice()

const isSpeak = computed(() => stage.value === 'speak')

const clock = computed(() => {
  const total = Math.max(0, remaining.value)
  const mm = String(Math.floor(total / 60)).padStart(2, '0')
  const ss = String(total % 60).padStart(2, '0')
  return `${mm}:${ss}`
})
</script>

<template>
  <section class="content timer-screen">
    <div class="step-label">{{ isSpeak ? '04 / ГОВОРИ' : '03 / ПОДГОТОВКА' }}</div>

    <p class="prepare-topic">{{ chosen?.question }}</p>

    <div class="timer-panel" :class="isSpeak ? 'timer-panel--speak' : 'timer-panel--prepare'">
      <div class="timer-content">
        <div class="timer-display">
          <div
            aria-hidden="true"
            class="timer-aura timer-aura--animated"
            :style="{ '--timer-angle': `${timerAngle}deg` }"
          />
          <div class="timer-display-content">
            <div class="timer-value" :aria-label="`Осталось ${clock}`">{{ clock }}</div>
            <div class="timer-caption">
              {{ isSpeak ? 'Говори' : 'Подумай, что хочешь сказать' }}
            </div>
          </div>
        </div>
        <div v-if="mood" class="mood-pill"><span class="mood-dot" />{{ mood }}</div>
      </div>
    </div>

    <div class="timer-actions">
      <button
        type="button"
        class="action-button action-button--primary"
        @click="isSpeak ? finish() : startSpeak()"
      >
        {{ isSpeak ? 'Закончить' : 'Говорить сейчас' }}
      </button>
      <button type="button" class="action-button action-button--secondary" @click="togglePause">
        {{ paused ? 'Продолжить' : 'Пауза' }}
      </button>
      <button type="button" class="action-button action-button--tertiary" @click="backToTopics">
        Закрыть
      </button>
    </div>
  </section>
</template>
