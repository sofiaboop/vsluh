import type { CheckMetrics } from "~/shared/api/checks";

export type VerdictLevel = "great" | "good" | "work";

/** На чём сильнее всего потеряли баллы: ключ подсказки в карточке истории. */
export type VerdictFocus =
  "clean" | "slow" | "fast" | "fillers" | "repeats" | "pauses";

export interface Verdict {
  score: number;
  level: VerdictLevel;
  focus: VerdictFocus;
}

/** Комфортный коридор темпа устной речи, слов в минуту. */
const PACE_MIN = 90;
const PACE_MAX = 170;

function penalties(metrics: CheckMetrics) {
  // Паразиты, повторы и паузы считаем на минуту речи: иначе длинный ответ
  // выглядел бы хуже короткого просто потому, что в нём больше слов.
  const minutes = Math.max(metrics.speechSec / 60, 0.1);
  const pace = metrics.wordsPerMin;

  return {
    slow: pace > 0 && pace < PACE_MIN ? Math.min(25, (PACE_MIN - pace) / 2) : 0,
    fast: pace > PACE_MAX ? Math.min(25, (pace - PACE_MAX) / 2) : 0,
    fillers: Math.min(30, (metrics.fillerCount / minutes) * 6),
    repeats: Math.min(20, (metrics.repeatCount / minutes) * 4),
    pauses: Math.min(15, (metrics.longPauses / minutes) * 5),
  };
}

/**
 * Оценка ответа от 0 до 100. Считается на фронте из уже измеренных бэкендом
 * метрик, поэтому одинакова и в отчёте, и в истории.
 */
export function verdictFor(metrics: CheckMetrics): Verdict {
  const lost = penalties(metrics);
  const score = Math.max(
    0,
    Math.round(100 - Object.values(lost).reduce((a, b) => a + b, 0)),
  );

  let focus: VerdictFocus = "clean";
  let worst = 3;
  for (const [key, value] of Object.entries(lost)) {
    if (value > worst) {
      worst = value;
      focus = key as VerdictFocus;
    }
  }

  const level: VerdictLevel =
    score >= 80 ? "great" : score >= 55 ? "good" : "work";
  return { score, level, focus };
}
