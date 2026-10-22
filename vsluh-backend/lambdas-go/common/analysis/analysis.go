// Package analysis считает метрики речи по расшифровке и таймингам слов.
// Пакет намеренно без внешних зависимостей: это чистая функция от входа,
// её поведение проверяется тестами и не зависит от доступности моделей.
package analysis

import (
	"sort"
	"strings"
	"unicode"

	"common/domain"
)

// Word — слово расшифровки с таймингами в секундах.
type Word struct {
	Text  string
	Begin float64
	End   float64
}

// longPauseSec — пауза, которую слушатель уже замечает как заминку.
const longPauseSec = 1.5

// comfortPaceMin и comfortPaceMax — коридор разговорного темпа в словах в минуту.
const (
	comfortPaceMin = 100
	comfortPaceMax = 160
)

// fillers — слова-паразиты русской разговорной речи. Многословные идут первыми:
// «как бы» должно посчитаться целиком, а не распасться на отдельные слова.
var fillers = []string{
	"как бы", "то есть", "в общем", "в принципе", "на самом деле", "так сказать",
	"типа того", "вот это вот", "скажем так", "короче говоря",
	"ну", "вот", "типа", "короче", "блин", "значит", "эээ", "ээ", "эм", "мм", "ммм",
	"как-то", "просто", "реально", "собственно", "допустим", "слушай",
}

// stopWords не считаем повторами: служебные слова обязаны повторяться.
var stopWords = map[string]bool{
	"и": true, "в": true, "не": true, "на": true, "я": true, "что": true, "с": true,
	"а": true, "по": true, "но": true, "к": true, "у": true, "за": true, "из": true,
	"от": true, "до": true, "для": true, "это": true, "как": true, "бы": true,
	"же": true, "ли": true, "то": true, "о": true, "об": true, "при": true, "так": true,
}

// Build считает метрики одной проверки. speakSec — отведённое время ответа,
// нужно для случая, когда тайминги слов не пришли.
func Build(transcript string, words []Word, speakSec float64) domain.Metrics {
	tokens := tokenize(transcript)

	m := domain.Metrics{
		DurationSec: speakSec,
		WordCount:   len(tokens),
	}

	if len(words) > 0 {
		m.SpeechSec = round1(words[len(words)-1].End - words[0].Begin)
		m.LongPauses, m.LongestPause = pauses(words)
	}
	if m.SpeechSec <= 0 {
		m.SpeechSec = round1(speakSec)
	}

	if m.SpeechSec > 0 {
		m.WordsPerMin = int(float64(m.WordCount) / m.SpeechSec * 60)
	}

	m.FillerWords = countFillers(transcript)
	for _, f := range m.FillerWords {
		m.FillerCount += f.Count
	}

	m.UniqueWords = uniqueCount(tokens)
	m.RepeatCount = repeats(tokens)

	return m
}

// Tips превращает метрики в короткие советы. Возвращает максимум три штуки:
// длинный список обесценивает каждый пункт.
func Tips(m domain.Metrics) []string {
	var tips []string

	switch {
	case m.FillerCount >= 5:
		tips = append(tips, "fillers_many")
	case m.FillerCount > 0:
		tips = append(tips, "fillers_few")
	default:
		tips = append(tips, "fillers_none")
	}

	switch {
	case m.WordsPerMin == 0:
		// Темп посчитать не удалось — совета о темпе не даём.
	case m.WordsPerMin > comfortPaceMax:
		tips = append(tips, "pace_fast")
	case m.WordsPerMin < comfortPaceMin:
		tips = append(tips, "pace_slow")
	default:
		tips = append(tips, "pace_good")
	}

	switch {
	case m.LongPauses >= 3:
		tips = append(tips, "pauses_many")
	case m.RepeatCount >= 3:
		tips = append(tips, "repeats_many")
	case m.WordCount < 40:
		tips = append(tips, "short_answer")
	default:
		tips = append(tips, "keep_going")
	}

	return tips
}

// tokenize режет текст на слова в нижнем регистре, выбрасывая пунктуацию.
func tokenize(text string) []string {
	fields := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-'
	})

	out := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.Trim(f, "-")
		if f != "" {
			out = append(out, f)
		}
	}
	return out
}

// countFillers ищет паразиты по нормализованному тексту с границами слов,
// чтобы «ну» не нашлось внутри «нужно».
func countFillers(transcript string) []domain.FillerWord {
	text := " " + strings.Join(tokenize(transcript), " ") + " "

	var found []domain.FillerWord
	for _, filler := range fillers {
		needle := " " + filler + " "
		count := 0
		for {
			idx := strings.Index(text, needle)
			if idx < 0 {
				break
			}
			count++
			// Затираем найденное, чтобы «как бы» не посчиталось ещё раз как «бы».
			text = text[:idx+1] + strings.Repeat("\x00", len(filler)) + text[idx+1+len(filler):]
		}
		if count > 0 {
			found = append(found, domain.FillerWord{Word: filler, Count: count})
		}
	}

	sort.Slice(found, func(i, j int) bool {
		if found[i].Count != found[j].Count {
			return found[i].Count > found[j].Count
		}
		return found[i].Word < found[j].Word
	})
	return found
}

// pauses считает заметные паузы между словами.
func pauses(words []Word) (int, float64) {
	count := 0
	longest := 0.0
	for i := 1; i < len(words); i++ {
		gap := words[i].Begin - words[i-1].End
		if gap > longest {
			longest = gap
		}
		if gap >= longPauseSec {
			count++
		}
	}
	return count, round1(longest)
}

func uniqueCount(tokens []string) int {
	seen := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		seen[t] = struct{}{}
	}
	return len(seen)
}

// repeats считает подряд идущие повторы значимых слов: «я я думаю».
func repeats(tokens []string) int {
	count := 0
	for i := 1; i < len(tokens); i++ {
		if tokens[i] == tokens[i-1] && !stopWords[tokens[i]] && len([]rune(tokens[i])) > 2 {
			count++
		}
	}
	return count
}

func round1(v float64) float64 {
	if v < 0 {
		return 0
	}
	return float64(int(v*10+0.5)) / 10
}
