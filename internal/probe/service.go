package probe

import (
	"context"

	"super-proxy-pool/internal/events"
)

type Service struct {
	events *events.Broker
}

func NewService(broker *events.Broker) *Service {
	return &Service{events: broker}
}

func (s *Service) Start(context.Context) {}

func (s *Service) EnqueueLatency(sourceType string, sourceNodeID int64) error {
	s.events.Publish("probe.queued", map[string]any{
		"source_type":    sourceType,
		"source_node_id": sourceNodeID,
		"test_type":      "latency",
	})
	return nil
}

func (s *Service) EnqueueSpeed(sourceType string, sourceNodeID int64) error {
	s.events.Publish("probe.queued", map[string]any{
		"source_type":    sourceType,
		"source_node_id": sourceNodeID,
		"test_type":      "speed",
	})
	return nil
}
