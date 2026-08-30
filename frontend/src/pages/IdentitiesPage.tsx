import React, { useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Search, Filter, Download, RefreshCw, Eye, ChevronRight,
  UserCheck, UserX, AlertTriangle, Fingerprint, Calendar,
  MapPin, Phone, Shield, Clock, ChevronDown, X
} from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

// ─── Types ──────────────────────────────────────────────────────────────────

type BiometricStatus = 'COMPLET' | 'PARTIEL' | 'MANQUANT';
type IdentityStatus  = 'VERIFIE' | 'SUSPECT' | 'EN_ATTENTE' | 'BLOQUE';
type AalLevel        = 'AAL1' | 'AAL2' | 'AAL3';

interface Identity {
  id: string;
  nni: string;
  firstName: string;
  lastName: string;
  dob: string;
  commune: string;
  department: string;
  phone?: string;
  status: IdentityStatus;
  biometric: BiometricStatus;
  aal: AalLevel;
  enrolledAt: string;
  lastVerified: string;
  riskScore: number;
}

// ─── Mock data ───────────────────────────────────────────────────────────────

const MOCK_IDENTITIES: Identity[] = [
  { id: '1', nni: 'HTI-2024-001-00001', firstName: 'Marie', lastName: 'Joseph', dob: '1985-03-12', commune: 'Port-au-Prince', department: 'Ouest', phone: '+509 3712-4455', status: 'VERIFIE', biometric: 'COMPLET', aal: 'AAL2', enrolledAt: '2024-01-15', lastVerified: '2024-05-10', riskScore: 0.04 },
  { id: '2', nni: 'HTI-2024-001-00002', firstName: 'Jean-Baptiste', lastName: 'Pierre', dob: '1990-07-22', commune: 'Pétionville', department: 'Ouest', phone: '+509 3801-2233', status: 'VERIFIE', biometric: 'COMPLET', aal: 'AAL2', enrolledAt: '2024-01-16', lastVerified: '2024-05-11', riskScore: 0.02 },
  { id: '3', nni: 'HTI-2024-002-00015', firstName: 'Claudette', lastName: 'Morisseau', dob: '1978-11-05', commune: 'Cap-Haïtien', department: 'Nord', phone: undefined, status: 'EN_ATTENTE', biometric: 'PARTIEL', aal: 'AAL1', enrolledAt: '2024-02-03', lastVerified: '2024-04-20', riskScore: 0.18 },
  { id: '4', nni: 'HTI-2024-003-00007', firstName: 'Pierre-Louis', lastName: 'Desroches', dob: '1965-09-18', commune: 'Les Cayes', department: 'Sud', phone: '+509 3645-7788', status: 'SUSPECT', biometric: 'COMPLET', aal: 'AAL1', enrolledAt: '2024-03-10', lastVerified: '2024-05-01', riskScore: 0.72 },
  { id: '5', nni: 'HTI-2024-004-00032', firstName: 'Roseline', lastName: 'Étienne', dob: '1995-01-30', commune: 'Jacmel', department: 'Sud-Est', phone: '+509 3990-1122', status: 'VERIFIE', biometric: 'COMPLET', aal: 'AAL3', enrolledAt: '2024-04-05', lastVerified: '2024-05-12', riskScore: 0.01 },
  { id: '6', nni: 'HTI-2024-001-00088', firstName: 'François', lastName: 'Lafortune', dob: '1972-06-14', commune: 'Delmas', department: 'Ouest', phone: '+509 3723-9900', status: 'BLOQUE', biometric: 'COMPLET', aal: 'AAL2', enrolledAt: '2024-01-28', lastVerified: '2024-03-15', riskScore: 0.91 },
  { id: '7', nni: 'HTI-2024-005-00003', firstName: 'Natacha', lastName: 'Belizaire', dob: '2000-04-08', commune: 'Gonaïves', department: 'Artibonite', phone: '+509 3855-4433', status: 'VERIFIE', biometric: 'COMPLET', aal: 'AAL2', enrolledAt: '2024-02-20', lastVerified: '2024-05-09', riskScore: 0.05 },
  { id: '8', nni: 'HTI-2024-006-00011', firstName: 'Emmanuel', lastName: 'Clermont', dob: '1988-12-25', commune: 'Fort-Liberté', department: 'Nord-Est', phone: undefined, status: 'EN_ATTENTE', biometric: 'MANQUANT', aal: 'AAL1', enrolledAt: '2024-03-30', lastVerified: '2024-04-02', riskScore: 0.32 },
];

// ─── Helpers ─────────────────────────────────────────────────────────────────

const STATUS_CONFIG: Record<IdentityStatus, { label: string; labelHT: string; color: string; icon: LucideIcon }> = {
  VERIFIE:    { label: 'Vérifié',    labelHT: 'Verifye',    color: 'text-emerald-400 bg-emerald-400/10 border-emerald-400/30', icon: UserCheck },
  SUSPECT:    { label: 'Suspect',    labelHT: 'Sispèk',     color: 'text-orange-400 bg-orange-400/10 border-orange-400/30',   icon: AlertTriangle },
  EN_ATTENTE: { label: 'En attente', labelHT: 'Ap tann',    color: 'text-yellow-400 bg-yellow-400/10 border-yellow-400/30',   icon: Clock },
  BLOQUE:     { label: 'Bloqué',     labelHT: 'Bloke',      color: 'text-red-400 bg-red-400/10 border-red-400/30',            icon: UserX },
};

const BIO_CONFIG: Record<BiometricStatus, { label: string; color: string }> = {
  COMPLET:  { label: 'Complet',  color: 'text-emerald-400' },
  PARTIEL:  { label: 'Partiel',  color: 'text-yellow-400'  },
  MANQUANT: { label: 'Manquant', color: 'text-red-400'     },
};

const AAL_COLOR: Record<AalLevel, string> = {
  AAL1: 'bg-slate-700 text-slate-300',
  AAL2: 'bg-blue-900/50 text-blue-300 border border-blue-700/40',
  AAL3: 'bg-purple-900/50 text-purple-300 border border-purple-700/40',
};

const getRiskColor = (score: number) => {
  if (score < 0.2)  return 'text-emerald-400';
  if (score < 0.5)  return 'text-yellow-400';
  if (score < 0.75) return 'text-orange-400';
  return 'text-red-400';
};

// ─── Subcomponents ───────────────────────────────────────────────────────────

const StatusBadge: React.FC<{ status: IdentityStatus; lang: string }> = ({ status, lang }) => {
  const cfg = STATUS_CONFIG[status];
  const Icon = cfg.icon;
  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border ${cfg.color}`}>
      <Icon size={11} />
      {lang === 'ht' ? cfg.labelHT : cfg.label}
    </span>
  );
};

const DetailPanel: React.FC<{ identity: Identity; onClose: () => void; lang: string }> = ({ identity, onClose, lang }) => {
  const cfg = STATUS_CONFIG[identity.status];
  const StatusIcon = cfg.icon;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" role="dialog" aria-modal="true" aria-label="Détail identité">
      <div className="bg-[#0D1B2A] border border-slate-700/60 rounded-2xl w-full max-w-2xl mx-4 shadow-2xl shadow-black/40 overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-slate-700/60 bg-[#1E3A5F]/30">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-[#1565C0]/20 border border-[#1565C0]/40 flex items-center justify-center text-[#1565C0] font-bold text-sm">
              {identity.firstName[0]}{identity.lastName[0]}
            </div>
            <div>
              <h2 className="text-white font-semibold">{identity.lastName}, {identity.firstName}</h2>
              <p className="text-slate-400 text-xs font-mono">{identity.nni}</p>
            </div>
          </div>
          <button onClick={onClose} className="p-2 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-[#1565C0]" aria-label="Fermer">
            <X size={18} />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 grid grid-cols-2 gap-6">
          {/* Left column */}
          <div className="space-y-4">
            <Section title={lang === 'ht' ? 'Enfòmasyon Pèsonèl' : 'Informations Personnelles'}>
              <Row icon={<Calendar size={14} />} label={lang === 'ht' ? 'Dat nesans' : 'Date de naissance'} value={new Date(identity.dob).toLocaleDateString('fr-HT')} />
              <Row icon={<MapPin size={14} />} label={lang === 'ht' ? 'Komin' : 'Commune'} value={`${identity.commune}, ${identity.department}`} />
              {identity.phone && <Row icon={<Phone size={14} />} label={lang === 'ht' ? 'Telefòn' : 'Téléphone'} value={identity.phone} />}
            </Section>

            <Section title={lang === 'ht' ? 'Byometrik' : 'Biométrie'}>
              <Row icon={<Fingerprint size={14} />} label={lang === 'ht' ? 'Eta' : 'Statut'} value={
                <span className={`${BIO_CONFIG[identity.biometric].color} font-medium`}>{BIO_CONFIG[identity.biometric].label}</span>
              } />
            </Section>
          </div>

          {/* Right column */}
          <div className="space-y-4">
            <Section title={lang === 'ht' ? 'Sekirite ak Konfyans' : 'Sécurité & Confiance'}>
              <Row icon={<Shield size={14} />} label="AAL" value={
                <span className={`text-xs px-2 py-0.5 rounded font-medium ${AAL_COLOR[identity.aal]}`}>{identity.aal}</span>
              } />
              <Row icon={<AlertTriangle size={14} />} label={lang === 'ht' ? 'Nivo risk' : 'Score de risque'} value={
                <span className={`font-mono font-bold ${getRiskColor(identity.riskScore)}`}>{(identity.riskScore * 100).toFixed(0)}%</span>
              } />
              <Row icon={<StatusIcon size={14} />} label={lang === 'ht' ? 'Eta Idantite' : 'Statut Identité'} value={
                <StatusBadge status={identity.status} lang={lang} />
              } />
            </Section>

            <Section title={lang === 'ht' ? 'Istwa' : 'Historique'}>
              <Row icon={<Calendar size={14} />} label={lang === 'ht' ? 'Enrejistre' : 'Enrôlé le'} value={new Date(identity.enrolledAt).toLocaleDateString('fr-HT')} />
              <Row icon={<Clock size={14} />} label={lang === 'ht' ? 'Dènye verifikasyon' : 'Dernière vérification'} value={new Date(identity.lastVerified).toLocaleDateString('fr-HT')} />
            </Section>
          </div>
        </div>

        {/* Risk bar */}
        <div className="px-6 pb-4">
          <div className="flex items-center justify-between text-xs text-slate-400 mb-1.5">
            <span>{lang === 'ht' ? 'Baw risk' : 'Barre de risque'}</span>
            <span className={`font-mono font-bold ${getRiskColor(identity.riskScore)}`}>{(identity.riskScore * 100).toFixed(1)}%</span>
          </div>
          <div className="h-2 bg-slate-800 rounded-full overflow-hidden">
            <div
              className="h-full rounded-full transition-all"
              style={{
                width: `${identity.riskScore * 100}%`,
                background: identity.riskScore < 0.2 ? '#2E7D32' : identity.riskScore < 0.5 ? '#E65100' : '#C62828'
              }}
            />
          </div>
        </div>

        {/* Footer actions */}
        <div className="flex justify-end gap-3 px-6 py-4 border-t border-slate-700/60 bg-slate-900/30">
          <button onClick={onClose} className="px-4 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 transition-colors focus:outline-none focus:ring-2 focus:ring-slate-600">
            {lang === 'ht' ? 'Fèmen' : 'Fermer'}
          </button>
          {identity.status === 'SUSPECT' && (
            <button className="px-4 py-2 rounded-lg text-sm bg-orange-600/20 text-orange-400 border border-orange-600/30 hover:bg-orange-600/30 transition-colors focus:outline-none focus:ring-2 focus:ring-orange-500">
              {lang === 'ht' ? 'Voye bay DCPJ' : 'Transmettre DCPJ'}
            </button>
          )}
          {identity.status !== 'BLOQUE' && (
            <button className="px-4 py-2 rounded-lg text-sm bg-[#1565C0] text-white hover:bg-blue-600 transition-colors focus:outline-none focus:ring-2 focus:ring-[#1565C0]">
              {lang === 'ht' ? 'Verifye' : 'Vérifier'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

const Section: React.FC<{ title: string; children: React.ReactNode }> = ({ title, children }) => (
  <div>
    <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">{title}</h3>
    <div className="space-y-2.5">{children}</div>
  </div>
);

const Row: React.FC<{ icon: React.ReactNode; label: string; value: React.ReactNode }> = ({ icon, label, value }) => (
  <div className="flex items-center justify-between gap-2">
    <div className="flex items-center gap-1.5 text-slate-400 text-xs min-w-0">
      {icon}
      <span className="truncate">{label}</span>
    </div>
    <div className="text-slate-200 text-xs text-right">{value}</div>
  </div>
);

// ─── Main Component ───────────────────────────────────────────────────────────

export const IdentitiesPage: React.FC = () => {
  const { t, i18n } = useTranslation();
  const lang = i18n.language;

  const [search, setSearch]           = useState('');
  const [statusFilter, setStatusFilter] = useState<IdentityStatus | 'ALL'>('ALL');
  const [bioFilter, setBioFilter]     = useState<BiometricStatus | 'ALL'>('ALL');
  const [sortField, setSortField]     = useState<keyof Identity>('lastName');
  const [sortAsc, setSortAsc]         = useState(true);
  const [selected, setSelected]       = useState<Identity | null>(null);
  const [showFilters, setShowFilters] = useState(false);

  const filtered = useMemo(() => {
    let list = [...MOCK_IDENTITIES];

    // Search across NNI, name, commune
    if (search.trim()) {
      const q = search.toLowerCase();
      list = list.filter(id =>
        id.nni.toLowerCase().includes(q) ||
        id.firstName.toLowerCase().includes(q) ||
        id.lastName.toLowerCase().includes(q) ||
        id.commune.toLowerCase().includes(q)
      );
    }

    if (statusFilter !== 'ALL') list = list.filter(id => id.status === statusFilter);
    if (bioFilter !== 'ALL')    list = list.filter(id => id.biometric === bioFilter);

    list.sort((a, b) => {
      const va = String(a[sortField]);
      const vb = String(b[sortField]);
      return sortAsc ? va.localeCompare(vb) : vb.localeCompare(va);
    });

    return list;
  }, [search, statusFilter, bioFilter, sortField, sortAsc]);

  const toggleSort = (field: keyof Identity) => {
    if (sortField === field) setSortAsc(!sortAsc);
    else { setSortField(field); setSortAsc(true); }
  };

  const stats = useMemo(() => ({
    total:     MOCK_IDENTITIES.length,
    verified:  MOCK_IDENTITIES.filter(i => i.status === 'VERIFIE').length,
    suspect:   MOCK_IDENTITIES.filter(i => i.status === 'SUSPECT').length,
    blocked:   MOCK_IDENTITIES.filter(i => i.status === 'BLOQUE').length,
    pending:   MOCK_IDENTITIES.filter(i => i.status === 'EN_ATTENTE').length,
  }), []);

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white">
            {lang === 'ht' ? 'Jèsyon Idantite' : t('nav.identities')}
          </h2>
          <p className="text-slate-400 text-sm mt-0.5">
            {lang === 'ht'
              ? `${MOCK_IDENTITIES.length.toLocaleString()} sitwayen anrejistre`
              : `${MOCK_IDENTITIES.length.toLocaleString()} citoyens enregistrés`}
          </p>
        </div>
        <div className="flex items-center gap-3">
          <button className="flex items-center gap-2 px-4 py-2 text-sm text-slate-300 bg-slate-800 hover:bg-slate-700 border border-slate-700 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-slate-600">
            <Download size={15} />
            {lang === 'ht' ? 'Ekspòte' : 'Exporter'}
          </button>
          <button className="flex items-center gap-2 px-4 py-2 text-sm text-white bg-[#1565C0] hover:bg-blue-600 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-[#1565C0]">
            <RefreshCw size={15} />
            {lang === 'ht' ? 'Rafraîchi' : 'Actualiser'}
          </button>
        </div>
      </div>

      {/* Stats Row */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        {[
          { label: lang === 'ht' ? 'Total' : 'Total',         value: stats.total,    color: 'text-slate-200'    },
          { label: lang === 'ht' ? 'Verifye' : 'Vérifiés',    value: stats.verified, color: 'text-emerald-400'  },
          { label: lang === 'ht' ? 'Sispèk' : 'Suspects',     value: stats.suspect,  color: 'text-orange-400'   },
          { label: lang === 'ht' ? 'Bloke' : 'Bloqués',       value: stats.blocked,  color: 'text-red-400'      },
        ].map(s => (
          <div key={s.label} className="bg-[#1E3A5F]/20 border border-slate-700/40 rounded-xl p-4 text-center">
            <div className={`text-2xl font-bold ${s.color}`}>{s.value}</div>
            <div className="text-slate-400 text-xs mt-1">{s.label}</div>
          </div>
        ))}
      </div>

      {/* Search & Filters */}
      <div className="bg-slate-900/60 border border-slate-700/40 rounded-xl p-4 space-y-3">
        <div className="flex gap-3">
          {/* Search bar */}
          <div className="relative flex-1">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              id="identity-search"
              type="search"
              placeholder={lang === 'ht' ? 'Chèche pa NNI, non, komin...' : 'Rechercher par NNI, nom, commune...'}
              value={search}
              onChange={e => setSearch(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-lg pl-9 pr-4 py-2.5 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-[#1565C0] focus:border-transparent"
            />
          </div>
          {/* Filter toggle */}
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`flex items-center gap-2 px-4 py-2.5 text-sm rounded-lg border transition-colors focus:outline-none focus:ring-2 focus:ring-[#1565C0] ${showFilters ? 'bg-[#1565C0]/20 border-[#1565C0]/40 text-blue-400' : 'bg-slate-800 border-slate-700 text-slate-300 hover:bg-slate-700'}`}
          >
            <Filter size={15} />
            {lang === 'ht' ? 'Filtre' : 'Filtres'}
            <ChevronDown size={14} className={`transition-transform ${showFilters ? 'rotate-180' : ''}`} />
          </button>
        </div>

        {/* Expanded filters */}
        {showFilters && (
          <div className="flex flex-wrap gap-3 pt-2 border-t border-slate-700/40">
            <div className="flex items-center gap-2">
              <label className="text-xs text-slate-400">{lang === 'ht' ? 'Eta:' : 'Statut:'}</label>
              <select
                value={statusFilter}
                onChange={e => setStatusFilter(e.target.value as IdentityStatus | 'ALL')}
                className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-[#1565C0]"
              >
                <option value="ALL">{lang === 'ht' ? 'Tout' : 'Tous'}</option>
                <option value="VERIFIE">{lang === 'ht' ? 'Verifye' : 'Vérifié'}</option>
                <option value="SUSPECT">{lang === 'ht' ? 'Sispèk' : 'Suspect'}</option>
                <option value="EN_ATTENTE">{lang === 'ht' ? 'Ap tann' : 'En attente'}</option>
                <option value="BLOQUE">{lang === 'ht' ? 'Bloke' : 'Bloqué'}</option>
              </select>
            </div>
            <div className="flex items-center gap-2">
              <label className="text-xs text-slate-400">{lang === 'ht' ? 'Byometrik:' : 'Biométrie:'}</label>
              <select
                value={bioFilter}
                onChange={e => setBioFilter(e.target.value as BiometricStatus | 'ALL')}
                className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-[#1565C0]"
              >
                <option value="ALL">{lang === 'ht' ? 'Tout' : 'Tous'}</option>
                <option value="COMPLET">{lang === 'ht' ? 'Konplè' : 'Complet'}</option>
                <option value="PARTIEL">{lang === 'ht' ? 'Pasyèl' : 'Partiel'}</option>
                <option value="MANQUANT">{lang === 'ht' ? 'Manke' : 'Manquant'}</option>
              </select>
            </div>
            {(statusFilter !== 'ALL' || bioFilter !== 'ALL' || search) && (
              <button
                onClick={() => { setStatusFilter('ALL'); setBioFilter('ALL'); setSearch(''); }}
                className="flex items-center gap-1.5 text-xs text-red-400 hover:text-red-300 transition-colors focus:outline-none"
              >
                <X size={12} />
                {lang === 'ht' ? 'Efase filtre' : 'Réinitialiser'}
              </button>
            )}
          </div>
        )}
      </div>

      {/* Table */}
      <div className="bg-slate-900/60 border border-slate-700/40 rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-700/40 text-slate-400 text-xs uppercase tracking-wider">
                <Th label="NNI"                      field="nni"          sort={sortField} asc={sortAsc} onClick={toggleSort} />
                <Th label={lang==='ht'?'Non':'Nom'}  field="lastName"     sort={sortField} asc={sortAsc} onClick={toggleSort} />
                <Th label={lang==='ht'?'Komin':'Commune'} field="commune" sort={sortField} asc={sortAsc} onClick={toggleSort} />
                <Th label={lang==='ht'?'Eta':'Statut'} field="status"     sort={sortField} asc={sortAsc} onClick={toggleSort} />
                <Th label={lang==='ht'?'Byometrik':'Biométrie'} field="biometric" sort={sortField} asc={sortAsc} onClick={toggleSort} />
                <th className="px-4 py-3 text-left">AAL</th>
                <Th label={lang==='ht'?'Risk':'Risque'} field="riskScore" sort={sortField} asc={sortAsc} onClick={toggleSort} />
                <th className="px-4 py-3 text-right">{lang==='ht'?'Aksyon':'Actions'}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-700/30">
              {filtered.length === 0 ? (
                <tr>
                  <td colSpan={8} className="px-6 py-12 text-center text-slate-500">
                    {lang === 'ht' ? 'Pa gen rezilta' : 'Aucun résultat'}
                  </td>
                </tr>
              ) : filtered.map(identity => (
                <tr
                  key={identity.id}
                  className="hover:bg-slate-800/30 transition-colors cursor-pointer group"
                  onClick={() => setSelected(identity)}
                >
                  <td className="px-4 py-3">
                    <span className="font-mono text-xs text-slate-300">{identity.nni}</span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2.5">
                      <div className="w-7 h-7 rounded-full bg-[#1E3A5F] border border-slate-600/50 flex items-center justify-center text-xs font-medium text-blue-300 flex-shrink-0">
                        {identity.firstName[0]}{identity.lastName[0]}
                      </div>
                      <div>
                        <div className="font-medium text-slate-100">{identity.lastName}</div>
                        <div className="text-slate-500 text-xs">{identity.firstName}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-slate-400 text-xs">
                    <div>{identity.commune}</div>
                    <div className="text-slate-600">{identity.department}</div>
                  </td>
                  <td className="px-4 py-3">
                    <StatusBadge status={identity.status} lang={lang} />
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-1.5">
                      <Fingerprint size={13} className={BIO_CONFIG[identity.biometric].color} />
                      <span className={`text-xs ${BIO_CONFIG[identity.biometric].color}`}>{BIO_CONFIG[identity.biometric].label}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <span className={`text-xs px-2 py-0.5 rounded font-medium ${AAL_COLOR[identity.aal]}`}>{identity.aal}</span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <div className="h-1.5 w-16 bg-slate-700 rounded-full overflow-hidden">
                        <div
                          className="h-full rounded-full"
                          style={{
                            width: `${identity.riskScore * 100}%`,
                            backgroundColor: identity.riskScore < 0.2 ? '#2E7D32' : identity.riskScore < 0.5 ? '#E65100' : '#C62828'
                          }}
                        />
                      </div>
                      <span className={`text-xs font-mono ${getRiskColor(identity.riskScore)}`}>
                        {(identity.riskScore * 100).toFixed(0)}%
                      </span>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      className="p-1.5 rounded-lg text-slate-500 hover:text-white hover:bg-slate-700 transition-colors opacity-0 group-hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-[#1565C0] focus:opacity-100"
                      aria-label="Voir détails"
                      onClick={e => { e.stopPropagation(); setSelected(identity); }}
                    >
                      <Eye size={15} />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Table footer */}
        <div className="flex items-center justify-between px-4 py-3 border-t border-slate-700/40 text-xs text-slate-500">
          <span>
            {lang === 'ht'
              ? `Montre ${filtered.length} nan ${MOCK_IDENTITIES.length} idantite`
              : `Affichage ${filtered.length} sur ${MOCK_IDENTITIES.length} identités`}
          </span>
          <div className="flex items-center gap-1 text-slate-400">
            <Shield size={12} className="text-[#1565C0]" />
            <span>{lang === 'ht' ? 'Jwenn enpòtasyon ABIS' : 'Données ABIS synchronisées'}</span>
          </div>
        </div>
      </div>

      {/* Detail panel modal */}
      {selected && <DetailPanel identity={selected} onClose={() => setSelected(null)} lang={lang} />}
    </div>
  );
};

// Sort column header
const Th: React.FC<{
  label: string; field: keyof Identity;
  sort: keyof Identity; asc: boolean;
  onClick: (f: keyof Identity) => void;
}> = ({ label, field, sort, asc, onClick }) => (
  <th
    className="px-4 py-3 text-left cursor-pointer hover:text-slate-200 select-none transition-colors"
    onClick={() => onClick(field)}
  >
    <div className="flex items-center gap-1">
      {label}
      <ChevronRight
        size={12}
        className={`transition-transform ${sort === field ? (asc ? 'rotate-90' : '-rotate-90') : 'opacity-30'}`}
      />
    </div>
  </th>
);
