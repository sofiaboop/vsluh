<script setup lang="ts">
import { ApiError } from '~/shared/api/apiInstance'
import { spendCredit } from '~/shared/api/auth'
import { ANALYSIS_DEMO } from '~/shared/constants/analysis'

const { t, tm, rt } = useI18n()
const { analysisEnabled, again, repeatTopic, backToTopics } = usePractice()
const { isAuthenticated, setCredits, openAuth, openBuy } = useAuth()
const features = useFeatures()

/** idle — разбор ещё не оплачен, empty — кредиты кончились. */
type AnalysisState = 'idle' | 'pending' | 'ready' | 'empty' | 'error'
const state = ref<AnalysisState>('idle')

const wantsAnalysis = computed(() => features.value.analysis && analysisEnabled.value)

/** Разбор стоит кредит: сначала списываем, и только потом показываем цифры. */
async function payForAnalysis() {
  if (state.value !== 'idle') return
  state.value = 'pending'
  try {
    setCredits(await spendCredit())
    state.value = 'ready'
  } catch (error) {
    state.value = error instanceof ApiError && error.status === 402 ? 'empty' : 'error'
  }
}

function retryAnalysis() {
  state.value = 'idle'
  void payForAnalysis()
}

// Вход мог случиться уже на этом экране — тогда разбор оплачивается сразу после него.
watch(
  [wantsAnalysis, isAuthenticated],
  ([wanted, authed]) => {
    if (wanted && authed) void payForAnalysis()
  },
  { immediate: true },
)

const fillerWords = computed(() =>
  (tm('analysis.demoWords') as unknown[]).map((word, index) => ({
    word: rt(word as string),
    count: ANALYSIS_DEMO.fillerWords[index] ?? 0,
  })),
)
</script>

<template>
  <section class="content finished-content">
    <div class="done-mark"><span /></div>

    <h1>{{ t('finished.title') }}</h1>
    <p class="lead">{{ t('finished.lead') }}</p>

    <div v-if="wantsAnalysis && state === 'ready'" class="analysis-summary">
      <div class="analysis-summary-heading">
        <span>{{ t('analysis.kicker') }}</span>
        <strong>{{ t('analysis.title') }}</strong>
      </div>

      <div class="analysis-summary-grid">
        <div class="analysis-summary-group">
          <span class="analysis-summary-label">{{ t('analysis.fillerPauses') }}</span>
          <div class="analysis-summary-row">
            <span>{{ t('analysis.demoPauses') }}</span>
            <strong>{{ ANALYSIS_DEMO.fillerPauses }}</strong>
          </div>
        </div>

        <div class="analysis-summary-group">
          <span class="analysis-summary-label">{{ t('analysis.fillerWords') }}</span>
          <div v-for="item in fillerWords" :key="item.word" class="analysis-summary-row">
            <span>{{ item.word }}</span>
            <strong>{{ item.count }}</strong>
          </div>
        </div>

        <div class="analysis-summary-metric">
          <span>{{ t('analysis.slips') }}</span>
          <strong>{{ ANALYSIS_DEMO.slips }}</strong>
        </div>

        <div class="analysis-summary-metric">
          <span>{{ t('analysis.repeats') }}</span>
          <strong>{{ ANALYSIS_DEMO.repeats }}</strong>
        </div>

        <div class="analysis-summary-metric">
          <span>{{ t('analysis.pace') }}</span>
          <strong>{{ t('analysis.paceValue') }}</strong>
        </div>
      </div>
    </div>

    <div v-else-if="wantsAnalysis && !isAuthenticated" class="analysis-auth-callout">
      <p>{{ t('analysis.authCallout') }}</p>
      <button type="button" @click="openAuth('login')">{{ t('analysis.authCta') }}</button>
    </div>

    <div v-else-if="wantsAnalysis && state === 'empty'" class="analysis-auth-callout">
      <strong>{{ t('credits.empty.title') }}</strong>
      <p>{{ t('credits.empty.text') }}</p>
      <button class="credits-buy" type="button" @click="openBuy">
        {{ t('credits.empty.cta') }}
      </button>
    </div>

    <div v-else-if="wantsAnalysis && state === 'error'" class="analysis-auth-callout">
      <p>{{ t('credits.error') }}</p>
      <button type="button" @click="retryAnalysis">
        {{ t('credits.retry') }}
      </button>
    </div>

    <p v-else-if="wantsAnalysis && state === 'pending'" class="analysis-pending" role="status">
      {{ t('credits.checking') }}
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
