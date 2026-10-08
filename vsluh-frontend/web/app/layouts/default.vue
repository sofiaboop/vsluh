<template>
  <div class="relative flex min-h-[100dvh] flex-col overflow-hidden bg-canvas">
    <BackdropGlow />
    <SiteHeader />
    <main class="relative z-10 flex flex-1 flex-col">
      <slot />
    </main>
    <SiteFooter />
  </div>
</template>
