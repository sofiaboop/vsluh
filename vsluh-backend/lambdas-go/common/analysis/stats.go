package analysis

import (
	"sort"
	"time"

	"common/domain"
)

// topFillersLimit — сколько слов-паразитов показываем в кабинете.
const topFillersLimit = 5

// dailyWindow — глубина графика динамики в днях.
const dailyWindow = 30

// Aggregate собирает статистику кабинета по истории проверок.
// now передаётся явно: серия дней зависит от сегодняшней даты, а её надо тестировать.
func Aggregate(checks []domain.Check, now time.Time) domain.Stats {
	stats := domain.Stats{
		TotalChecks: len(checks),
		TopFillers:  []domain.FillerWord{},
		Daily:       []domain.DailyPoint{},
	}
	if len(checks) == 0 {
		return stats
	}

	fillers := map[string]int{}
	byDay := map[string]*domain.DailyPoint{}
	days := map[string]bool{}
	paceSum, paceCount := 0, 0

	for _, c := range checks {
		stats.TotalSpeechSec += c.Metrics.SpeechSec
		stats.TotalWords += c.Metrics.WordCount
		stats.TotalFillers += c.Metrics.FillerCount

		if c.Metrics.WordsPerMin > 0 {
			paceSum += c.Metrics.WordsPerMin
			paceCount++
		}
		for _, f := range c.Metrics.FillerWords {
			fillers[f.Word] += f.Count
		}

		day := time.Unix(c.CreatedAt, 0).UTC().Format("2006-01-02")
		days[day] = true

		point, ok := byDay[day]
		if !ok {
			point = &domain.DailyPoint{Date: day}
			byDay[day] = point
		}
		point.Checks++
		point.SpeechSec += c.Metrics.SpeechSec
		point.FillerCount += c.Metrics.FillerCount
		// В точке дня держим суммарный темп, после обхода поделим на число проверок.
		point.WordsPerMin += c.Metrics.WordsPerMin
	}

	stats.TotalSpeechSec = round1(stats.TotalSpeechSec)
	stats.AvgFillersPerCheck = round1(float64(stats.TotalFillers) / float64(len(checks)))
	if paceCount > 0 {
		stats.AvgWordsPerMin = paceSum / paceCount
	}

	for word, count := range fillers {
		stats.TopFillers = append(stats.TopFillers, domain.FillerWord{Word: word, Count: count})
	}
	sort.Slice(stats.TopFillers, func(i, j int) bool {
		if stats.TopFillers[i].Count != stats.TopFillers[j].Count {
			return stats.TopFillers[i].Count > stats.TopFillers[j].Count
		}
		return stats.TopFillers[i].Word < stats.TopFillers[j].Word
	})
	if len(stats.TopFillers) > topFillersLimit {
		stats.TopFillers = stats.TopFillers[:topFillersLimit]
	}

	for _, point := range byDay {
		if point.Checks > 0 {
			point.WordsPerMin /= point.Checks
			point.SpeechSec = round1(point.SpeechSec)
		}
		stats.Daily = append(stats.Daily, *point)
	}
	sort.Slice(stats.Daily, func(i, j int) bool { return stats.Daily[i].Date < stats.Daily[j].Date })
	if len(stats.Daily) > dailyWindow {
		stats.Daily = stats.Daily[len(stats.Daily)-dailyWindow:]
	}

	stats.StreakDays = streak(days, now)
	return stats
}

// streak считает, сколько дней подряд человек практиковался.
// Серия не рвётся, если сегодня ещё не занимался: отсчёт начинаем со вчера.
func streak(days map[string]bool, now time.Time) int {
	day := now.UTC().Truncate(24 * time.Hour)
	if !days[day.Format("2006-01-02")] {
		day = day.AddDate(0, 0, -1)
	}

	count := 0
	for days[day.Format("2006-01-02")] {
		count++
		day = day.AddDate(0, 0, -1)
	}
	return count
}
