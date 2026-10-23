package domain

// MaxAudioBytes ограничивает размер присылаемой дорожки. Лямбда получает тело
// запроса в base64 и упирается в лимит API Gateway (10 МБ), поэтому режем раньше.
const MaxAudioBytes = 6 << 20

// FillerWord — слово-паразит и сколько раз оно прозвучало.
type FillerWord struct {
	Word  string `json:"word"  dynamodbav:"word"`
	Count int    `json:"count" dynamodbav:"count"`
}

// Metrics — измеримая часть разбора. Считается по расшифровке и таймингам слов,
// поэтому не зависит от доступности текстовой модели.
type Metrics struct {
	DurationSec  float64      `json:"duration_sec"  dynamodbav:"duration_sec"`
	SpeechSec    float64      `json:"speech_sec"    dynamodbav:"speech_sec"`
	WordCount    int          `json:"word_count"    dynamodbav:"word_count"`
	UniqueWords  int          `json:"unique_words"  dynamodbav:"unique_words"`
	WordsPerMin  int          `json:"words_per_min" dynamodbav:"words_per_min"`
	FillerCount  int          `json:"filler_count"  dynamodbav:"filler_count"`
	FillerWords  []FillerWord `json:"filler_words"  dynamodbav:"filler_words"`
	RepeatCount  int          `json:"repeat_count"  dynamodbav:"repeat_count"`
	LongPauses   int          `json:"long_pauses"   dynamodbav:"long_pauses"`
	LongestPause float64      `json:"longest_pause" dynamodbav:"longest_pause"`
}

// Check — одна проверка: что сказал пользователь, как это измерено и где лежит запись.
type Check struct {
	ID         string  `json:"id"          dynamodbav:"check_id"`
	UserEmail  string  `json:"-"           dynamodbav:"user_email"`
	CreatedAt  int64   `json:"created_at"  dynamodbav:"created_at"`
	SortKey    string  `json:"-"           dynamodbav:"sort_key"`
	TopicID    string  `json:"topic_id"    dynamodbav:"topic_id"`
	TopicLabel string  `json:"topic_label" dynamodbav:"topic_label"`
	Mode       Mode    `json:"mode"        dynamodbav:"mode"`
	Transcript string  `json:"transcript"  dynamodbav:"transcript"`
	Metrics    Metrics `json:"metrics"     dynamodbav:"metrics"`
	// Advice — связный совет текстовой модели. Пустой, когда модель недоступна:
	// в этом случае интерфейс показывает подсказки из Tips.
	Advice string `json:"advice" dynamodbav:"advice"`
	// Tips — короткие рекомендации по правилам, считаются всегда.
	Tips []string `json:"tips" dynamodbav:"tips"`
	// AudioKey — объект в S3. Наружу отдаётся не ключ, а подписанная ссылка.
	AudioKey string `json:"-" dynamodbav:"audio_key"`
	AudioURL string `json:"audio_url,omitempty" dynamodbav:"-"`
}

// Stats — агрегаты личного кабинета.
type Stats struct {
	TotalChecks    int          `json:"total_checks"`
	TotalSpeechSec float64      `json:"total_speech_sec"`
	TotalWords     int          `json:"total_words"`
	TotalFillers   int          `json:"total_fillers"`
	AvgWordsPerMin int          `json:"avg_words_per_min"`
	// AvgFillersPerCheck — сколько слов-паразитов приходится на одну проверку,
	// с одним знаком после запятой: на коротких ответах целое всегда давало ноль.
	AvgFillersPerCheck float64 `json:"avg_fillers_per_check"`
	StreakDays         int     `json:"streak_days"`
	TopFillers     []FillerWord `json:"top_fillers"`
	Daily          []DailyPoint `json:"daily"`
}

// DailyPoint — одна точка графика динамики. Дата в UTC, формат YYYY-MM-DD.
type DailyPoint struct {
	Date        string  `json:"date"`
	Checks      int     `json:"checks"`
	SpeechSec   float64 `json:"speech_sec"`
	WordsPerMin int     `json:"words_per_min"`
	FillerCount int     `json:"filler_count"`
}
