package util

type TypingResult struct {
	Wpm float64
}

func CalcTypingResult(durationMs int64, words []Word) TypingResult {
	wpm := float64(len(words)) / (float64(durationMs) / 1000 / 60)

	return TypingResult{Wpm: wpm}
}
