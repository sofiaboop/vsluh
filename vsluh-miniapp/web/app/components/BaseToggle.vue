<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
  label: string
}>()

const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>

<template>
  <button
    type="button"
    role="switch"
    :aria-checked="props.modelValue"
    :aria-label="props.label"
    class="relative h-7 w-[52px] shrink-0 rounded-pill border border-chip transition-colors duration-200"
    :class="props.modelValue ? 'bg-brand' : 'bg-[#D6DCE6]'"
    @click="emit('update:modelValue', !props.modelValue)"
  >
    <span
      class="absolute top-1/2 h-[22px] w-[22px] -translate-y-1/2 rounded-full bg-white
             shadow-chip transition-all duration-200"
      :class="props.modelValue ? 'left-[27px]' : 'left-[2px]'"
    />
  </button>
</template>
