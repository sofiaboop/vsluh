<script setup lang="ts">
import { ApiError } from '~/shared/api/apiInstance'
import { createCheck, type Check } from '~/shared/api/checks'
import { toBase64 } from '~/shared/lib/audioRecorder'

const { t } = useI18n()
const { locale } = useI18n()
const { analysisWanted, analysisEnabled, chosen, mode, modeOption, again, repeatTopic, backToTopics } =
  usePractice()
const { isAuthenticated, setCredits, openAuth, openBuy } = useAuth()
const { state: recorderState, take } = useRecorder()
const features = useFeatures()

/** idle — разбор ещё не запускали, silent — в записи не нашлось речи. */
type AnalysisState = 'idle' | 'pending' | 'ready' | 'empty' | 'silent' | 'error'
const state = ref<AnalysisState>('idle')
const check = ref<Check | null>(null)
const audioUrl = ref('')
const audioFormat = ref('webm')

const needsAnalysis = computed(() => features.value.analysis && analysisEnabled.value)
const micBlocked = computed(
  () => recorderState.value === 'denied' || recorderState.value === 'unsupported',
)

async function runAnalysis() {
  if (state.value === 'pending') return

  const recording = await take()
  if (!recording || recording.blob.size === 0) {
    state.value = micBlocked.value ? 'idle' : 'silent'
    return
  }

  state.value = 'pending'
  // Запись показываем сразу: ссылка из хранилища придёт только вместе с разбором.
  audioUrl.value = URL.createObjectURL(recording.blob)
  audioFormat.value = recording.format

  try {
    const result = await createCheck({
      audioBase64: await toBase64(recording.blob),
      format: recording.format,
      topicId: chosen.value?.topic ?? '',
      topicLabel: chosen.value?.question ?? '',
      mode: mode.value,
      speakSec: modeOption.value.speakSec,
      language: locale.value,
    })
    check.value = result.check
    setCredits(result.credits)
    state.value = 'ready'
  } catch (error) {
    if (error instanceof ApiError && error.status === 402) {
      state.value = 'empty'
      return
    }
    if (error instanceof ApiError && error.status === 422) {
      state.value = 'silent'
      return
    }
    state.value = 'error'
  }
}

watch(
  [analysisWanted, () => recorderState.value],
  ([wanted]) => {
    if (wanted && state.value === 'idle') void runAnalysis()
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  if (audioUrl.value) URL.revokeObjectURL(audioUrl.value)
})

const metrics = computed(() => check.value?.metrics ?? null)
const downloadName = computed(() => `vsluh-${check.value?.id ?? 'answer'}.${audioFormat.value}`)
</script>

<template>
  <section class="content finished-content">
    <div class="done-mark"><span /></div>

    <h1>{{ t('finished.title') }}</h1>
    <p class="lead">{{ t('finished.lead') }}</p>

    <div v-if="needsAnalysis && state === 'ready' && check && metrics" class="analysis-summary">
      <div class="analysis-summary-heading">
        <span>{{ t('analysis.kicker') }}</span>
        <strong>{{ t('analysis.title') }}</strong>
      </div>

      <audio v-if="audioUrl" class="analysis-audio" :src="audioUrl" controls preload="metadata" />
      <a v-if="audioUrl" class="analysis-download" :href="audioUrl" :download="downloadName">{{
        t('report.download')
      }}</a>

      <p v-if="check.transcript" class="analysis-transcript">{{ check.transcript }}</p>

      <div class="analysis-summary-grid">
        <div class="analysis-summary-metric">
          <span>{{ t('analysis.pace') }}</span>
          <strong>{{ t('report.wpmValue', { value: metrics.wordsPerMin }) }}</strong>
        </div>

        <div class="analysis-summary-metric">
          <span>{{ t('report.words') }}</span>
          <strong>{{ metrics.wordCount }}</strong>
        </div>

        <div class="analysis-summary-metric">
          <span>{{ t('analysis.fillerPauses') }}</span>
          <strong>{{ metrics.longPauses }}</strong>
        </div>

        <div class="analysis-summary-metric">
          <span>{{ t('analysis.repeats') }}</span>
          <strong>{{ metrics.repeatCount }}</strong>
        </div>

        <div class="analysis-summary-group">
          <span class="analysis-summary-label">{{ t('analysis.fillerWords') }}</span>
          <div v-if="!metrics.fillerWords.length" class="analysis-summary-row">
            <span>{{ t('report.noFillers') }}</span>
            <strong>0</strong>
          </div>
          <div v-for="item in metrics.fillerWords" :key="item.word" class="analysis-summary-row">
            <span>{{ item.word }}</span>
            <strong>{{ item.count }}</strong>
          </div>
        </div>
      </div>

      <p v-if="check.advice" class="analysis-advice">{{ check.advice }}</p>
      <ul v-else-if="check.tips.length" class="analysis-tips">
        <li v-for="tip in check.tips" :key="tip">{{ t(`report.tips.${tip}`) }}</li>
      </ul>

      <NuxtLink class="analysis-history-link" to="/history">{{ t('report.historyLink') }}</NuxtLink>
    </div>

    <div v-else-if="needsAnalysis && !isAuthenticated" class="analysis-auth-callout">
      <p>{{ t('analysis.authCallout') }}</p>
      <button type="button" @click="openAuth('login')">{{ t('analysis.authCta') }}</button>
    </div>

    <div v-else-if="needsAnalysis && micBlocked" class="analysis-auth-callout">
      <strong>{{ t('report.micTitle') }}</strong>
      <p>{{ t('report.micText') }}</p>
    </div>

    <div v-else-if="needsAnalysis && state === 'empty'" class="analysis-auth-callout">
      <strong>{{ t('credits.empty.title') }}</strong>
      <p>{{ t('credits.empty.text') }}</p>
      <button class="credits-buy" type="button" @click="openBuy">
        {{ t('credits.empty.cta') }}
      </button>
    </div>

    <div v-else-if="needsAnalysis && state === 'silent'" class="analysis-auth-callout">
      <strong>{{ t('report.silentTitle') }}</strong>
      <p>{{ t('report.silentText') }}</p>
    </div>

    <div v-else-if="needsAnalysis && state === 'error'" class="analysis-auth-callout">
      <p>{{ t('credits.error') }}</p>
      <button type="button" @click="runAnalysis">{{ t('credits.retry') }}</button>
    </div>

    <p v-else-if="needsAnalysis && state === 'pending'" class="analysis-pending" role="status">
      {{ t('report.pending') }}
    </p>

    <div class="finished-actions">
      <button type="button" class="action-button action-button--primary" @click="again">
        {{ t('finished.again') }}
      </button>
      <button type="button" class="action-button action-button--secondary" @click="repeatTopic">
        {{ t('finished.repeat') }}
      </button>
      <button type="button" class="action-button action-button--tertiary" @click="backToTopics">
        {{ t('finished.changeTopics') }}
      </button>
    </div>
  </section>
</template>
