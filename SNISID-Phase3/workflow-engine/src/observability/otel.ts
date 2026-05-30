/**
 * SNISID — OpenTelemetry bootstrap.
 * Exports traces + metrics + logs to the central OTel collector.
 */
import { NodeSDK } from '@opentelemetry/sdk-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import { config } from '../config.js';

const sdk = new NodeSDK({
  serviceName: 'snisid-workflow-engine',
  instrumentations: [getNodeAutoInstrumentations()]
});

export function startObservability() {
  sdk.start();
  process.on('SIGTERM', () => sdk.shutdown().catch(console.error));
}
