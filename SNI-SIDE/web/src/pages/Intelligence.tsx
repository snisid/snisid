import { useState } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import { GitBranch, FileText, Search, AlertTriangle, Copy, ChevronDown, ChevronUp } from 'lucide-react'
import { api } from '../services/api'

const REPORT_TYPES = [
  { value: 'ENTITY_PROFILE', label: 'Profil Entité' },
  { value: 'LINK_ANALYSIS', label: 'Analyse de Liens' },
  { value: 'NETWORK_MAP', label: 'Carte du Réseau' },
  { value: 'FINANCIAL_FLOW', label: 'Flux Financiers' },
  { value: 'TEMPORAL_ANALYSIS', label: 'Analyse Temporelle' },
  { value: 'CROSS_ENTITY', label: 'Relation Croisée' },
]

export default function Intelligence() {
  const [entityId, setEntityId] = useState('')
  const [entityType, setEntityType] = useState('Citizen')
  const [reportType, setReportType] = useState('ENTITY_PROFILE')
  const [entityId2, setEntityId2] = useState('')
  const [expandedSections, setExpandedSections] = useState<Record<string, boolean>>({})

  const mutation = useMutation({
    mutationFn: () =>
      api.generateIntelligence({
        entity_id: entityId,
        entity_type: entityType,
        report_type: reportType,
        entity_id2: entityId2 || undefined,
      }),
  })

  const report = mutation.data

  const toggleSection = (key: string) =>
    setExpandedSections((s) => ({ ...s, [key]: !s[key] }))

  const Section = ({ title, key: sectionKey, children }: any) => (
    <div className="border border-gray-800 rounded-lg overflow-hidden">
      <button
        onClick={() => toggleSection(sectionKey)}
        className="w-full flex items-center justify-between px-4 py-3 bg-gray-800/50 hover:bg-gray-800 transition-colors"
      >
        <span className="font-medium text-sm">{title}</span>
        {expandedSections[sectionKey] ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
      </button>
      {expandedSections[sectionKey] && <div className="p-4">{children}</div>}
    </div>
  )

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Génération d'Intelligence</h1>

      <div className="card space-y-4">
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1 block">ID Entité</label>
            <input
              value={entityId}
              onChange={(e) => setEntityId(e.target.value)}
              placeholder="NIU, plaque, IP..."
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm
                         focus:outline-none focus:border-primary-500"
            />
          </div>
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1 block">Type</label>
            <select
              value={entityType}
              onChange={(e) => setEntityType(e.target.value)}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm
                         focus:outline-none focus:border-primary-500"
            >
              {['Citizen', 'Vehicle', 'Phone', 'Gang', 'BankAccount', 'IP', 'Domain', 'Wallet'].map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1 block">Type de Rapport</label>
            <select
              value={reportType}
              onChange={(e) => setReportType(e.target.value)}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm
                         focus:outline-none focus:border-primary-500"
            >
              {REPORT_TYPES.map((t) => (
                <option key={t.value} value={t.value}>{t.label}</option>
              ))}
            </select>
          </div>
        </div>

        {reportType === 'CROSS_ENTITY' && (
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1 block">ID Seconde Entité</label>
            <input
              value={entityId2}
              onChange={(e) => setEntityId2(e.target.value)}
              placeholder="NIU, plaque..."
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm"
            />
          </div>
        )}

        <button
          onClick={() => mutation.mutate()}
          disabled={!entityId || mutation.isPending}
          className="bg-primary-600 hover:bg-primary-500 disabled:bg-gray-700 disabled:text-gray-500
                     text-white px-6 py-2 rounded-lg text-sm font-medium transition-colors"
        >
          {mutation.isPending ? 'Génération en cours...' : 'Générer le Rapport'}
        </button>
      </div>

      {mutation.isError && (
        <div className="card text-red-400 flex items-center gap-2">
          <AlertTriangle className="w-4 h-4" />
          Erreur: {(mutation.error as Error).message}
        </div>
      )}

      {report && (
        <div className="space-y-3">
          <div className="card bg-primary-600/5 border-primary-500/20">
            <p className="text-lg font-semibold mb-1">{report.executive_summary}</p>
            <div className="flex items-center gap-4 mt-2 text-xs text-gray-400">
              <span>Confiance: {(report.confidence_score * 100).toFixed(0)}%</span>
              <span className={`badge-${report.risk_assessment?.overall_risk?.toLowerCase() || 'medium'}`}>
                Risque: {report.risk_assessment?.overall_risk || 'N/A'}
              </span>
              <span>Nœuds: {report.graph_context?.node_count}</span>
              <span>Relations: {report.graph_context?.relationship_count}</span>
            </div>
          </div>

          <Section title="Constatations Clés" key_="findings">
            <ul className="space-y-2">
              {report.key_findings?.map((f: string, i: number) => (
                <li key={i} className="flex items-start gap-2 text-sm">
                  <span className="text-primary-400 mt-0.5">•</span>
                  {f}
                </li>
              ))}
            </ul>
          </Section>

          <Section title="Analyse des Connexions" key_="connections">
            <p className="text-sm text-gray-300">{report.connections_analysis || report.network_metrics?.density !== undefined ? (
              <pre className="text-xs text-gray-400">{JSON.stringify(report.network_metrics, null, 2)}</pre>
            ) : 'Non disponible'}</p>
          </Section>

          {report.suspicious_patterns?.length > 0 && (
            <Section title="Patterns Suspects" key_="patterns">
              <div className="space-y-3">
                {report.suspicious_patterns.map((p: any, i: number) => (
                  <div key={i} className="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg">
                    <span className={`badge-${p.severity?.toLowerCase() || 'medium'}`}>{p.severity}</span>
                    <div>
                      <p className="text-sm">{p.pattern}</p>
                      <p className="text-xs text-gray-500">{p.entities?.join(', ')}</p>
                    </div>
                  </div>
                ))}
              </div>
            </Section>
          )}

          <Section title="Recommandations" key_="recommendations">
            <ol className="space-y-2">
              {report.recommendations?.map((r: string, i: number) => (
                <li key={i} className="flex items-start gap-2 text-sm">
                  <span className="text-primary-400 font-mono min-w-[20px]">{i + 1}.</span>
                  {r}
                </li>
              ))}
            </ol>
          </Section>

          <Section title="Indicateurs" key_="indicators">
            <pre className="text-xs text-gray-400 font-mono whitespace-pre-wrap">
              {JSON.stringify(report.indicators || report.aml_indicators || report.timeline, null, 2)}
            </pre>
          </Section>
        </div>
      )}
    </div>
  )
}
