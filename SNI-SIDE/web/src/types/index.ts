export interface AlertEvent {
  event_id: string
  alert_type: string
  severity: 'CRITICAL' | 'HIGH' | 'MEDIUM' | 'LOW' | 'INFO'
  title: string
  description: string
  entities_involved: EntityRef[]
  confidence_score: number
  timestamp: number
  source: string
}

export interface EntityRef {
  entity_type: string
  entity_id: string
  entity_name?: string
}

export interface IntelligenceReport {
  report_id: string
  report_type: string
  generated_at: string
  target_entity: { type: string; id: string; label: string }
  executive_summary: string
  key_findings: string[]
  risk_assessment: { overall_risk: string; risk_factors: string[]; confidence: number }
  graph_context: { node_count: number; relationship_count: number }
  recommendations: string[]
  confidence_score: number
}

export interface SearchResult {
  entity_id: string
  entity_type: string
  label: string
  score: number
  source: string
  details: Record<string, unknown>
  match_type?: string
}

export interface UnifiedSearchResponse {
  query: string
  detected_type: string
  results: SearchResult[]
  graph_context?: { nodes: unknown[]; relationships: unknown[] }
  execution_time_ms: number
}

export interface DashboardStats {
  total_alerts_24h: number
  critical_alerts: number
  persons_wanted: number
  vehicles_wanted: number
  events_per_second: number
  active_cases: number
  graph_nodes: number
  graph_relationships: number
  system_health: 'healthy' | 'degraded' | 'critical'
}
