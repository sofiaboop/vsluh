package analysis

import (
	"testing"
	"time"

	"common/domain"
)

func check(day string, speechSec float64, wpm int, fillers ...domain.FillerWord) domain.Check {
	ts, err := time.Parse("2006-01-02", day)
	if err != nil {
		panic(err)
	}
	count := 0
	for _, f := range fillers {
		count += f.Count
	}
	return domain.Check{
		CreatedAt: ts.UTC().Unix(),
		Metrics: domain.Metrics{
			SpeechSec:   speechSec,
			WordCount:   int(speechSec),
			WordsPerMin: wpm,
			FillerWords: fillers,
			FillerCount: count,
		},
	}
}

func TestAggregateEmptyHistory(t *testing.T) {
	stats := Aggregate(nil, time.Now())

	if stats.TotalChecks != 0 || stats.StreakDays != 0 {
		t.Fatalf("пустая история дала %#v", stats)
	}
	if stats.TopFillers == nil || stats.Daily == nil {
		t.Fatal("пустые списки должны быть не nil, иначе в JSON уедет null")
	}
}

func TestAggregateCountsTotalsAndPace(t *testing.T) {
	checks := []domain.Check{
		check("2026-10-20", 50, 120),
		check("2026-10-21", 40, 100),
	}

	stats := Aggregate(checks, time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC))

	if stats.TotalChecks != 2 {
		t.Fatalf("проверок = %d", stats.TotalChecks)
	}
	if stats.TotalSpeechSec != 90 {
		t.Fatalf("секунд = %v", stats.TotalSpeechSec)
	}
	if stats.AvgWordsPerMin != 110 {
		t.Fatalf("средний темп = %d, ожидалось 110", stats.AvgWordsPerMin)
	}
	if len(stats.Daily) != 2 || stats.Daily[0].Date != "2026-10-20" {
		t.Fatalf("динамика собрана неверно: %#v", stats.Daily)
	}
}

func TestAggregateAveragesPaceInsideDay(t *testing.T) {
	checks := []domain.Check{
		check("2026-10-21", 30, 100),
		check("2026-10-21", 30, 140),
	}

	stats := Aggregate(checks, time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC))

	if len(stats.Daily) != 1 {
		t.Fatalf("дней = %d, ожидался 1", len(stats.Daily))
	}
	if stats.Daily[0].WordsPerMin != 120 {
		t.Fatalf("темп дня = %d, ожидалось 120", stats.Daily[0].WordsPerMin)
	}
	if stats.Daily[0].Checks != 2 {
		t.Fatalf("проверок за день = %d", stats.Daily[0].Checks)
	}
}

func TestAggregateTopFillers(t *testing.T) {
	checks := []domain.Check{
		check("2026-10-21", 30, 120, domain.FillerWord{Word: "ну", Count: 2}, domain.FillerWord{Word: "вот", Count: 1}),
		check("2026-10-21", 30, 120, domain.FillerWord{Word: "вот", Count: 5}),
	}

	stats := Aggregate(checks, time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC))

	if stats.TopFillers[0].Word != "вот" || stats.TopFillers[0].Count != 6 {
		t.Fatalf("топ паразитов неверный: %#v", stats.TopFillers)
	}
	if stats.TotalFillers != 8 {
		t.Fatalf("всего паразитов = %d, ожидалось 8", stats.TotalFillers)
	}
	if stats.AvgFillersPerCheck != 4 {
		t.Fatalf("паразитов на проверку = %v, ожидалось 4", stats.AvgFillersPerCheck)
	}
}

func TestAggregateRoundsFillersPerCheck(t *testing.T) {
	// Один паразит на три проверки: целое округление показало бы ноль.
	checks := []domain.Check{
		check("2026-10-21", 30, 120, domain.FillerWord{Word: "ну", Count: 1}),
		check("2026-10-21", 30, 120),
		check("2026-10-21", 30, 120),
	}

	stats := Aggregate(checks, time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC))

	if stats.AvgFillersPerCheck != 0.3 {
		t.Fatalf("паразитов на проверку = %v, ожидалось 0.3", stats.AvgFillersPerCheck)
	}
}

func TestStreakCountsConsecutiveDays(t *testing.T) {
	now := time.Date(2026, 10, 21, 10, 0, 0, 0, time.UTC)
	checks := []domain.Check{
		check("2026-10-21", 30, 120),
		check("2026-10-20", 30, 120),
		check("2026-10-19", 30, 120),
		check("2026-10-17", 30, 120),
	}

	if got := Aggregate(checks, now).StreakDays; got != 3 {
		t.Fatalf("серия = %d, ожидалось 3", got)
	}
}

func TestStreakSurvivesDayWithoutPractice(t *testing.T) {
	// Сегодня ещё не занимался — вчерашняя серия не должна обнуляться.
	now := time.Date(2026, 10, 21, 10, 0, 0, 0, time.UTC)
	checks := []domain.Check{
		check("2026-10-20", 30, 120),
		check("2026-10-19", 30, 120),
	}

	if got := Aggregate(checks, now).StreakDays; got != 2 {
		t.Fatalf("серия = %d, ожидалось 2", got)
	}
}
