package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type HandlerFunc func(ctx context.Context, msg map[string]any) error

type Consumer struct {
	reader  *kafka.Reader
	logger  *zap.Logger
	handler HandlerFunc
}

type ConsumerConfig struct {
	GroupID     string
	Topic       string
	Description string
	SLAMaxLag   int
}

var ConsumerGroups = []ConsumerConfig{
	{GroupID: "snisid-bio-sdis-sync", Topic: "snisid.bio.profile.created", Description: "LDIS->SDIS sync", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-ndis-matcher", Topic: "snisid.bio.profile.uploaded", Description: "NDIS matching trigger", SLAMaxLag: 500},
	{GroupID: "snisid-bio-alert-dispatcher", Topic: "snisid.bio.hits", Description: "Hit alert dispatch", SLAMaxLag: 100},
	{GroupID: "snisid-lapi-responder", Topic: "snisid.bio.lapi.query", Description: "LAPI real-time response", SLAMaxLag: 50},
	{GroupID: "snisid-bio-wanted-sync", Topic: "snisid.bio.wanted.events", Description: "Wanted index sync", SLAMaxLag: 500},
	{GroupID: "snisid-bio-missing-sync", Topic: "snisid.bio.missing.events", Description: "Missing index sync", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-vehicle-sync", Topic: "snisid.bio.vehicle.stolen", Description: "Stolen vehicle index", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-vehicle-recovered", Topic: "snisid.bio.vehicle.recovered", Description: "Vehicle recovery + DIDComm notify", SLAMaxLag: 500},
	{GroupID: "snisid-bio-document-sync", Topic: "snisid.bio.document.stolen", Description: "Stolen document index", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-vessel-sync", Topic: "snisid.bio.vessel.stolen", Description: "Stolen vessel index", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-arm-sync", Topic: "snisid.bio.arm.stolen", Description: "Stolen firearm index", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-arm-hit", Topic: "snisid.bio.arm.hit", Description: "Firearm crime scene hit cross-ref", SLAMaxLag: 100},
	{GroupID: "snisid-bio-expunge", Topic: "snisid.bio.expunge.events", Description: "Expunge handler", SLAMaxLag: 5000},
	{GroupID: "snisid-bio-audit-writer", Topic: "snisid.bio.audit.events", Description: "Audit log writer", SLAMaxLag: 5000},
	{GroupID: "snisid-bio-article-sync", Topic: "snisid.bio.article.stolen", Description: "Stolen article index", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-security-sync", Topic: "snisid.bio.security.stolen", Description: "Stolen security index", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-oni-doc-revoked", Topic: "snisid.oni.document.revoked", Description: "ONI document revocation -> BIE-DOC", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-fugitive-sync", Topic: "snisid.bio.fugitive.events", Description: "PER-FUG foreign fugitive sync", SLAMaxLag: 500},
	{GroupID: "snisid-bio-unidentified-sync", Topic: "snisid.bio.unidentified.events", Description: "PER-NID unidentified persons sync", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-terrorism-sync", Topic: "snisid.bio.terrorism.events", Description: "PER-TER terrorism watch sync", SLAMaxLag: 500},
	{GroupID: "snisid-bio-protection-sync", Topic: "snisid.bio.protection.events", Description: "PER-OPR protection order sync", SLAMaxLag: 500},
	{GroupID: "snisid-bio-supervised-sync", Topic: "snisid.bio.supervised.events", Description: "PER-LIB supervised release sync", SLAMaxLag: 500},
	{GroupID: "snisid-bio-lab-duplicate", Topic: "snisid.bio.lab.duplicate", Description: "Duplicate specimen alert", SLAMaxLag: 500},
	{GroupID: "snisid-bio-lab-equipment", Topic: "snisid.bio.lab.equipment", Description: "Equipment registry sync", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-lab-training", Topic: "snisid.bio.lab.training", Description: "Staff training records sync", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-lab-upload", Topic: "snisid.bio.lab.upload", Description: "LDIS upload completion", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-ndis-crossdept", Topic: "snisid.bio.ndis.crossdept.hit", Description: "NDIS cross-dept hit notification", SLAMaxLag: 100},
	{GroupID: "snisid-bio-ndis-interpol", Topic: "snisid.bio.ndis.interpol", Description: "INTERPOL submission", SLAMaxLag: 1000},
	{GroupID: "snisid-bio-ndis-reports", Topic: "snisid.bio.ndis.reports", Description: "NDIS weekly reports", SLAMaxLag: 5000},
	{GroupID: "snisid-bio-violence-sync", Topic: "snisid.bio.violence.events", Description: "PER-VIO known violence sync", SLAMaxLag: 500},
	{GroupID: "snisid-bio-identitytheft-sync", Topic: "snisid.bio.identitytheft.events", Description: "PER-IDV identity theft sync", SLAMaxLag: 500},
	{GroupID: "snisid-bio-identity-linker", Topic: "snisid.bio.identity.linked", Description: "BioIdentityLink dissociation sync", SLAMaxLag: 500},
}

func NewConsumer(cfg ConsumerConfig, brokers string, logger *zap.Logger, handler HandlerFunc) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokers},
		GroupID:     cfg.GroupID,
		Topic:       cfg.Topic,
		MinBytes:    1,
		MaxBytes:    10e6,
		MaxWait:     1 * time.Second,
		StartOffset: kafka.LastOffset,
	})
	return &Consumer{reader: r, logger: logger, handler: handler}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.FetchMessage(ctx)
				if err != nil {
					c.logger.Error("kafka fetch error", zap.Error(err))
					time.Sleep(1 * time.Second)
					continue
				}

				var event map[string]any
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					c.logger.Warn("kafka unmarshal error", zap.Error(err))
					c.reader.CommitMessages(ctx, msg)
					continue
				}

				if err := c.handler(ctx, event); err != nil {
					log.Printf("handler error: %v", err)
				}

				c.reader.CommitMessages(ctx, msg)
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func StartAllConsumers(ctx context.Context, brokers string, logger *zap.Logger, handlers map[string]HandlerFunc) []*Consumer {
	var consumers []*Consumer
	for _, cfg := range ConsumerGroups {
		h, ok := handlers[cfg.Topic]
		if !ok {
			h = func(ctx context.Context, msg map[string]any) error {
				log.Printf("[%s] received: %v", cfg.GroupID, msg["event_type"])
				return nil
			}
		}
		c := NewConsumer(cfg, brokers, logger, h)
		c.Start(ctx)
		consumers = append(consumers, c)
		logger.Info("consumer started", zap.String("group", cfg.GroupID), zap.String("topic", cfg.Topic))
	}
	return consumers
}
