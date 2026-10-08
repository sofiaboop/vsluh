---
applyTo: "vsluh-frontend/**,vsluh-miniapp/**"
---

# Nuxt 4 / Vue 3

Структура — Feature-Sliced Design. Одинакова для лендинга и мини-аппа.

```
app/
├── pages/        # роуты, тонкие — делегируют в features
├── widgets/      # составные блоки экрана
├── features/     # фичи, UI в features/<name>/ui/
├── entities/     # доменные сущности, UI в entities/<name>/ui/
├── shared/       # api/, lib/, constants/, types/
├── composables/  # useX
├── stores/       # Pinia
└── assets/css/
```

## Правила

- Только `<script setup lang="ts">`. Пропсы и эмиты типизированы дженериками:
  `defineProps<{ title: string }>()`, `defineEmits<{ select: [id: string] }>()`.
- Импорты — абсолютные через `~/`. Кросс-слойные относительные импорты запрещены.
- Цвета, тени и радиусы — только из `tailwind.config.js` (токены `brand`, `surface`, `ink`).
  Произвольные hex в шаблонах запрещены.
- Состояние — Pinia, типизированное (`state`, `getters`, `actions`).
- Любой текст интерфейса — русский, в `shared/constants/`, не в разметке вперемешку с логикой.
- Доступность: интерактивные элементы — `<button>`/`<a>`, с `aria-pressed` / `aria-label`.
