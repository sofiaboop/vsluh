<script setup lang="ts">
import {
  fetchCheck,
  fetchChecks,
  fetchStats,
  type Check,
  type Stats,
} from "~/shared/api/checks";
import { verdictFor } from "~/shared/lib/verdict";

const { t, locale } = useI18n();
const { isAuthenticated, openAuth } = useAuth();
const features = useFeatures();

const stats = ref<Stats | null>(null);
const checks = ref<Check[]>([]);
const cursor = ref("");
const loading = ref(false);
const failed = ref(false);
/** Проверка, открытая в полном отчёте: детали просим отдельным запросом. */
const openCheck = ref<Check | null>(null);
const openPending = ref("");
const openFailed = ref(false);

const hasMore = computed(() => cursor.value !== "");

/** Вердикт считаем один раз на проверку, а не на каждое обращение из шаблона. */
const cards = computed(() =>
  checks.value.map((check) => ({ check, verdict: verdictFor(check.metrics) })),
);

async function load(reset = false) {
  if (loading.value) return;
  loading.value = true;
  failed.value = false;

  try {
    if (reset) {
      cursor.value = "";
      checks.value = [];
      stats.value = await fetchStats();
    }
    const page = await fetchChecks(cursor.value);
    checks.value = [...checks.value, ...page.checks];
    cursor.value = page.nextCursor;
  } catch {
    failed.value = true;
  } finally {
    loading.value = false;
  }
}

async function openReport(check: Check) {
  openPending.value = check.id;
  openFailed.value = false;
  try {
    openCheck.value = await fetchCheck(check.id);
  } catch {
    openFailed.value = true;
  } finally {
    openPending.value = "";
  }
}

function closeReport() {
  openCheck.value = null;
}

watch(
  isAuthenticated,
  (authed) => {
    if (authed) void load(true);
    else {
      checks.value = [];
      stats.value = null;
      openCheck.value = null;
    }
  },
  { immediate: true },
);

const dateFormat = computed(
  () =>
    new Intl.DateTimeFormat(locale.value, {
      day: "numeric",
      month: "short",
      hour: "2-digit",
      minute: "2-digit",
    }),
);

function formatDate(seconds: number) {
  return dateFormat.value.format(new Date(seconds * 1000));
}

function formatMinutes(seconds: number) {
  const minutes = seconds / 60;
  // Первые практики короче минуты: округление до целого показало бы ноль.
  return minutes < 10 ? Math.round(minutes * 10) / 10 : Math.round(minutes);
}

/** Высота столбика графика в процентах от самого длинного дня. */
function barHeight(value: number) {
  const peak = Math.max(
    ...(stats.value?.daily ?? []).map((d) => d.speechSec),
    1,
  );
  return Math.max(6, Math.round((value / peak) * 100));
}
</script>

<template>
  <div class="screen screen--history">
    <section class="content history-content">
      <h1>{{ t("history.title") }}</h1>
      <p class="lead">{{ t("history.lead") }}</p>

      <div
        v-if="!features.auth || !isAuthenticated"
        class="analysis-auth-callout"
      >
        <p>{{ t("history.authCallout") }}</p>
        <button v-if="features.auth" type="button" @click="openAuth('login')">
          {{ t("analysis.authCta") }}
        </button>
      </div>

      <template v-else>
        <div v-if="stats" class="overview">
          <div class="overview-head">
            <span>{{ t("history.stats.overview") }}</span>
            <strong>{{ stats.totalChecks }}</strong>
            <p>{{ t("history.stats.overviewHint") }}</p>
          </div>

          <div class="stats-grid">
            <div class="stats-card">
              <span>{{ t("history.stats.minutes") }}</span>
              <strong>{{ formatMinutes(stats.totalSpeechSec) }}</strong>
            </div>
            <div class="stats-card">
              <span>{{ t("history.stats.words") }}</span>
              <strong>{{ stats.totalWords }}</strong>
            </div>
            <div class="stats-card">
              <span>{{ t("history.stats.pace") }}</span>
              <strong>{{ stats.avgWordsPerMin }}</strong>
            </div>
            <div class="stats-card">
              <span>{{ t("history.stats.streak") }}</span>
              <strong>{{ stats.streakDays }}</strong>
            </div>
            <div class="stats-card">
              <span>{{ t("history.stats.fillersTotal") }}</span>
              <strong>{{ stats.totalFillers }}</strong>
            </div>
            <div class="stats-card">
              <span>{{ t("history.stats.fillersAvg") }}</span>
              <strong>{{ stats.avgFillersPerCheck }}</strong>
            </div>
          </div>
        </div>

        <div v-if="stats && stats.daily.length" class="stats-chart">
          <span class="analysis-summary-label">{{
            t("history.stats.dynamics")
          }}</span>
          <div class="stats-chart-bars">
            <div
              v-for="point in stats.daily"
              :key="point.date"
              class="stats-chart-bar"
              :style="{ height: barHeight(point.speechSec) + '%' }"
              :title="
                t('history.stats.barHint', {
                  date: point.date,
                  checks: point.checks,
                })
              "
            />
          </div>
        </div>

        <div v-if="stats && stats.topFillers.length" class="stats-fillers">
          <span class="analysis-summary-label">{{
            t("history.stats.fillers")
          }}</span>
          <div
            v-for="item in stats.topFillers"
            :key="item.word"
            class="analysis-summary-row"
          >
            <span>{{ item.word }}</span>
            <strong>{{ item.count }}</strong>
          </div>
        </div>

        <p v-if="failed" class="analysis-pending" role="status">
          {{ t("history.failed") }}
        </p>

        <p v-else-if="!loading && !checks.length" class="analysis-pending">
          {{ t("history.empty") }}
        </p>

        <ul v-if="cards.length" class="history-list">
          <li v-for="card in cards" :key="card.check.id" class="history-item">
            <div class="history-item-head">
              <span class="history-item-date">{{
                formatDate(card.check.createdAt)
              }}</span>
              <span class="history-item-topic">
                {{ card.check.topicLabel || t("history.noTopic") }}
              </span>
              <span class="history-item-meta">
                {{
                  t("report.wpmValue", {
                    value: card.check.metrics.wordsPerMin,
                  })
                }}
                ·
                {{
                  t("history.fillersShort", {
                    count: card.check.metrics.fillerCount,
                  })
                }}
              </span>
            </div>

            <div class="verdict" :class="`verdict--${card.verdict.level}`">
              <div class="verdict-head">
                <strong>{{
                  t(`history.verdict.${card.verdict.level}`)
                }}</strong>
                <span class="verdict-score">{{
                  t("history.verdict.score", { score: card.verdict.score })
                }}</span>
              </div>
              <p>{{ t(`history.verdict.focus.${card.verdict.focus}`) }}</p>
            </div>

            <audio
              v-if="card.check.audioUrl"
              class="analysis-audio"
              :src="card.check.audioUrl"
              controls
              preload="none"
            />

            <button
              class="history-item-more"
              type="button"
              @click="openReport(card.check)"
            >
              {{
                openPending === card.check.id
                  ? t("history.loading")
                  : t("history.details")
              }}
            </button>
          </li>
        </ul>

        <p v-if="openFailed" class="analysis-pending" role="status">
          {{ t("history.failed") }}
        </p>

        <button
          v-if="hasMore"
          class="action-button action-button--secondary"
          type="button"
          @click="load()"
        >
          {{ loading ? t("history.loading") : t("history.more") }}
        </button>
      </template>

      <NuxtLink class="analysis-history-link" to="/">{{
        t("history.back")
      }}</NuxtLink>
    </section>

    <CheckReportDialog
      v-if="openCheck"
      :check="openCheck"
      @close="closeReport"
    />
  </div>
</template>
