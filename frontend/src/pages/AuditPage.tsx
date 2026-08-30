import React, { useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Search, Filter, Download, RefreshCw, Shield,
  User, ChevronRight, ChevronDown, Activity, Eye,
  Lock, Key, Database, Server, Fingerprint, X
} from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

// ─── Types ──────────────────────────────────────────────────────────────────

type EventSeverity  = 'INFO' | 'WARNING' | 'CRITICAL' | 'ALERT';
type EventCategory  = 'AUTH' | 'BIOMETRIC' | 'DATA_ACCESS' | 'ADMIN' | 'SECURITY' | 'SYSTEM';

interface AuditEvent {
  id: string;
  timestamp: string;
  severity: EventSeverity;
  category: EventCategory;
  actor: string;
  actorNni?: string;
  action: string;
  resource: string;
  outcome: 'SUCCESS' | 'FAILURE' | 'BLOCKED';
  ipAddress: string;
  details: string;
  traceId: string;
  aal: string;
}

// ─── Mock data ───────────────────────────────────────────────────────────────

const MOCK_EVENTS: AuditEvent[] = [
  { id: 'EVT-001', timestamp: '2024-05-12T14:32:11Z', severity: 'CRITICAL', category: 'SECURITY', actor: 'SYSTÈME', action: 'TENTATIVE_ACCES_ROOT_CA', resource: 'root-ca.snisid.ht', outcome: 'BLOCKED', ipAddress: '192.168.1.54', details: 'Tentative d\'accès non autorisé à la Root CA depuis une adresse IP inconnue. War Room déclenché automatiquement. PAM Teleport a bloqué la session.', traceId: 'TRC-A9F3-E2B1', aal: 'AAL3' },
  { id: 'EVT-002', timestamp: '2024-05-12T13:15:44Z', severity: 'WARNING', category: 'AUTH', actor: 'jean.pierre@oni.ht', actorNni: 'HTI-2024-001-00002', action: 'MFA_ECHEC_3_FOIS', resource: 'keycloak/token', outcome: 'FAILURE', ipAddress: '10.0.1.22', details: 'Échec de l\'authentification MFA 3 fois de suite. Compte verrouillé temporairement pour 15 minutes. UEBA a enregistré l\'anomalie.', traceId: 'TRC-B2C4-F1D8', aal: 'AAL2' },
  { id: 'EVT-003', timestamp: '2024-05-12T12:55:02Z', severity: 'INFO', category: 'BIOMETRIC', actor: 'operateur.cap@oni.ht', action: 'ENROLEMENT_CITOYEN', resource: 'abis/enroll', outcome: 'SUCCESS', ipAddress: '10.2.0.45', details: 'Enrôlement biométrique réussi pour NNI HTI-2024-002-00099. Score qualité empreintes: 94%. Déduplication 1:N passée (0 doublon détecté).', traceId: 'TRC-C5D7-G3E9', aal: 'AAL2' },
  { id: 'EVT-004', timestamp: '2024-05-12T12:30:17Z', severity: 'ALERT', category: 'DATA_ACCESS', actor: 'audit.interne@dcpj.ht', actorNni: 'HTI-2024-003-00015', action: 'EXPORT_MASSE_TENTATIVE', resource: 'api/citizens/export', outcome: 'BLOCKED', ipAddress: '10.3.0.12', details: 'Tentative d\'export massif de données — 45,000 enregistrements en 2 minutes. DLP a déclenché une alerte. Session PAM gelée. CISO notifié.', traceId: 'TRC-D8E0-H4F2', aal: 'AAL2' },
  { id: 'EVT-005', timestamp: '2024-05-12T11:44:59Z', severity: 'INFO', category: 'AUTH', actor: 'marie.joseph@dgi.ht', actorNni: 'HTI-2024-001-00001', action: 'LOGIN_REUSSI', resource: 'keycloak/sessions', outcome: 'SUCCESS', ipAddress: '10.1.2.100', details: 'Connexion réussie via FIDO2 YubiKey. Niveau AAL2 accordé. Session OIDC créée.', traceId: 'TRC-E1F3-I5G4', aal: 'AAL2' },
  { id: 'EVT-006', timestamp: '2024-05-12T11:20:33Z', severity: 'WARNING', category: 'SECURITY', actor: 'FALCO_eBPF', action: 'SHELL_SPAWNE_CONTAINER', resource: 'k8s/pod/api-gateway-pod-x7k2', outcome: 'BLOCKED', ipAddress: '10.0.5.12', details: 'Shell anormal détecté dans le pod api-gateway (syscall execve). Pod isolé par Cilium NetworkPolicy en < 30 secondes. Forensique container initiée.', traceId: 'TRC-F4G6-J6H5', aal: 'SYSTÈME' },
  { id: 'EVT-007', timestamp: '2024-05-12T10:58:21Z', severity: 'INFO', category: 'ADMIN', actor: 'admin.pki@snisid.ht', action: 'CERT_RENOUVELE', resource: 'pki/issuing-ca/tls-3b', outcome: 'SUCCESS', ipAddress: '10.0.0.5', details: 'Certificat TLS Issuing CA 3B renouvelé via ACME. Durée de validité: 24h. Nouvelles données distribuées aux microservices via Vault agent injection.', traceId: 'TRC-G7H9-K7I6', aal: 'AAL3' },
  { id: 'EVT-008', timestamp: '2024-05-12T10:15:05Z', severity: 'CRITICAL', category: 'BIOMETRIC', actor: 'ABIS_ENGINE', action: 'HARD_MATCH_DETECTE', resource: 'abis/dedup/1n', outcome: 'BLOCKED', ipAddress: 'N/A', details: 'Hard Match ABIS: score 97.3% (seuil ≥95%). NNI HTI-2024-006-00022 en tentative de double enrôlement. Workflow arrêté. DCPJ forensique déclenché.', traceId: 'TRC-H0I2-L8J7', aal: 'SYSTÈME' },
  { id: 'EVT-009', timestamp: '2024-05-12T09:42:17Z', severity: 'INFO', category: 'SYSTEM', actor: 'ARGOCD', action: 'DEPLOY_STAGING', resource: 'k8s/namespace/snisid-staging', outcome: 'SUCCESS', ipAddress: '10.0.0.10', details: 'Déploiement GitOps staging réussi. Image: snisid-api:v2.4.1@sha256:a1b2c3... Cosign signature vérifiée. SLSA Level 3 attestation validée.', traceId: 'TRC-I3J5-M9K8', aal: 'SYSTÈME' },
  { id: 'EVT-010', timestamp: '2024-05-12T09:10:44Z', severity: 'WARNING', category: 'AUTH', actor: 'claude.morisseau@oni.ht', action: 'SESSION_EXPIREE_AAL2', resource: 'keycloak/sessions', outcome: 'SUCCESS', ipAddress: '10.2.1.77', details: 'Session AAL2 expirée après 5 minutes d\'inactivité (conformément aux règles UX). Utilisateur redirigé vers re-authentification MFA.', traceId: 'TRC-J6K8-N0L9', aal: 'AAL2' },
];

// ─── Config ───────────────────────────────────────────────────────────────────

const SEVERITY_CONFIG: Record<EventSeverity, { color: string; bg: string; label: string }> = {
  INFO:     { color: 'text-slate-300',  bg: 'bg-slate-600/20 border-slate-600/30', label: 'INFO'     },
  WARNING:  { color: 'text-yellow-400', bg: 'bg-yellow-500/10 border-yellow-500/30', label: 'WARNING' },
  ALERT:    { color: 'text-orange-400', bg: 'bg-orange-500/10 border-orange-500/30', label: 'ALERTE'  },
  CRITICAL: { color: 'text-red-400',    bg: 'bg-red-500/10 border-red-500/30',       label: 'CRITIQUE'},
};

const CATEGORY_ICON: Record<EventCategory, LucideIcon> = {
  AUTH:        Key,
  BIOMETRIC:   Fingerprint,
  DATA_ACCESS: Database,
  ADMIN:       Server,
  SECURITY:    Shield,
  SYSTEM:      Activity,
};

const OUTCOME_CONFIG = {
  SUCCESS: { color: 'text-emerald-400', label: 'Succès'  },
  FAILURE: { color: 'text-red-400',     label: 'Échec'   },
  BLOCKED: { color: 'text-orange-400',  label: 'Bloqué'  },
};

// ─── Subcomponents ────────────────────────────────────────────────────────────

const SeverityBadge: React.FC<{ s: EventSeverity }> = ({ s }) => {
  const cfg = SEVERITY_CONFIG[s];
  return (
    <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold border ${cfg.bg} ${cfg.color}`}>
      {cfg.label}
    </span>
  );
};

const CategoryIcon: React.FC<{ cat: EventCategory }> = ({ cat }) => {
  const Icon = CATEGORY_ICON[cat];
  const colors: Record<EventCategory, string> = {
    AUTH: 'text-blue-400', BIOMETRIC: 'text-purple-400', DATA_ACCESS: 'text-cyan-400',
    ADMIN: 'text-slate-400', SECURITY: 'text-red-400', SYSTEM: 'text-emerald-400'
  };
  return <Icon size={15} className={colors[cat]} />;
};

const EventDetailModal: React.FC<{ event: AuditEvent; onClose: () => void; lang: string }> = ({ event, onClose, lang }) => (
  <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" role="dialog" aria-modal="true">
    <div className="bg-[#0D1B2A] border border-slate-700/60 rounded-2xl w-full max-w-2xl mx-4 shadow-2xl overflow-hidden">
      <div className="flex items-center justify-between px-6 py-4 border-b border-slate-700/60 bg-[#1E3A5F]/20">
        <div className="flex items-center gap-3">
          <div className={`p-2 rounded-lg border ${SEVERITY_CONFIG[event.severity].bg}`}>
            <CategoryIcon cat={event.category} />
          </div>
          <div>
            <h2 className="font-semibold text-white text-sm">{event.action}</h2>
            <p className="text-xs font-mono text-slate-400">{event.traceId}</p>
          </div>
        </div>
        <button onClick={onClose} className="p-2 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-[#1565C0]">
          <X size={16} />
        </button>
      </div>

      <div className="p-6 space-y-4">
        {/* Metadata grid */}
        <div className="grid grid-cols-2 gap-4">
          {[
            { label: lang==='ht'?'Dat ak Lè':'Date & Heure',        value: new Date(event.timestamp).toLocaleString('fr-HT') },
            { label: 'Sévérité',                                      value: <SeverityBadge s={event.severity} /> },
            { label: lang==='ht'?'Aktè':'Acteur',                    value: event.actor },
            { label: lang==='ht'?'Rezilta':'Résultat',               value: <span className={OUTCOME_CONFIG[event.outcome].color + ' font-medium'}>{OUTCOME_CONFIG[event.outcome].label}</span> },
            { label: lang==='ht'?'Resous':'Ressource',               value: <span className="font-mono text-xs">{event.resource}</span> },
            { label: 'Adresse IP',                                    value: <span className="font-mono text-xs">{event.ipAddress}</span> },
            { label: 'AAL',                                           value: event.aal },
            { label: 'Catégorie',                                     value: event.category },
          ].map(item => (
            <div key={item.label} className="bg-slate-800/40 rounded-lg p-3">
              <div className="text-xs text-slate-500 mb-1">{item.label}</div>
              <div className="text-slate-200 text-sm">{item.value}</div>
            </div>
          ))}
        </div>

        {/* Details */}
        <div className="bg-slate-800/40 rounded-lg p-4">
          <h3 className="text-xs text-slate-500 mb-2 font-medium">{lang==='ht'?'Detay':'Détails'}</h3>
          <p className="text-slate-300 text-sm leading-relaxed">{event.details}</p>
        </div>

        {/* Immutability attestation */}
        <div className="flex items-center gap-2 text-xs text-emerald-400 bg-emerald-400/5 border border-emerald-400/20 rounded-lg p-3">
          <Lock size={13} />
          <span>
            {lang === 'ht'
              ? 'Evènman sa a imyab — Merkle Tree verifye — Siyati ECDSA-P384 valid'
              : 'Événement immuable — Merkle Tree vérifié — Signature ECDSA-P384 valide'}
          </span>
        </div>
      </div>

      <div className="flex justify-end gap-3 px-6 py-4 border-t border-slate-700/60">
        <button onClick={onClose} className="px-4 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 transition-colors focus:outline-none focus:ring-2 focus:ring-slate-600">
          {lang==='ht'?'Fèmen':'Fermer'}
        </button>
        <button className="px-4 py-2 rounded-lg text-sm bg-[#1565C0] text-white hover:bg-blue-600 transition-colors focus:outline-none focus:ring-2 focus:ring-[#1565C0]">
          {lang==='ht'?'Telechaje':'Télécharger'}
        </button>
      </div>
    </div>
  </div>
);

// ─── KPI Stats ────────────────────────────────────────────────────────────────

const StatCard: React.FC<{ label: string; value: string | number; sub?: string; color?: string }> = ({ label, value, sub, color = 'text-slate-100' }) => (
  <div className="bg-[#1E3A5F]/15 border border-slate-700/40 rounded-xl p-4">
    <div className={`text-2xl font-bold ${color}`}>{value}</div>
    <div className="text-slate-400 text-xs mt-1">{label}</div>
    {sub && <div className="text-slate-600 text-xs mt-0.5">{sub}</div>}
  </div>
);

// ─── Main Component ───────────────────────────────────────────────────────────

export const AuditPage: React.FC = () => {
  const { i18n } = useTranslation();
  const lang = i18n.language;

  const [search, setSearch]               = useState('');
  const [severityFilter, setSeverityFilter] = useState<EventSeverity | 'ALL'>('ALL');
  const [categoryFilter, setCategoryFilter] = useState<EventCategory | 'ALL'>('ALL');
  const [outcomeFilter, setOutcomeFilter]   = useState<'ALL'|'SUCCESS'|'FAILURE'|'BLOCKED'>('ALL');
  const [showFilters, setShowFilters]       = useState(false);
  const [selectedEvent, setSelectedEvent]   = useState<AuditEvent | null>(null);

  const filtered = useMemo(() => {
    let list = [...MOCK_EVENTS];
    if (search.trim()) {
      const q = search.toLowerCase();
      list = list.filter(e =>
        e.actor.toLowerCase().includes(q) ||
        e.action.toLowerCase().includes(q) ||
        e.resource.toLowerCase().includes(q) ||
        e.traceId.toLowerCase().includes(q)
      );
    }
    if (severityFilter !== 'ALL') list = list.filter(e => e.severity === severityFilter);
    if (categoryFilter !== 'ALL') list = list.filter(e => e.category === categoryFilter);
    if (outcomeFilter  !== 'ALL') list = list.filter(e => e.outcome === outcomeFilter);
    return list;
  }, [search, severityFilter, categoryFilter, outcomeFilter]);

  const stats = useMemo(() => ({
    total:    MOCK_EVENTS.length,
    critical: MOCK_EVENTS.filter(e => e.severity === 'CRITICAL').length,
    blocked:  MOCK_EVENTS.filter(e => e.outcome === 'BLOCKED').length,
    mttd:     '3m 24s',
  }), []);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white flex items-center gap-3">
            <Shield size={24} className="text-[#1565C0]" />
            {lang === 'ht' ? 'Odit ak Envestigasyon' : 'Audit & Investigation Forensique'}
          </h2>
          <p className="text-slate-400 text-sm mt-0.5">
            {lang === 'ht'
              ? 'Jounal imyab — Merkle Tree ECDSA-P384 — Hot tier 30 jou'
              : 'Journal immuable — Merkle Tree ECDSA-P384 — Hot tier 30 jours'}
          </p>
        </div>
        <div className="flex items-center gap-3">
          <button className="flex items-center gap-2 px-4 py-2 text-sm text-slate-300 bg-slate-800 hover:bg-slate-700 border border-slate-700 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-slate-600">
            <Download size={15} />
            {lang === 'ht' ? 'Ekspòte' : 'Exporter'}
          </button>
          <button className="flex items-center gap-2 px-4 py-2 text-sm text-white bg-[#1565C0] hover:bg-blue-600 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-[#1565C0]">
            <RefreshCw size={15} />
            Live
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <StatCard label={lang==='ht'?'Total evènman':'Événements totaux'} value={stats.total} />
        <StatCard label={lang==='ht'?'Evènman Kritik':'Critiques'} value={stats.critical} color="text-red-400" sub="Dernières 24h" />
        <StatCard label={lang==='ht'?'Bloke':'Bloqués'} value={stats.blocked} color="text-orange-400" sub="Auto-réponse SOAR" />
        <StatCard label="MTTD" value={stats.mttd} color="text-emerald-400" sub={lang==='ht'?'Objektif < 5min':'Objectif < 5 min'} />
      </div>

      {/* Search & Filters */}
      <div className="bg-slate-900/60 border border-slate-700/40 rounded-xl p-4 space-y-3">
        <div className="flex gap-3">
          <div className="relative flex-1">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              id="audit-search"
              type="search"
              placeholder={lang === 'ht' ? 'Chèche pa aktè, aksyon, Trace ID...' : 'Rechercher par acteur, action, Trace ID...'}
              value={search}
              onChange={e => setSearch(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-lg pl-9 pr-4 py-2.5 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-[#1565C0]"
            />
          </div>
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`flex items-center gap-2 px-4 py-2.5 text-sm rounded-lg border transition-colors focus:outline-none focus:ring-2 focus:ring-[#1565C0] ${showFilters ? 'bg-[#1565C0]/20 border-[#1565C0]/40 text-blue-400' : 'bg-slate-800 border-slate-700 text-slate-300 hover:bg-slate-700'}`}
          >
            <Filter size={15} />
            {lang === 'ht' ? 'Filtre' : 'Filtres'}
            <ChevronDown size={14} className={`transition-transform ${showFilters ? 'rotate-180' : ''}`} />
          </button>
        </div>

        {showFilters && (
          <div className="flex flex-wrap gap-3 pt-2 border-t border-slate-700/40">
            <div className="flex items-center gap-2">
              <label className="text-xs text-slate-400">{lang==='ht'?'Severite:':'Sévérité:'}</label>
              <select value={severityFilter} onChange={e => setSeverityFilter(e.target.value as any)}
                className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-[#1565C0]">
                <option value="ALL">{lang==='ht'?'Tout':'Tous'}</option>
                <option value="INFO">INFO</option>
                <option value="WARNING">WARNING</option>
                <option value="ALERT">ALERTE</option>
                <option value="CRITICAL">CRITIQUE</option>
              </select>
            </div>
            <div className="flex items-center gap-2">
              <label className="text-xs text-slate-400">{lang==='ht'?'Kategori:':'Catégorie:'}</label>
              <select value={categoryFilter} onChange={e => setCategoryFilter(e.target.value as any)}
                className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-[#1565C0]">
                <option value="ALL">{lang==='ht'?'Tout':'Tous'}</option>
                <option value="AUTH">Authentification</option>
                <option value="BIOMETRIC">Biométrie</option>
                <option value="DATA_ACCESS">Accès Données</option>
                <option value="ADMIN">Administration</option>
                <option value="SECURITY">Sécurité</option>
                <option value="SYSTEM">Système</option>
              </select>
            </div>
            <div className="flex items-center gap-2">
              <label className="text-xs text-slate-400">{lang==='ht'?'Rezilta:':'Résultat:'}</label>
              <select value={outcomeFilter} onChange={e => setOutcomeFilter(e.target.value as any)}
                className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-[#1565C0]">
                <option value="ALL">{lang==='ht'?'Tout':'Tous'}</option>
                <option value="SUCCESS">{lang==='ht'?'Siksè':'Succès'}</option>
                <option value="FAILURE">{lang==='ht'?'Echèk':'Échec'}</option>
                <option value="BLOCKED">{lang==='ht'?'Bloke':'Bloqué'}</option>
              </select>
            </div>
            {(severityFilter !== 'ALL' || categoryFilter !== 'ALL' || outcomeFilter !== 'ALL' || search) && (
              <button onClick={() => { setSeverityFilter('ALL'); setCategoryFilter('ALL'); setOutcomeFilter('ALL'); setSearch(''); }}
                className="flex items-center gap-1.5 text-xs text-red-400 hover:text-red-300 focus:outline-none">
                <X size={12} /> {lang==='ht'?'Efase filtre':'Réinitialiser'}
              </button>
            )}
          </div>
        )}
      </div>

      {/* Event Timeline */}
      <div className="bg-slate-900/60 border border-slate-700/40 rounded-xl overflow-hidden">
        <div className="px-4 py-3 border-b border-slate-700/40 flex items-center justify-between">
          <h3 className="text-sm font-medium text-slate-300 flex items-center gap-2">
            <Activity size={15} className="text-[#00BCD4]" />
            {lang === 'ht' ? 'Tan reyèl — Journal Forensique' : 'Temps réel — Journal Forensique'}
          </h3>
          <div className="flex items-center gap-2 text-xs text-slate-500">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
            {lang === 'ht' ? 'Live Sync' : 'Sync Live'}
          </div>
        </div>

        <div className="divide-y divide-slate-700/30">
          {filtered.length === 0 ? (
            <div className="py-16 text-center text-slate-500">
              {lang === 'ht' ? 'Pa gen rezilta' : 'Aucun événement trouvé'}
            </div>
          ) : filtered.map((event, idx) => {
            const Icon = CATEGORY_ICON[event.category];
            const sevCfg = SEVERITY_CONFIG[event.severity];
            const outCfg = OUTCOME_CONFIG[event.outcome];

            return (
              <div
                key={event.id}
                className="flex items-start gap-4 px-4 py-4 hover:bg-slate-800/30 transition-colors cursor-pointer group"
                onClick={() => setSelectedEvent(event)}
              >
                {/* Timeline connector */}
                <div className="flex flex-col items-center flex-shrink-0 mt-1">
                  <div className={`w-8 h-8 rounded-full border flex items-center justify-center flex-shrink-0 ${sevCfg.bg}`}>
                    <Icon size={14} className={sevCfg.color} />
                  </div>
                  {idx < filtered.length - 1 && (
                    <div className="w-px flex-1 bg-slate-700/50 mt-2 min-h-4" />
                  )}
                </div>

                {/* Content */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <SeverityBadge s={event.severity} />
                        <span className="text-xs text-slate-500 font-mono">{event.category}</span>
                        <span className={`text-xs font-medium ${outCfg.color}`}>→ {outCfg.label}</span>
                      </div>
                      <p className="text-slate-100 text-sm font-medium mt-1 truncate">{event.action}</p>
                      <p className="text-slate-500 text-xs mt-0.5 truncate">{event.resource}</p>
                    </div>
                    <div className="text-right flex-shrink-0">
                      <div className="text-xs text-slate-400">
                        {new Date(event.timestamp).toLocaleTimeString('fr-HT', { hour: '2-digit', minute: '2-digit', second: '2-digit' })}
                      </div>
                      <div className="text-xs text-slate-600 mt-0.5">
                        {new Date(event.timestamp).toLocaleDateString('fr-HT')}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center gap-4 mt-2">
                    <div className="flex items-center gap-1.5 text-xs text-slate-500">
                      <User size={11} />
                      <span className="truncate max-w-[180px]">{event.actor}</span>
                    </div>
                    <div className="flex items-center gap-1.5 text-xs text-slate-600 font-mono">
                      <span>{event.ipAddress}</span>
                    </div>
                    <div className="ml-auto opacity-0 group-hover:opacity-100 transition-opacity">
                      <button className="flex items-center gap-1 text-xs text-[#1565C0] hover:text-blue-300 focus:outline-none">
                        <Eye size={12} />
                        {lang==='ht'?'Detay':'Détails'}
                        <ChevronRight size={12} />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between px-4 py-3 border-t border-slate-700/40 text-xs text-slate-500">
          <div className="flex items-center gap-2">
            <Lock size={12} className="text-emerald-400" />
            <span>
              {lang === 'ht'
                ? `${filtered.length} evènman — Merkle Tree verifye — WORM imyab`
                : `${filtered.length} événements — Merkle Tree vérifié — WORM immuable`}
            </span>
          </div>
          <div className="text-slate-600">
            {lang === 'ht' ? 'Retansyon: Hot 30j / Warm 1an / Cold 10an' : 'Rétention: Hot 30j / Warm 1an / Cold 10 ans'}
          </div>
        </div>
      </div>

      {/* Event detail modal */}
      {selectedEvent && <EventDetailModal event={selectedEvent} onClose={() => setSelectedEvent(null)} lang={lang} />}
    </div>
  );
};
