<script setup lang="ts">
const { stage } = usePractice()

const screenClass = computed(() => {
  const map: Record<string, string> = {
    topics: 'screen--topics',
    wheel: 'screen--wheel',
    prepare: 'screen--prepare',
    speak: 'screen--speak',
    finished: 'screen--finished',
  }
  return map[stage.value]
})
</script>

<template>
  <div class="screen" :class="screenClass">
    <TopicSelection v-if="stage === 'topics'" />
    <TopicWheel v-else-if="stage === 'wheel'" />
    <PracticeTimer v-else-if="stage === 'prepare' || stage === 'speak'" />
    <FinishedScreen v-else />
  </div>
</template>
