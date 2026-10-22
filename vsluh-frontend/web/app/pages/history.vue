<script setup lang="ts">
import { fetchCheck, fetchChecks, fetchStats, type Check, type Stats } from '~/shared/api/checks'

const { t, locale } = useI18n()
const { isAuthenticated, openAuth } = useAuth()
const features = useFeatures()

const stats = ref<Stats | null>(null)
const checks = ref<Check[]>([])
const cursor = ref('')
const loading = ref(false)
const failed = ref(false)
/** Открытая карточка истории: ссылку на запись просим отдельно, она временная. */
const openId = ref('')
const openCheck = ref<Check | null>(null)
const openFailed = ref(false)

const hasMore = computed(() => cursor.value !== '')

async function load(reset = false) {
  if (loading.value) return
  loading.value = true
  failed.value = false

  try {
    if (reset) {
      cursor.value = ''
      checks.value = []
      stats.value = await fetchStats()
    }
    const page = await fetchChecks(cursor.value)
    checks.value = [...checks.value, ...page.checks]
    cursor.value = page.nextCursor
  } catch {
    failed.value = true
  } finally {
    loading.value = false
  }
}

async function toggle(check: Check) {
  if (openId.value === check.id) {
    openId.value = ''
    openCheck.value = null
    return
  }

  openId.value = check.id
  openCheck.value = null
  openFailed.value = false
  try {
    openCheck.value = await fetchCheck(check.id)
  } catch {
    openFailed.value = true
  }
}

watch(
  isAuthenticated,
  (authed) => {
    if (authed) void load(true)
    else {
      checks.value = []
      stats.value = null
    }
  },
  { immediate: true },
)

const dateFormat = computed(
  () =>
    new Intl.DateTimeFormat(locale.value, {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    }),
)

function formatDate(seconds: number) {
  return dateFormat.value.format(new Date(seconds * 1000))
}

function formatMinutes(seconds: number) {
  const minutes = seconds / 60
  // Первые практики короче минуты: округление до целого показало бы ноль.
  return minutes < 10 ? Math.round(minutes * 10) / 10 : Math.round(minutes)
}

/** Высота столбика графика в процентах от самого длинного дня. */
function barHeight(value: number) {
  const peak = Math.max(...(stats.value?.daily ?? []).map((d) => d.speechSec), 1)
  return Math.max(6, Math.round((value / peak) * 100))
}
</script>

<template>
  <div class="screen screen--history">
    <section class="content history-content">
      <h1>{{ t('history.title') }}</h1>
      <p class="lead">{{ t('history.lead') }}</p>

      <div v-if="!features.auth || !isAuthenticated" class="analysis-auth-callout">
        <p>{{ t('history.authCallout') }}</p>
        <button v-if="features.auth" type="button" @click="openAuth('login')">
          {{ t('analysis.authCta') }}
        </button>
      </div>

      <template v-else>
        <div v-if="stats" class="stats-grid">
          <div class="stats-card">
            <span>{{ t('history.stats.checks') }}</span>
            <strong>{{ stats.totalChecks }}</strong>
          </div>
          <div class="stats-card">
            <span>{{ t('history.stats.minutes') }}</span>
            <strong>{{ formatMinutes(stats.totalSpeechSec) }}</strong>
          </div>
          <div class="stats-card">
            <span>{{ t('history.stats.pace') }}</span>
            <strong>{{ stats.avgWordsPerMin }}</strong>
          </div>
          <div class="stats-card">
            <span>{{ t('history.stats.streak') }}</span>
            <strong>{{ stats.streakDays }}</strong>
          </div>
        </div>

        <div v-if="stats && stats.daily.length" class="stats-chart">
          <span class="analysis-summary-label">{{ t('history.stats.dynamics') }}</span>
          <div class="stats-chart-bars">
            <div
              v-for="point in stats.daily"
              :key="point.date"
              class="stats-chart-bar"
              :style="{ height: barHeight(point.speechSec) + '%' }"
              :title="t('history.stats.barHint', { date: point.date, checks: point.checks })"
            />
          </div>
        </div>

        <div v-if="stats && stats.topFillers.length" class="stats-fillers">
          <span class="analysis-summary-label">{{ t('history.stats.fillers') }}</span>
          <div v-for="item in stats.topFillers" :key="item.word" class="analysis-summary-row">
            <span>{{ item.word }}</span>
            <strong>{{ item.count }}</strong>
          </div>
        </div>

        <p v-if="failed" class="analysis-pending" role="status">{{ t('history.failed') }}</p>

        <p v-else-if="!loading && !checks.length" class="analysis-pending">
          {{ t('history.empty') }}
        </p>

        <ul v-if="checks.length" class="history-list">
          <li v-for="check in checks" :key="check.id" class="history-item">
            <button class="history-item-head" type="button" @click="toggle(check)">
              <span class="history-item-date">{{ formatDate(check.createdAt) }}</span>
              <span class="history-item-topic">{{ check.topicLabel || t('history.noTopic') }}</span>
              <span class="history-item-meta">
                {{ t('report.wpmValue', { value: check.metrics.wordsPerMin }) }} ·
                {{ t('history.fillersShort', { count: check.metrics.fillerCount }) }}
              </span>
            </button>

            <div v-if="openId === check.id" class="history-item-body">
              <p v-if="openFailed" class="analysis-pending">{{ t('history.failed') }}</p>
              <p v-else-if="!openCheck" class="analysis-pending" role="status">
                {{ t('history.loading') }}
              </p>
              <template v-else>
                <a
                  v-if="openCheck.audioUrl"
                  class="analysis-download"
                  :href="openCheck.audioUrl"
                  download
                  >{{ t('report.download') }}</a
                >
                <audio
                  v-if="openCheck.audioUrl"
                  class="analysis-audio"
                  :src="openCheck.audioUrl"
                  controls
                  preload="none"
                />
                <p v-if="openCheck.transcript" class="analysis-transcript">
                  {{ openCheck.transcript }}
                </p>
                <div class="analysis-summary-grid">
                  <div class="analysis-summary-metric">
                    <span>{{ t('report.words') }}</span>
                    <strong>{{ openCheck.metrics.wordCount }}</strong>
                  </div>
                  <div class="analysis-summary-metric">
                    <span>{{ t('analysis.fillerPauses') }}</span>
                    <strong>{{ openCheck.metrics.longPauses }}</strong>
                  </div>
                  <div class="analysis-summary-metric">
                    <span>{{ t('analysis.repeats') }}</span>
                    <strong>{{ openCheck.metrics.repeatCount }}</strong>
                  </div>
                </div>
                <p v-if="openCheck.advice" class="analysis-advice">{{ openCheck.advice }}</p>
                <ul v-else-if="openCheck.tips.length" class="analysis-tips">
                  <li v-for="tip in openCheck.tips" :key="tip">{{ t(`report.tips.${tip}`) }}</li>
                </ul>
              </template>
            </div>
          </li>
        </ul>

        <button v-if="hasMore" class="action-button action-button--secondary" type="button" @click="load()">
          {{ loading ? t('history.loading') : t('history.more') }}
        </button>
      </template>

      <NuxtLink class="analysis-history-link" to="/">{{ t('history.back') }}</NuxtLink>
    </section>
  </div>
</template>
