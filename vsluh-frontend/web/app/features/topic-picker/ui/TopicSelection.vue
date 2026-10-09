<script setup lang="ts">
import { MODES, MOOD_TOOLTIP, TOPICS } from '~/shared/constants/practice'

const { mode, modeOption, selected, moodEnabled, allSelected, canContinue, toggleTopic, toggleAll, startWheel } =
  usePractice()

const tooltipOpen = ref(false)
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
            {{ option.label }}
          </button>
        </template>
      </div>
    </div>

    <div class="selection-duration">{{ modeOption.duration }}</div>

    <div class="selection-main">
      <h1>О чём поговорим?</h1>
      <p class="lead">Выбери одну или несколько тем</p>

      <div class="chips">
        <button
          type="button"
          class="chip chip--special"
          :class="{ 'chip--selected': allSelected }"
          :aria-pressed="allSelected"
          @click="toggleAll"
        >
          Выбрать все
        </button>
        <button
          v-for="topic in TOPICS"
          :key="topic.id"
          type="button"
          class="chip"
          :class="{ 'chip--selected': selected.includes(topic.id) }"
          :aria-pressed="selected.includes(topic.id)"
          @click="toggleTopic(topic.id)"
        >
          {{ topic.label }}
        </button>
      </div>

      <div class="mood-control">
        <button
          type="button"
          class="toggle"
          :class="{ 'toggle--active': moodEnabled }"
          role="switch"
          :aria-checked="moodEnabled"
          aria-label="Добавить настроение"
          @click="moodEnabled = !moodEnabled"
        >
          <span />
        </button>
        <span>Добавить настроение</span>
        <button
          type="button"
          class="help-button"
          aria-label="Что значит добавить настроение"
          :aria-expanded="tooltipOpen"
          @click="tooltipOpen = !tooltipOpen"
        >
          ?
        </button>
        <div v-if="tooltipOpen" class="tooltip" role="tooltip">{{ MOOD_TOOLTIP }}</div>
      </div>

      <div class="selection-action">
        <button
          type="button"
          class="action-button action-button--primary"
          :disabled="!canContinue"
          @click="startWheel"
        >
          Продолжить
        </button>
      </div>
    </div>
  </section>
</template>
