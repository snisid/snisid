/**
 * SNISID — Fraud detection activities.
 * Combines DMN rules + ML model + graph anomalies.
 */
import axios from 'axios';
import { kafka } from '../kafka/producer.js';

const FRAUD = process.env.FRAUD_GATEWAY ?? 'https://fraud.snisid.ht';

export async function runRules(input: Record<string, unknown>) {
  const r = await axios.post(`${FRAUD}/rules`, input, { timeout: 5_000 });
  return r.data as { triggered: string[]; ruleScore: number };
}

export async function mlScore(input: Record<string, unknown>) {
  const r = await axios.post(`${FRAUD}/ml/score`, input, { timeout: 5_000 });
  return r.data as { mlScore: number; modelVersion: string };
}

export async function graphAnalyze(input: Record<string, unknown>) {
  const r = await axios.post(`${FRAUD}/graph/analyze`, input, { timeout: 10_000 });
  return r.data as { anomalies: string[]; graphScore: number };
}

export interface FraudOutcome {
  fraudScore: number;
  triggered: string[];
  anomalies: string[];
  recommended: 'ALLOW' | 'REVIEW' | 'SUSPEND' | 'REJECT' | 'INVESTIGATE';
  modelVersion: string;
}

export async function detect(input: {
  workflowId: string;
  workflowInstanceId: string;
  data: Record<string, unknown>;
}): Promise<FraudOutcome> {
  const [rules, ml, graph] = await Promise.all([
    runRules(input.data),
    mlScore(input.data),
    graphAnalyze(input.data)
  ]);
  const fraudScore = Math.min(
    1,
    0.4 * rules.ruleScore + 0.4 * ml.mlScore + 0.2 * graph.graphScore
  );
  const recommended: FraudOutcome['recommended'] =
    fraudScore < 0.3 ? 'ALLOW' :
    fraudScore < 0.6 ? 'REVIEW' :
    fraudScore < 0.8 ? 'SUSPEND' :
    fraudScore < 0.9 ? 'INVESTIGATE' : 'REJECT';

  await kafka.emit({
    topic: 'fraud.detected.v1',
    eventType: 'fraud.detected.v1',
    correlation: {
      workflowId: input.workflowId,
      workflowInstanceId: input.workflowInstanceId
    },
    payload: {
      fraudScore,
      rulesTriggered: rules.triggered,
      mlModelVersion: ml.modelVersion,
      graphAnomalies: graph.anomalies,
      recommendedAction: recommended
    }
  });
  return {
    fraudScore,
    triggered: rules.triggered,
    anomalies: graph.anomalies,
    recommended,
    modelVersion: ml.modelVersion
  };
}
