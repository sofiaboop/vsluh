<template>
  <main class="app-shell">
    <div class="ambient ambient--one" />
    <div class="ambient ambient--two" />

    <TopBar />

    <slot />

    <footer class="footer">
      <span>Практика речи</span>
    </footer>
  </main>
</template>
