export type SpeakingMode = 'spontaneous' | 'deep'

export interface ModeOption {
  id: SpeakingMode
  label: string
  duration: string
  /** Секунды на подготовку */
  prepareSec: number
  /** Секунды на речь */
  speakSec: number
}

export const MODES: ModeOption[] = [
  {
    id: 'spontaneous',
    label: 'спонтанно',
    duration: 'У вас будет 30 секунд подготовки, чтобы говорить 1 минуту',
    prepareSec: 30,
    speakSec: 60,
  },
  {
    id: 'deep',
    label: 'глубоко',
    duration: 'У вас будет 10 минут подготовки, чтобы говорить 2 минуты',
    prepareSec: 600,
    speakSec: 120,
  },
]

export type TopicId =
  | 'personal'
  | 'people'
  | 'work'
  | 'life'
  | 'opinions'
  | 'stories'
  | 'fantasy'
  | 'unexpected'

export interface Topic {
  id: TopicId
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

export const MOODS = [
  'с любопытством',
  'максимально серьёзно',
  'как будто рассказываешь секрет',
  'с лёгкой иронией',
  'очень уверенно',
]

/** Геометрия рулетки. Значения из макета: 18 строк, старт на 3-й, финиш на 15-й. */
export const WHEEL = {
  rows: 18,
  startIndex: 2,
  selectedIndex: 14,
  rowHeight: 120,
  rowHeightMobile: 108,
  spinMs: 4350,
}

export const MOOD_TOOLTIP =
  'Сайт задаст случайную манеру подачи — от «максимально серьёзно» до «как будто рассказываешь секрет».'

export const SUPPORT = {
  title: 'Помоги «Вслух» расти',
  text: 'Если этот тренажёр оказался тебе полезен, ты можешь поддержать его развитие любой суммой.',
  payUrl: 'https://pay.cloudtips.ru/p/c76e8470',
  authorUrl: 'https://www.instagram.com/gurl_sonya',
  authorHandle: '@gurl_sonya',
}
