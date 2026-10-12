<script setup lang="ts">
import { MODES, TOPIC_IDS } from '~/shared/constants/practice'

const { t } = useI18n()
const { mode, selected, moodEnabled, analysisEnabled, allSelected, canContinue, toggleTopic, toggleAll, startWheel } =
  usePractice()

const tooltipOpen = ref(false)
const analysisTooltipOpen = ref(false)
</script>

<template>
  <section class="content selection-content">
    <div class="selection-meta">
      <div class="selection-mode">
        <template v-for="(option, index) in MODES" :key="option.id">
          <i v-if="index > 0">/</i>
          <button
            type="button"
            class="selection-mode-option"
            :class="{ 'selection-mode-option--active': mode === option.id }"
            @click="mode = option.id"
          >
            {{ t(`modes.${option.id}.label`) }}
          </button>
        </template>
      </div>
    </div>

    <div class="selection-duration">{{ t(`modes.${mode}.duration`) }}</div>

    <div class="selection-main">
      <h1>{{ t(`modes.${mode}.heading`) }}</h1>
      <p class="lead">{{ t('selection.lead') }}</p>

      <div class="chips">
        <button
          type="button"
          class="chip chip--special"
          :class="{ 'chip--selected': allSelected }"
          :aria-pressed="allSelected"
          @click="toggleAll"
        >
          {{ allSelected ? t('selection.clearAll') : t('selection.selectAll') }}
        </button>
        <button
          v-for="id in TOPIC_IDS"
          :key="id"
          type="button"
          class="chip"
          :class="{ 'chip--selected': selected.includes(id) }"
          :aria-pressed="selected.includes(id)"
          @click="toggleTopic(id)"
        >
          {{ t(`categories.${id}`) }}
        </button>
      </div>

      <div class="mood-control">
        <button
          type="button"
          class="toggle"
          :class="{ 'toggle--active': moodEnabled }"
          role="switch"
          :aria-checked="moodEnabled"
          :aria-label="t('mood.toggle')"
          @click="moodEnabled = !moodEnabled"
        >
          <span />
        </button>
        <span>{{ t('mood.toggle') }}</span>
        <button
          type="button"
          class="help-button"
          :aria-label="t('mood.help')"
          :aria-expanded="tooltipOpen"
          @click="tooltipOpen = !tooltipOpen"
        >
          ?
        </button>
        <div v-if="tooltipOpen" class="tooltip" role="tooltip">{{ t('mood.tooltip') }}</div>
      </div>

      <div class="mood-control analysis-control">
        <button
          type="button"
          class="toggle"
          :class="{ 'toggle--active': analysisEnabled }"
          role="switch"
          :aria-checked="analysisEnabled"
          :aria-label="t('analysis.toggle')"
          @click="analysisEnabled = !analysisEnabled"
        >
          <span />
        </button>
        <span>{{ t('analysis.toggle') }}</span>
        <button
          type="button"
          class="help-button"
          :aria-label="t('analysis.help')"
          :aria-expanded="analysisTooltipOpen"
          @click="analysisTooltipOpen = !analysisTooltipOpen"
        >
          ?
        </button>
        <div v-if="analysisTooltipOpen" class="tooltip tooltip--analysis" role="tooltip">
          <strong>{{ t('analysis.tooltipTitle') }}</strong>
          <span>{{ t('analysis.tooltipText') }}</span>
        </div>
      </div>

      <div class="selection-action">
        <button
          type="button"
          class="action-button action-button--primary"
          :disabled="!canContinue"
          @click="startWheel"
        >
          {{ t('selection.continue') }}
        </button>
      </div>
    </div>
  </section>
</template>
