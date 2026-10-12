import type { SpeakingMode, TopicId } from '~/shared/constants/practice'
import { DEEP_TOPICS_AR } from '~/shared/data/deepTopics.ar'
import { DEEP_TOPICS_DE } from '~/shared/data/deepTopics.de'
import { DEEP_TOPICS_EN } from '~/shared/data/deepTopics.en'
import { DEEP_TOPICS_ES } from '~/shared/data/deepTopics.es'
import { DEEP_TOPICS_FR } from '~/shared/data/deepTopics.fr'
import { DEEP_TOPICS_IT } from '~/shared/data/deepTopics.it'
import { DEEP_TOPICS_PL } from '~/shared/data/deepTopics.pl'
import { DEEP_TOPICS_PT } from '~/shared/data/deepTopics.pt'
import { DEEP_TOPICS_RU } from '~/shared/data/deepTopics.ru'
import { DEEP_TOPICS_TR } from '~/shared/data/deepTopics.tr'
import { QUESTIONS_AR } from '~/shared/data/questions.ar'
import { QUESTIONS_DE } from '~/shared/data/questions.de'
import { QUESTIONS_EN } from '~/shared/data/questions.en'
import { QUESTIONS_ES } from '~/shared/data/questions.es'
import { QUESTIONS_FR } from '~/shared/data/questions.fr'
import { QUESTIONS_IT } from '~/shared/data/questions.it'
import { QUESTIONS_PL } from '~/shared/data/questions.pl'
import { QUESTIONS_PT } from '~/shared/data/questions.pt'
import { QUESTIONS_RU } from '~/shared/data/questions.ru'
import { QUESTIONS_TR } from '~/shared/data/questions.tr'

export type TopicBank = Record<TopicId, string[]>

const SPONTANEOUS: Record<string, TopicBank> = {
  ru: QUESTIONS_RU,
  en: QUESTIONS_EN,
  es: QUESTIONS_ES,
  pt: QUESTIONS_PT,
  fr: QUESTIONS_FR,
  de: QUESTIONS_DE,
  it: QUESTIONS_IT,
  pl: QUESTIONS_PL,
  tr: QUESTIONS_TR,
  ar: QUESTIONS_AR,
}

const DEEP: Record<string, TopicBank> = {
  ru: DEEP_TOPICS_RU,
  en: DEEP_TOPICS_EN,
  es: DEEP_TOPICS_ES,
  pt: DEEP_TOPICS_PT,
  fr: DEEP_TOPICS_FR,
  de: DEEP_TOPICS_DE,
  it: DEEP_TOPICS_IT,
  pl: DEEP_TOPICS_PL,
  tr: DEEP_TOPICS_TR,
  ar: DEEP_TOPICS_AR,
}

/** Язык интерфейса может быть любым, а банк тем — только переведённым: остальные откатываются на английский. */
export function topicBank(locale: string, mode: SpeakingMode): TopicBank {
  const banks = mode === 'deep' ? DEEP : SPONTANEOUS
  return banks[locale] ?? banks.en!
}
