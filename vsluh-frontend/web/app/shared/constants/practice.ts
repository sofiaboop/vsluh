export type SpeakingMode = 'spontaneous' | 'deep'

export interface ModeOption {
  id: SpeakingMode
  label: string
  hint: string
}

export const MODES: ModeOption[] = [
  {
    id: 'spontaneous',
    label: 'спонтанно',
    hint: 'У вас будет 30 секунд подготовки, чтобы говорить 1 минуту',
  },
  {
    id: 'deep',
    label: 'глубоко',
    hint: 'У вас будет 2 минуты подготовки, чтобы говорить 3 минуты',
  },
]

export interface Topic {
  id: string
  label: string
}

export const TOPICS: Topic[] = [
  { id: 'personal', label: 'Личное' },
  { id: 'people', label: 'Люди' },
  { id: 'work', label: 'Работа и деньги' },
  { id: 'life', label: 'Жизнь' },
  { id: 'opinions', label: 'Мнения' },
  { id: 'stories', label: 'Истории' },
  { id: 'fantasy', label: 'Фантазия' },
  { id: 'unexpected', label: 'Неожиданное' },
]

export const MOOD_HINT =
  'К теме добавится случайное настроение — например «с иронией» или «с нежностью». ' +
  'Так проще уйти от заученных формулировок.'
