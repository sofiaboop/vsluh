<script setup lang="ts">
import { ANALYSIS_DEMO } from '~/shared/constants/analysis'
import { FEATURES } from '~/shared/constants/features'

const { t, tm, rt } = useI18n()
const { analysisEnabled, again, repeatTopic, backToTopics } = usePractice()
const { isAuthenticated, openAuth } = useAuth()

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

    <div v-if="FEATURES.analysis && analysisEnabled" class="analysis-summary">
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

      <div v-if="FEATURES.auth && !isAuthenticated" class="analysis-auth-callout">
        <p>{{ t('analysis.authCallout') }}</p>
        <button type="button" @click="openAuth('login')">{{ t('analysis.authCta') }}</button>
      </div>
    </div>

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
