package router

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/snisid/platform/internal/platform/events"
	"github.com/snisid/platform/internal/platform/logger"
	"go.uber.org/zap"
)

type RouterService struct {
	consumer  events.ConsumerInterface
	engine    *Engine
	producers map[string]events.ProducerInterface
	brokers   []string
	mu        sync.RWMutex
}

func NewRouterService(brokers []string, ingressTopic string, groupID string) (*RouterService, error) {
	engine, err := NewEngine()
	if err != nil {
		return nil, err
	}

	consumer := events.NewConsumer(brokers, groupID, ingressTopic)

	return &RouterService{
		consumer:  consumer,
		engine:    engine,
		producers: make(map[string]events.ProducerInterface),
		brokers:   brokers,
	}, nil
}

func (s *RouterService) Start(ctx context.Context) error {
	return s.consumer.Start(ctx, s.handleEvent)
}

func (s *RouterService) handleEvent(ctx context.Context, payload []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	targets := s.engine.Evaluate(ctx, event)
	if len(targets) == 0 {
		logger.Info(ctx, "Event matched no rules, dropping", zap.Any("event_type", event["type"]))
		return nil
	}

	for _, target := range targets {
		producer := s.getOrCreateProducer(target)
		// We use a generic key or pull from event if available
		key := fmt.Sprintf("%v", event["eventId"])
		if err := producer.Publish(ctx, key, event); err != nil {
			logger.Error(ctx, "Failed to dispatch event", err, zap.String("target", target))
		}
	}

	return nil
}

func (s *RouterService) getOrCreateProducer(topic string) events.ProducerInterface {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p, ok := s.producers[topic]; ok {
		return p
	}

	p := events.NewProducer(s.brokers, topic)
	s.producers[topic] = p
	return p
}

func (s *RouterService) ReloadRules(rules []Rule) error {
	return s.engine.UpdateRules(rules)
}

func (s *RouterService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.producers {
		_ = p.Close()
	}
	return s.consumer.Close()
}
