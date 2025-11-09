package processor

import (
	"strings"

	"github.com/google/uuid"

	"scrapper/internal/models"
)

var (
	metadataWarningIgnorePrefixes = []string{
		"strategy ",
		"fallback to root ",
		"persist attempt ",
	}
	metadataWarningIgnoreExact = map[string]struct{}{
		"no strategy succeeded": {},
	}
)

func BuildMetadata(outcome models.RequestOutcome, contentID string) models.MetadataMessage {
	best := outcome.Best.Result

	metadata := models.MetadataMessage{
		RequestID:   outcome.Request.RequestID,
		ContentID:   contentID,
		Status:      best.Status,
		Failure:     best.Failure,
		Error:       best.Error,
		Warnings:    filterMetadataWarnings(best.Warnings),
		URL:         outcome.Best.URL,
		Domain:      best.Domain,
		Strategy:    best.Strategy,
		Metric:      best.Metric,
		UserAgent:   best.UserAgent,
		ProcessedAt: outcome.ProcessedAt,
	}

	if len(metadata.Warnings) == 0 {
		metadata.Warnings = nil
	}

	return metadata
}

func filterMetadataWarnings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	filtered := make([]string, 0, len(values))

valueLoop:
	for _, item := range values {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, skip := metadataWarningIgnoreExact[item]; skip {
			continue
		}
		for _, prefix := range metadataWarningIgnorePrefixes {
			if strings.HasPrefix(item, prefix) {
				continue valueLoop
			}
		}
		filtered = append(filtered, item)
	}

	return filtered
}

func BuildContent(outcome models.RequestOutcome, contentID string) (models.ContentMessage, bool) {
	best := outcome.Best.Result

	if strings.TrimSpace(best.Content) == "" {
		return models.ContentMessage{}, false
	}

	if best.Status != models.StatusOK {
		return models.ContentMessage{}, false
	}

	if best.Failure != models.FailureNone && best.Failure != "" {
		return models.ContentMessage{}, false
	}

	contentMsg := models.ContentMessage{
		RequestID:   outcome.Request.RequestID,
		ContentID:   contentID,
		URL:         outcome.Best.URL,
		Strategy:    best.Strategy,
		Metric:      best.Metric,
		Content:     best.Content,
		Status:      best.Status,
		Failure:     best.Failure,
		Error:       best.Error,
		UserAgent:   best.UserAgent,
		ProcessedAt: outcome.ProcessedAt,
	}

	return contentMsg, true
}

func EnsureContentID(outcome models.RequestOutcome) string {
	contentID := outcome.ContentID
	best := outcome.Best.Result
	if strings.TrimSpace(contentID) == "" && best.Status == models.StatusOK && strings.TrimSpace(best.Content) != "" {
		return uuid.NewString()
	}
	return contentID
}
