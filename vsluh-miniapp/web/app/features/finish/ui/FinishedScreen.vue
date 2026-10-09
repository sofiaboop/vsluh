<script setup lang="ts">
const { again, repeatTopic, backToTopics } = usePractice()
</script>

<template>
  <section class="content finished-content">
    <div class="done-mark"><span /></div>

    <h1>Готово.</h1>
    <p class="lead">Главное — ты сказал это вслух.</p>

    <div class="finished-actions">
      <button type="button" class="action-button action-button--primary" @click="again">
        Ещё одну
      </button>
      <button type="button" class="action-button action-button--secondary" @click="repeatTopic">
        Повторить тему
      </button>
      <button type="button" class="action-button action-button--tertiary" @click="backToTopics">
        Изменить темы
      </button>
    </div>
  </section>
</template>
