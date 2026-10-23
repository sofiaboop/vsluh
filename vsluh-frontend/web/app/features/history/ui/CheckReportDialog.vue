<script setup lang="ts">
import type { Check } from "~/shared/api/checks";
import { verdictFor } from "~/shared/lib/verdict";

const props = defineProps<{ check: Check }>();
const emit = defineEmits<{ close: [] }>();

const { t, locale } = useI18n();

const verdict = computed(() => verdictFor(props.check.metrics));
/** Формата записи в ответе нет — берём расширение из ключа подписанной ссылки. */
const downloadName = computed(() => {
  const path = props.check.audioUrl.split("?")[0] ?? "";
  const ext = path.includes(".")
    ? path.slice(path.lastIndexOf(".") + 1)
    : "webm";
  return `vsluh-${props.check.id}.${ext}`;
});
const createdAt = computed(() =>
  new Intl.DateTimeFormat(locale.value, {
    day: "numeric",
    month: "long",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(props.check.createdAt * 1000)),
);

useEventListener(window, "keydown", (e: KeyboardEvent) => {
  if (e.key === "Escape") emit("close");
});
</script>

<template>
  <Teleport to="body">
    <div
      :aria-label="t('history.report.title')"
      aria-modal="true"
      class="buy-overlay"
      role="dialog"
      @mousedown.self="emit('close')"
    >
      <div class="buy-dialog report-dialog">
        <button
          :aria-label="t('auth.close')"
          class="support-close"
          type="button"
          @click="emit('close')"
        />

        <span class="history-item-date">{{ createdAt }}</span>
        <h2>{{ check.topicLabel || t("history.noTopic") }}</h2>

        <div class="verdict" :class="`verdict--${verdict.level}`">
          <div class="verdict-head">
            <strong>{{ t(`history.verdict.${verdict.level}`) }}</strong>
            <span class="verdict-score">{{
              t("history.verdict.score", { score: verdict.score })
            }}</span>
          </div>
          <p>{{ t(`history.verdict.focus.${verdict.focus}`) }}</p>
        </div>

        <audio
          v-if="check.audioUrl"
          class="analysis-audio"
          :src="check.audioUrl"
          controls
          preload="none"
        />
        <a
          v-if="check.audioUrl"
          class="analysis-download"
          :download="downloadName"
          :href="check.audioUrl"
          >{{ t("report.download") }}</a
        >

        <p v-if="check.transcript" class="analysis-transcript">
          {{ check.transcript }}
        </p>

        <div class="analysis-summary-grid">
          <div class="analysis-summary-metric">
            <span>{{ t("analysis.pace") }}</span>
            <strong>{{
              t("report.wpmValue", { value: check.metrics.wordsPerMin })
            }}</strong>
          </div>
          <div class="analysis-summary-metric">
            <span>{{ t("report.words") }}</span>
            <strong>{{ check.metrics.wordCount }}</strong>
          </div>
          <div class="analysis-summary-metric">
            <span>{{ t("analysis.fillerPauses") }}</span>
            <strong>{{ check.metrics.longPauses }}</strong>
          </div>
          <div class="analysis-summary-metric">
            <span>{{ t("analysis.repeats") }}</span>
            <strong>{{ check.metrics.repeatCount }}</strong>
          </div>

          <div class="analysis-summary-group">
            <span class="analysis-summary-label">{{
              t("analysis.fillerWords")
            }}</span>
            <div
              v-if="!check.metrics.fillerWords.length"
              class="analysis-summary-row"
            >
              <span>{{ t("report.noFillers") }}</span>
              <strong>0</strong>
            </div>
            <div
              v-for="item in check.metrics.fillerWords"
              :key="item.word"
              class="analysis-summary-row"
            >
              <span>{{ item.word }}</span>
              <strong>{{ item.count }}</strong>
            </div>
          </div>
        </div>

        <p v-if="check.advice" class="analysis-advice">{{ check.advice }}</p>
        <ul v-else-if="check.tips.length" class="analysis-tips">
          <li v-for="tip in check.tips" :key="tip">
            {{ t(`report.tips.${tip}`) }}
          </li>
        </ul>
      </div>
    </div>
  </Teleport>
</template>
