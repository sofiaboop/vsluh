<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    active?: boolean
    pressedLabel?: string
  }>(),
  { active: false, pressedLabel: undefined },
)

const emit = defineEmits<{ toggle: [] }>()
</script>

<template>
  <button
    type="button"
    class="chip"
    :class="{ 'chip--active': props.active }"
    :aria-pressed="props.active"
    @click="emit('toggle')"
  >
    {{ props.label }}
  </button>
</template>
