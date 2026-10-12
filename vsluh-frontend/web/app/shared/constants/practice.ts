export type SpeakingMode = 'spontaneous' | 'deep'

export interface ModeOption {
  id: SpeakingMode
  /** Секунды на подготовку */
  prepareSec: number
  /** Секунды на речь */
  speakSec: number
}

export const MODES: ModeOption[] = [
  { id: 'spontaneous', prepareSec: 30, speakSec: 60 },
  { id: 'deep', prepareSec: 600, speakSec: 60 },
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

export const TOPIC_IDS: TopicId[] = [
  'personal',
  'people',
  'work',
  'life',
  'opinions',
  'stories',
  'fantasy',
  'unexpected',
]

/** Сколько вариантов настроения лежит в каждой локали (`mood.options`). */
export const MOOD_COUNT = 5

/** Геометрия рулетки. Значения из макета: 18 строк, старт на 3-й, финиш на 15-й. */
export const WHEEL = {
  rows: 18,
  startIndex: 2,
  selectedIndex: 14,
  rowHeight: 120,
  rowHeightMobile: 108,
  spinMs: 4350,
}

export const SUPPORT = {
  payUrl: 'https://pay.cloudtips.ru/p/c76e8470',
  authorUrl: 'https://www.instagram.com/gurl_sonya',
  authorHandle: '@gurl_sonya',
}
