package quality

import (
	"context"
	"strings"
	"unicode/utf8"

	"scrapper/internal/domain"
)

type SimpleEvaluator struct {
	minLength int
}

func NewSimpleEvaluator(minLength int) *SimpleEvaluator {
	if minLength < 0 {
		minLength = 0
	}

	return &SimpleEvaluator{minLength: minLength}
}

func (e *SimpleEvaluator) Evaluate(ctx context.Context, data domain.ContentData) (domain.QualityEvaluation, error) {
	_ = ctx

	text := strings.TrimSpace(data.Text)
	runeCount := utf8.RuneCountInString(text)

	evaluation := domain.QualityEvaluation{
		Score: float64(runeCount),
	}

	if e.minLength <= 0 {
		evaluation.Accepted = runeCount > 0
		return evaluation, nil
	}

	evaluation.Accepted = runeCount >= e.minLength
	return evaluation, nil
}
