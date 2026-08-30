import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ArrowLeft, User, AlertTriangle, GitBranch, FileText, ExternalLink } from 'lucide-react'
import { api } from '../services/api'

const entityLabels: Record<string, string> = {
  citizen: 'Citoyen', vehicle: 'Véhicule', phone: 'Téléphone',
  ip: 'Adresse IP', domain: 'Domaine', wallet: 'Portefeuille Crypto',
  gang: 'Gang', bank_account: 'Compte Bancaire', case: 'Cas',
}

export default function EntityDetail() {
  const { type, id } = useParams<{ type: string; id: string }>()
  const navigate = useNavigate()

  const { data: report, isLoading } = useQuery({
    queryKey: ['entity-report', type, id],
    queryFn: () => api.getReport('ENTITY_PROFILE', id!),
    enabled: !!id,
  })

  return (
    <div className="space-y-6">
      <button onClick={() => navigate(-1)} className="flex items-center gap-2 text-sm text-gray-400 hover:text-white">
        <ArrowLeft className="w-4 h-4" /> Retour
      </button>

      <div className="card">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <div className="p-3 rounded-lg bg-primary-600/10">
              <User className="w-6 h-6 text-primary-400" />
            </div>
            <div>
              <h1 className="text-xl font-bold">{report?.target_entity?.label || id}</h1>
              <p className="text-sm text-gray-400">
                {entityLabels[type?.toLowerCase() || ''] || type} · {id}
              </p>
            </div>
          </div>
          {report?.risk_assessment?.overall_risk && (
            <span className={`badge-${report.risk_assessment.overall_risk.toLowerCase()}`}>
              {report.risk_assessment.overall_risk}
            </span>
          )}
        </div>
      </div>

      {isLoading && (
        <div className="flex justify-center py-12">
          <div className="w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin" />
        </div>
      )}

      {report && (
        <>
          <div className="card">
            <h2 className="text-lg font-semibold mb-3">Résumé Exécutif</h2>
            <p className="text-gray-300">{report.executive_summary}</p>
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div className="card">
              <p className="stat-label">Niveau de Risque</p>
              <p className={`stat-value ${
                report.risk_assessment.overall_risk === 'CRITICAL' ? 'text-red-400' :
                report.risk_assessment.overall_risk === 'HIGH' ? 'text-orange-400' :
                'text-yellow-400'
              }`}>{report.risk_assessment.overall_risk}</p>
            </div>
            <div className="card">
              <p className="stat-label">Confiance</p>
              <p className="stat-value text-primary-400">
                {(report.confidence_score * 100).toFixed(0)}%
              </p>
            </div>
            <div className="card">
              <p className="stat-label">Contexte Graphe</p>
              <p className="stat-value text-gray-200">
                {report.graph_context?.node_count || 0}
                <span className="text-sm text-gray-500 ml-1">nœuds</span>
              </p>
            </div>
          </div>

          {report.key_findings && report.key_findings.length > 0 && (
            <div className="card">
              <h2 className="text-lg font-semibold mb-3">Constatations</h2>
              <ul className="space-y-2">
                {report.key_findings.map((f: string, i: number) => (
                  <li key={i} className="flex items-start gap-2 text-sm">
                    <span className="text-primary-400 mt-1">•</span>
                    {f}
                  </li>
                ))}
              </ul>
            </div>
          )}

          {report.recommendations && report.recommendations.length > 0 && (
            <div className="card">
              <h2 className="text-lg font-semibold mb-3">Recommandations</h2>
              <ol className="space-y-2">
                {report.recommendations.map((r: string, i: number) => (
                  <li key={i} className="flex items-start gap-2 text-sm">
                    <span className="text-primary-400 font-mono min-w-[24px]">{i + 1}.</span>
                    {r}
                  </li>
                ))}
              </ol>
            </div>
          )}

          <div className="flex gap-3">
            <button
              onClick={() => navigate(`/search?q=${id}`)}
              className="flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 rounded-lg text-sm transition-colors"
            >
              <GitBranch className="w-4 h-4" /> Analyse de Liens
            </button>
            <button
              onClick={() => navigate(`/intelligence?entity=${id}&type=ENTITY_PROFILE`)}
              className="flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 rounded-lg text-sm transition-colors"
            >
              <FileText className="w-4 h-4" /> Rapport Complet
            </button>
          </div>
        </>
      )}
    </div>
  )
}
