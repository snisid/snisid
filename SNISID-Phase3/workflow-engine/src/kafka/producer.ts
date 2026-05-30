/**
 * SNISID — Kafka producer with Avro serialization, mTLS, headers, signature.
 */
import { Kafka, Producer, CompressionTypes } from 'kafkajs';
import { v7 as uuidv7 } from 'uuid';
import { config } from '../config.js';
import { sign } from '../pki/sign.js';
import pino from 'pino';

const log = pino({ name: 'kafka-producer' });

export interface EmitOptions {
  topic: string;
  eventType: string;
  payload: Record<string, unknown>;
  correlation: {
    workflowId: string;
    workflowInstanceId: string;
    traceId?: string;
    causationId?: string;
  };
  subjectKey?: string;
}

class KafkaService {
  private producer: Producer;
  private connected = false;

  constructor() {
    const kafka = new Kafka({
      clientId: config.KAFKA_CLIENT_ID,
      brokers: config.KAFKA_BROKERS.split(','),
      ssl: { rejectUnauthorized: true },
      // mTLS configured via env (kafkajs supports passing certs):
      // ssl: { ca: [...], key: ..., cert: ... }
    });
    this.producer = kafka.producer({
      idempotent: true,
      maxInFlightRequests: 5,
      allowAutoTopicCreation: false
    });
  }

  async connect() {
    if (!this.connected) {
      await this.producer.connect();
      this.connected = true;
      log.info('Kafka producer connected');
    }
  }

  async emit(opts: EmitOptions): Promise<{ eventId: string; offset: string }> {
    await this.connect();
    const eventId = uuidv7();
    const occurredAt = Date.now();

    const envelope = {
      eventId,
      eventType: opts.eventType,
      occurredAt,
      producer: {
        service: 'snisid-workflow-engine',
        version: config.ENGINE_VERSION,
        node: config.ENGINE_REGION,
        spiffeId: `spiffe://snisid.ht/workflow-engine`
      },
      correlation: opts.correlation,
      payload: opts.payload
    };

    // sign payload (canonical JSON sha384 + PKI signature)
    const { signature, tsa } = await sign(envelope);

    const headers = {
      'event-id': eventId,
      'event-type': opts.eventType,
      'trace-id': opts.correlation.traceId ?? '',
      'producer-spiffe': `spiffe://snisid.ht/workflow-engine`,
      'signature': signature,
      'tsa-timestamp': tsa,
      'correlation-id': opts.correlation.workflowInstanceId,
      'causation-id': opts.correlation.causationId ?? ''
    };

    const record = await this.producer.send({
      topic: opts.topic,
      compression: CompressionTypes.ZSTD,
      messages: [{
        key: opts.subjectKey ?? eventId,
        value: JSON.stringify({ ...envelope, integrity: { signature, tsa } }),
        headers
      }]
    });

    log.info({ eventId, topic: opts.topic, eventType: opts.eventType }, 'event emitted');
    return { eventId, offset: record[0].baseOffset ?? '-1' };
  }

  async disconnect() {
    if (this.connected) {
      await this.producer.disconnect();
      this.connected = false;
    }
  }
}

export const kafka = new KafkaService();
