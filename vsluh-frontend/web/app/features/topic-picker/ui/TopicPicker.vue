<script setup lang="ts">
import { MODES, MOOD_HINT, TOPICS, type SpeakingMode } from '~/shared/constants/practice'

const mode = ref<SpeakingMode>('spontaneous')
const selected = ref<Set<string>>(new Set())
const moodEnabled = ref(false)
const hintOpen = ref(false)
const notice = ref(false)

const activeHint = computed(() => MODES.find((m) => m.id === mode.value)?.hint ?? '')
const allSelected = computed(() => selected.value.size === TOPICS.length)
const canContinue = computed(() => selected.value.size > 0)

function toggleTopic(id: string) {
  const next = new Set(selected.value)
  next.has(id) ? next.delete(id) : next.add(id)
  selected.value = next
  notice.value = false
}

function toggleAll() {
  selected.value = allSelected.value ? new Set() : new Set(TOPICS.map((t) => t.id))
  notice.value = false
}

// MVP: практики ещё нет, показываем заглушку вместо перехода
function onContinue() {
  if (!canContinue.value) return
  notice.value = true
}
</script>

<template>
  <section class="mx-auto flex w-full max-w-stage flex-col items-center px-5 pb-6 text-center sm:px-8">
    <div class="flex items-center gap-2.5 text-[10px] font-medium text-ink-muted animate-fade-up">
      <template v-for="(option, index) in MODES" :key="option.id">
        <span v-if="index > 0" class="text-ink-soft/60">/</span>
        <button
          type="button"
          class="pb-0.5 transition-colors duration-200"
          :class="
            mode === option.id
              ? 'border-b border-ink-muted text-ink-muted'
              : 'border-b border-transparent text-ink-soft hover:text-ink-muted'
          "
          :aria-pressed="mode === option.id"
          @click="mode = option.id"
        >
          {{ option.label }}
        </button>
      </template>
    </div>

    <p class="mt-5 text-[13px] text-ink-soft animate-fade-up [animation-delay:60ms]">
      {{ activeHint }}
    </p>

    <h1
      class="mt-6 text-display font-medium tracking-display text-ink animate-fade-up [animation-delay:120ms]"
    >
      О чём поговорим?
    </h1>

    <p class="mt-3 text-[17px] text-ink-soft animate-fade-up [animation-delay:180ms]">
      Выбери одну или несколько тем
    </p>

    <div
      class="mt-8 flex flex-wrap items-center justify-center gap-3 animate-fade-up [animation-delay:240ms]"
    >
      <BaseChip label="Выбрать все" :active="allSelected" @toggle="toggleAll" />
      <BaseChip
        v-for="topic in TOPICS"
        :key="topic.id"
        :label="topic.label"
        :active="selected.has(topic.id)"
        @toggle="toggleTopic(topic.id)"
      />
    </div>

    <div
      class="relative mt-8 flex items-center gap-3 animate-fade-up [animation-delay:300ms]"
      @mouseleave="hintOpen = false"
    >
      <BaseToggle v-model="moodEnabled" label="Добавить настроение" />
      <span class="text-sm text-ink">Добавить настроение</span>
      <button
        type="button"
        class="flex h-6 w-6 items-center justify-center rounded-full border border-chip
               bg-chip text-xs text-ink-soft transition-colors hover:text-ink"
        aria-label="Что значит добавить настроение"
        :aria-expanded="hintOpen"
        @click="hintOpen = !hintOpen"
        @mouseenter="hintOpen = true"
      >
        ?
      </button>

      <Transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="opacity-0 translate-y-1"
        leave-active-class="transition duration-150 ease-in"
        leave-to-class="opacity-0 translate-y-1"
      >
        <p
          v-if="hintOpen"
          role="tooltip"
          class="absolute left-1/2 top-full z-20 mt-3 w-[min(22rem,80vw)] -translate-x-1/2
                 rounded-2xl border border-chip bg-white/90 px-4 py-3 text-left text-[13px]
                 leading-relaxed text-ink-muted shadow-cta backdrop-blur"
        >
          {{ MOOD_HINT }}
        </p>
      </Transition>
    </div>

    <button
      type="button"
      :disabled="!canContinue"
      class="mt-8 h-[58px] w-[min(280px,100%)] rounded-pill text-[15px] font-medium transition-all duration-200"
      :class="
        canContinue
          ? 'bg-brand text-white shadow-cta hover:-translate-y-0.5 hover:bg-brand-strong'
          : 'cursor-not-allowed bg-white/50 text-ink-soft/70 shadow-chip'
      "
      @click="onContinue"
    >
      Продолжить
    </button>

    <p v-if="notice" class="mt-5 text-sm text-ink-soft" role="status">
      Практика скоро откроется. Мы запускаемся — следите за обновлениями.
    </p>
  </section>
</template>
