export type SpeakingMode = 'spontaneous' | 'deep'

export interface ModeOption {
  id: SpeakingMode
  label: string
  duration: string
}

export const MODES: ModeOption[] = [
  {
    id: 'spontaneous',
    label: 'спонтанно',
    duration: 'У вас будет 30 секунд подготовки, чтобы говорить 1 минуту',
  },
  {
    id: 'deep',
    label: 'глубоко',
    duration: 'У вас будет 10 минут подготовки, чтобы говорить 2 минуты',
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

export const MOOD_TOOLTIP =
  'Сайт задаст случайную манеру подачи — от «максимально серьёзно» до «как будто рассказываешь секрет».'

export const SUPPORT = {
  title: 'Помоги «Вслух» расти',
  text: 'Если этот тренажёр оказался тебе полезен, ты можешь поддержать его развитие любой суммой.',
  payUrl: 'https://pay.cloudtips.ru/p/c76e8470',
  authorUrl: 'https://www.instagram.com/gurl_sonya',
  authorHandle: '@gurl_sonya',
}
