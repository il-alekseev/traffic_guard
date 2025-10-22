package processor

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"scrapper/internal/models"
)

func DecodeKafkaMessage(data []byte) (models.ProcessingRequest, error) {
	var req models.ProcessingRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return req, err
	}
	return NormalizeRequest(req)
}

func NormalizeRequest(req models.ProcessingRequest) (models.ProcessingRequest, error) {
	req = normalizeRequest(req)

	if strings.TrimSpace(req.RequestID) == "" {
		return req, errors.New("request_id is required")
	}

	if strings.TrimSpace(req.Dst.Resource) == "" {
		return req, errors.New("destination.resource is required")
	}

	switch req.Dst.Type {
	case models.DestinationIP, models.DestinationDomain, models.DestinationURL:
	default:
		return req, errors.New("unsupported destination type: " + string(req.Dst.Type))
	}

	if req.Received.IsZero() {
		req.Received = time.Now().UTC()
	}

	if req.Raw == nil {
		req.Raw = make(map[string]any)
	}

	return req, nil
}

func normalizeRequest(req models.ProcessingRequest) models.ProcessingRequest {
	req.Dst.Type = normalizeDestinationType(req.Dst.Type)
	if req.Dst.Proto != nil {
		proto := models.Protocol(strings.ToLower(string(*req.Dst.Proto)))
		req.Dst.Proto = &proto
	}
	return req
}

func normalizeDestinationType(value models.DestinationType) models.DestinationType {
	switch strings.ToLower(string(value)) {
	case "ip":
		return models.DestinationIP
	case "domain":
		return models.DestinationDomain
	case "url", "":
		return models.DestinationURL
	default:
		return models.DestinationType(strings.ToLower(string(value)))
	}
}
