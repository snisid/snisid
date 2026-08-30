import { FileText, FolderOpen } from 'lucide-react'

const MOCK_CASES = [
  { id: 'CASE-2026-001', title: 'Enquête trafic international de stupéfiants', type: 'NARCOTICS', status: 'ACTIVE', priority: 'CRITICAL', opened: '2026-01-15', subjects: 12, agency: 'BNCD' },
  { id: 'CASE-2026-002', title: 'Réseau de blanchiment d\'argent — Secteur financier', type: 'FINANCIAL', status: 'ACTIVE', priority: 'HIGH', opened: '2026-02-20', subjects: 8, agency: 'UCREF' },
  { id: 'CASE-2026-003', title: 'Disparition mineure — Alerte AMBER', type: 'MISSING', status: 'ACTIVE', priority: 'CRITICAL', opened: '2026-03-05', subjects: 1, agency: 'PNH' },
  { id: 'CASE-2026-004', title: 'Cyberattaque infrastructures gouvernementales', type: 'CYBER', status: 'ACTIVE', priority: 'HIGH', opened: '2026-04-10', subjects: 5, agency: 'CERT-HT' },
  { id: 'CASE-2025-098', title: 'Trafic d\'armes — Gang 400 Mawozo', type: 'FIREARMS', status: 'CLOSED', priority: 'HIGH', opened: '2025-06-01', closed: '2026-02-15', subjects: 23, agency: 'PNH' },
  { id: 'CASE-2025-045', title: 'Enquête préliminaire — Corruption publique', type: 'CORRUPTION', status: 'SUSPENDED', priority: 'MEDIUM', opened: '2025-03-10', subjects: 3, agency: 'ULCC' },
]

export default function Cases() {
  const priorityClass = (p: string) => `badge-${p.toLowerCase()}`
  const statusClass = (s: string) =>
    s === 'ACTIVE' ? 'text-green-400 bg-green-500/10' :
    s === 'CLOSED' ? 'text-gray-400 bg-gray-800' :
    'text-yellow-400 bg-yellow-500/10'

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Cas Criminels</h1>
      <div className="grid grid-cols-3 gap-4">
        <div className="card">
          <p className="stat-label">Cas Actifs</p>
          <p className="stat-value text-primary-400">{MOCK_CASES.filter((c) => c.status === 'ACTIVE').length}</p>
        </div>
        <div className="card">
          <p className="stat-label">Priorité Critique</p>
          <p className="stat-value text-red-400">{MOCK_CASES.filter((c) => c.priority === 'CRITICAL').length}</p>
        </div>
        <div className="card">
          <p className="stat-label">Sujets Impliqués</p>
          <p className="stat-value text-gray-200">{MOCK_CASES.reduce((a, c) => a + (c.subjects || 0), 0)}</p>
        </div>
      </div>
      <div className="space-y-3">
        {MOCK_CASES.map((c) => (
          <div key={c.id} className="card hover:border-gray-600 transition-colors cursor-pointer">
            <div className="flex items-start justify-between">
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-gray-800 mt-1">
                  <FolderOpen className="w-5 h-5 text-primary-400" />
                </div>
                <div>
                  <p className="font-medium">{c.title}</p>
                  <div className="flex items-center gap-2 mt-1">
                    <span className="text-xs text-gray-500">{c.id}</span>
                    <span className={priorityClass(c.priority)}>{c.priority}</span>
                    <span className={`text-xs px-2 py-0.5 rounded ${statusClass(c.status)}`}>{c.status}</span>
                    <span className="text-xs text-gray-500">{c.type}</span>
                  </div>
                  <div className="flex items-center gap-3 mt-2 text-xs text-gray-500">
                    <span>Ouvert: {c.opened}</span>
                    {c.closed && <span>Fermé: {c.closed}</span>}
                    <span>Sujets: {c.subjects}</span>
                    <span>{c.agency}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
