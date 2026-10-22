package analysis

import "testing"

func TestBuildCountsWordsAndPace(t *testing.T) {
	words := []Word{
		{Text: "я", Begin: 0, End: 0.4},
		{Text: "думаю", Begin: 0.5, End: 1.0},
		{Text: "что", Begin: 2.8, End: 3.0},
	}

	m := Build("Я думаю, что...", words, 60)

	if m.WordCount != 3 {
		t.Fatalf("word count = %d, ожидалось 3", m.WordCount)
	}
	if m.SpeechSec != 3 {
		t.Fatalf("speech sec = %v, ожидалось 3", m.SpeechSec)
	}
	if m.WordsPerMin != 60 {
		t.Fatalf("wpm = %d, ожидалось 60", m.WordsPerMin)
	}
	if m.LongPauses != 1 {
		t.Fatalf("long pauses = %d, ожидалась 1", m.LongPauses)
	}
	if m.LongestPause != 1.8 {
		t.Fatalf("longest pause = %v, ожидалось 1.8", m.LongestPause)
	}
}

func TestBuildFallsBackToPlannedDuration(t *testing.T) {
	m := Build("одно слово", nil, 60)

	if m.SpeechSec != 60 {
		t.Fatalf("speech sec = %v, ожидалось 60", m.SpeechSec)
	}
	if m.WordsPerMin != 2 {
		t.Fatalf("wpm = %d, ожидалось 2", m.WordsPerMin)
	}
}

func TestCountFillersPrefersLongestMatch(t *testing.T) {
	m := Build("Ну как бы я вот не знаю, как бы объяснить", nil, 60)

	counts := map[string]int{}
	for _, f := range m.FillerWords {
		counts[f.Word] = f.Count
	}

	if counts["как бы"] != 2 {
		t.Fatalf("«как бы» = %d, ожидалось 2", counts["как бы"])
	}
	if counts["ну"] != 1 || counts["вот"] != 1 {
		t.Fatalf("неожиданный разбор паразитов: %#v", counts)
	}
	if m.FillerCount != 4 {
		t.Fatalf("всего паразитов = %d, ожидалось 4", m.FillerCount)
	}
}

func TestFillersDoNotMatchInsideWords(t *testing.T) {
	m := Build("Нужно вотум доверия", nil, 60)

	if m.FillerCount != 0 {
		t.Fatalf("паразиты внутри слов: %#v", m.FillerWords)
	}
}

func TestRepeatsIgnoreStopWords(t *testing.T) {
	m := Build("и и я я думаю думаю", nil, 60)

	if m.RepeatCount != 1 {
		t.Fatalf("повторы = %d, ожидался 1", m.RepeatCount)
	}
}

func TestTipsCoverPaceAndFillers(t *testing.T) {
	clean := Build("Сегодня я расскажу о том, как работает привычка говорить вслух каждый день подряд", nil, 40)
	tips := Tips(clean)

	if len(tips) == 0 || tips[0] != "fillers_none" {
		t.Fatalf("без паразитов ожидался fillers_none, получили %v", tips)
	}

	noisy := Build("ну ну вот вот типа короче блин значит", nil, 60)
	if got := Tips(noisy); got[0] != "fillers_many" {
		t.Fatalf("при восьми паразитах ожидался fillers_many, получили %v", got)
	}
}
