import { apiFetch } from '~/shared/api/apiInstance'

export interface FillerWord {
  word: string
  count: number
}

export interface CheckMetrics {
  durationSec: number
  speechSec: number
  wordCount: number
  uniqueWords: number
  wordsPerMin: number
  fillerCount: number
  fillerWords: FillerWord[]
  repeatCount: number
  longPauses: number
  longestPause: number
}

export interface Check {
  id: string
  createdAt: number
  topicId: string
  topicLabel: string
  mode: string
  transcript: string
  metrics: CheckMetrics
  /** Связный совет текстовой модели. Пустой, когда модель недоступна. */
  advice: string
  /** Коды коротких рекомендаций по правилам: ключи переводов analysis.tips.*. */
  tips: string[]
  audioUrl: string
}

export interface DailyPoint {
  date: string
  checks: number
  speechSec: number
  wordsPerMin: number
  fillerCount: number
}

export interface Stats {
  totalChecks: number
  totalSpeechSec: number
  totalWords: number
  avgWordsPerMin: number
  streakDays: number
  topFillers: FillerWord[]
  daily: DailyPoint[]
}

interface MetricsPayload {
  duration_sec: number
  speech_sec: number
  word_count: number
  unique_words: number
  words_per_min: number
  filler_count: number
  filler_words: FillerWord[] | null
  repeat_count: number
  long_pauses: number
  longest_pause: number
}

interface CheckPayload {
  id: string
  created_at: number
  topic_id: string
  topic_label: string
  mode: string
  transcript: string
  metrics: MetricsPayload
  advice: string
  tips: string[] | null
  audio_url?: string
}

interface StatsPayload {
  total_checks: number
  total_speech_sec: number
  total_words: number
  avg_words_per_min: number
  streak_days: number
  top_fillers: FillerWord[] | null
  daily: DailyPoint[] | null
}

function toCheck(payload: CheckPayload): Check {
  return {
    id: payload.id,
    createdAt: payload.created_at,
    topicId: payload.topic_id,
    topicLabel: payload.topic_label,
    mode: payload.mode,
    transcript: payload.transcript,
    advice: payload.advice ?? '',
    tips: payload.tips ?? [],
    audioUrl: payload.audio_url ?? '',
    metrics: {
      durationSec: payload.metrics.duration_sec,
      speechSec: payload.metrics.speech_sec,
      wordCount: payload.metrics.word_count,
      uniqueWords: payload.metrics.unique_words,
      wordsPerMin: payload.metrics.words_per_min,
      fillerCount: payload.metrics.filler_count,
      fillerWords: payload.metrics.filler_words ?? [],
      repeatCount: payload.metrics.repeat_count,
      longPauses: payload.metrics.long_pauses,
      longestPause: payload.metrics.longest_pause,
    },
  }
}

/** Отправляет запись на разбор. Кредит списывает бэкенд, остаток приходит в ответе. */
export async function createCheck(input: {
  audioBase64: string
  format: string
  topicId: string
  topicLabel: string
  mode: string
  speakSec: number
  language: string
}): Promise<{ check: Check; credits: number }> {
  const payload = await apiFetch<{ check: CheckPayload; credits: number }>('/v1/checks', {
    method: 'POST',
    body: JSON.stringify({
      audio_base64: input.audioBase64,
      format: input.format,
      topic_id: input.topicId,
      topic_label: input.topicLabel,
      mode: input.mode,
      speak_sec: input.speakSec,
      language: input.language,
    }),
  })
  return { check: toCheck(payload.check), credits: payload.credits }
}

export async function fetchChecks(cursor = ''): Promise<{ checks: Check[]; nextCursor: string }> {
  const query = cursor ? `?cursor=${encodeURIComponent(cursor)}` : ''
  const payload = await apiFetch<{ checks: CheckPayload[]; next_cursor: string }>(
    `/v1/checks${query}`,
  )
  return { checks: payload.checks.map(toCheck), nextCursor: payload.next_cursor ?? '' }
}

export async function fetchCheck(id: string): Promise<Check> {
  const payload = await apiFetch<{ check: CheckPayload }>(`/v1/checks/${id}`)
  return toCheck(payload.check)
}

export async function fetchStats(): Promise<Stats> {
  const payload = await apiFetch<{ stats: StatsPayload }>('/v1/stats')
  return {
    totalChecks: payload.stats.total_checks,
    totalSpeechSec: payload.stats.total_speech_sec,
    totalWords: payload.stats.total_words,
    avgWordsPerMin: payload.stats.avg_words_per_min,
    streakDays: payload.stats.streak_days,
    topFillers: payload.stats.top_fillers ?? [],
    daily: payload.stats.daily ?? [],
  }
}
