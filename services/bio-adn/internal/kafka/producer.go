package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewProducer(brokers string, logger *zap.Logger) *Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers),
		Balancer:               &kafka.LeastBytes{},
		BatchTimeout:           10 * time.Millisecond,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
	}
	return &Producer{writer: w, logger: logger}
}

func (p *Producer) PublishProfileCreated(ctx context.Context, e any) error {
	return p.publish(ctx, "snisid.bio.profile.created", e)
}

func (p *Producer) PublishProfileUploaded(ctx context.Context, e any) error {
	return p.publish(ctx, "snisid.bio.profile.uploaded", e)
}

func (p *Producer) PublishHitDetected(ctx context.Context, e *DNAHitDetected) error {
	e.EventEnvelope = newEnvelope("DNAHitDetected")
	return p.publish(ctx, "snisid.bio.hits", e)
}

func (p *Producer) PublishWantedEvent(ctx context.Context, e *WantedPersonCreated) error {
	e.EventEnvelope = newEnvelope("WantedPersonCreated")
	return p.publish(ctx, "snisid.bio.wanted.events", e)
}

func (p *Producer) PublishMissingEvent(ctx context.Context, e *MissingEvent) error {
	e.EventEnvelope = newEnvelope("MissingPersonReported")
	return p.publish(ctx, "snisid.bio.missing.events", e)
}

func (p *Producer) PublishVehicleStolen(ctx context.Context, e *StolenVehicleEvent) error {
	e.EventEnvelope = newEnvelope("VehicleStolen")
	return p.publish(ctx, "snisid.bio.vehicle.stolen", e)
}

func (p *Producer) PublishVehicleRecovered(ctx context.Context, e *VehicleRecoveredEvent) error {
	e.EventEnvelope = newEnvelope("VehicleRecovered")
	return p.publish(ctx, "snisid.bio.vehicle.recovered", e)
}

func (p *Producer) PublishDocumentStolen(ctx context.Context, e *StolenDocumentEvent) error {
	e.EventEnvelope = newEnvelope("DocumentStolen")
	return p.publish(ctx, "snisid.bio.document.stolen", e)
}

func (p *Producer) PublishVesselStolen(ctx context.Context, e *StolenVesselEvent) error {
	e.EventEnvelope = newEnvelope("VesselStolen")
	return p.publish(ctx, "snisid.bio.vessel.stolen", e)
}

func (p *Producer) PublishArmStolen(ctx context.Context, e *StolenFirearmEvent) error {
	e.EventEnvelope = newEnvelope("FirearmStolen")
	return p.publish(ctx, "snisid.bio.arm.stolen", e)
}

func (p *Producer) PublishArmCrimeSceneHit(ctx context.Context, e *ArmCrimeSceneHitEvent) error {
	e.EventEnvelope = newEnvelope("ArmCrimeSceneHit")
	return p.publish(ctx, "snisid.bio.arm.hit", e)
}

func (p *Producer) PublishLAPIQuery(ctx context.Context, e *LAPIPlateQuery) error {
	e.EventEnvelope = newEnvelope("LAPIPlateQuery")
	return p.publish(ctx, "snisid.bio.lapi.query", e)
}

func (p *Producer) PublishLAPIResponse(ctx context.Context, e *LAPIPlateResponse) error {
	e.EventEnvelope = newEnvelope("LAPIPlateResponse")
	return p.publish(ctx, "snisid.bio.lapi.query", e)
}

func (p *Producer) PublishExpungeEvent(ctx context.Context, e *ExpungeEvent) error {
	e.EventEnvelope = newEnvelope("ProfileExpunged")
	return p.publish(ctx, "snisid.bio.expunge.events", e)
}

func (p *Producer) PublishArticleStolen(ctx context.Context, e *StolenArticleEvent) error {
	e.EventEnvelope = newEnvelope("ArticleStolen")
	return p.publish(ctx, "snisid.bio.article.stolen", e)
}

func (p *Producer) PublishSecurityStolen(ctx context.Context, e *StolenSecurityEvent) error {
	e.EventEnvelope = newEnvelope("SecurityStolen")
	return p.publish(ctx, "snisid.bio.security.stolen", e)
}

func (p *Producer) PublishONIDocumentRevoked(ctx context.Context, e *ONIDocumentRevokedEvent) error {
	e.EventEnvelope = newEnvelope("ONIDocumentRevoked")
	return p.publish(ctx, "snisid.oni.document.revoked", e)
}

func (p *Producer) PublishForeignFugitive(ctx context.Context, e *ForeignFugitiveCreated) error {
	e.EventEnvelope = newEnvelope("ForeignFugitiveCreated")
	return p.publish(ctx, "snisid.bio.fugitive.events", e)
}

func (p *Producer) PublishUnidentifiedPerson(ctx context.Context, e *UnidentifiedPersonCreated) error {
	e.EventEnvelope = newEnvelope("UnidentifiedPersonCreated")
	return p.publish(ctx, "snisid.bio.unidentified.events", e)
}

func (p *Producer) PublishTerrorismWatch(ctx context.Context, e *TerrorismWatchCreated) error {
	e.EventEnvelope = newEnvelope("TerrorismWatchCreated")
	return p.publish(ctx, "snisid.bio.terrorism.events", e)
}

func (p *Producer) PublishProtectionOrder(ctx context.Context, e *ProtectionOrderCreated) error {
	e.EventEnvelope = newEnvelope("ProtectionOrderCreated")
	return p.publish(ctx, "snisid.bio.protection.events", e)
}

func (p *Producer) PublishSupervisedRelease(ctx context.Context, e *SupervisedReleaseCreated) error {
	e.EventEnvelope = newEnvelope("SupervisedReleaseCreated")
	return p.publish(ctx, "snisid.bio.supervised.events", e)
}

func (p *Producer) PublishAuditEvent(ctx context.Context, e map[string]any) error {
	e["event_type"] = "AuditEvent"
	e["timestamp"] = time.Now().UnixMilli()
	return p.publish(ctx, "snisid.bio.audit.events", e)
}

func (p *Producer) PublishDuplicateSpecimen(ctx context.Context, e *DuplicateSpecimenDetected) error {
	e.EventEnvelope = newEnvelope("DuplicateSpecimenDetected")
	return p.publish(ctx, "snisid.bio.lab.duplicate", e)
}

func (p *Producer) PublishEquipmentRegistered(ctx context.Context, e *EquipmentRegistered) error {
	e.EventEnvelope = newEnvelope("EquipmentRegistered")
	return p.publish(ctx, "snisid.bio.lab.equipment", e)
}

func (p *Producer) PublishTrainingRecorded(ctx context.Context, e *TrainingRecorded) error {
	e.EventEnvelope = newEnvelope("TrainingRecorded")
	return p.publish(ctx, "snisid.bio.lab.training", e)
}

func (p *Producer) PublishUploadCompleted(ctx context.Context, e *LDISUploadCompleted) error {
	e.EventEnvelope = newEnvelope("LDISUploadCompleted")
	return p.publish(ctx, "snisid.bio.lab.upload", e)
}

func (p *Producer) PublishCrossDeptHit(ctx context.Context, e *CrossDeptHitDetected) error {
	e.EventEnvelope = newEnvelope("CrossDeptHitDetected")
	return p.publish(ctx, "snisid.bio.ndis.crossdept.hit", e)
}

func (p *Producer) PublishInterpolSubmission(ctx context.Context, e *InterpolSubmissionRequested) error {
	e.EventEnvelope = newEnvelope("InterpolSubmissionRequested")
	return p.publish(ctx, "snisid.bio.ndis.interpol", e)
}

func (p *Producer) PublishNDISReport(ctx context.Context, e *NDISReportGenerated) error {
	e.EventEnvelope = newEnvelope("NDISReportGenerated")
	return p.publish(ctx, "snisid.bio.ndis.reports", e)
}

func (p *Producer) PublishViolenceRecord(ctx context.Context, e *ViolenceRecordCreated) error {
	e.EventEnvelope = newEnvelope("ViolenceRecordCreated")
	return p.publish(ctx, "snisid.bio.violence.events", e)
}

func (p *Producer) PublishIdentityTheft(ctx context.Context, e *IdentityTheftRecorded) error {
	e.EventEnvelope = newEnvelope("IdentityTheftRecorded")
	return p.publish(ctx, "snisid.bio.identitytheft.events", e)
}

func (p *Producer) PublishIdentityLinked(ctx context.Context, e *BioIdentityLinked) error {
	e.EventEnvelope = newEnvelope("BioIdentityLinked")
	return p.publish(ctx, "snisid.bio.identity.linked", e)
}

func (p *Producer) PublishSLABreach(ctx context.Context, plate string, responseMs int) {
	_ = p.publish(ctx, "snisid.bio.lapi.query", map[string]any{
		"event_type":  "LAPI_SLA_BREACH",
		"plate":       plate,
		"response_ms": responseMs,
		"timestamp":   time.Now().UnixMilli(),
	})
}

func (p *Producer) publish(ctx context.Context, topic string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	msg := kafka.Message{
		Topic: topic,
		Value: data,
		Time:  time.Now(),
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error("kafka publish failed", zap.String("topic", topic), zap.Error(err))
		return fmt.Errorf("write message: %w", err)
	}
	p.logger.Info("kafka published", zap.String("topic", topic))
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
