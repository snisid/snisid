import React, { useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type { LucideIcon } from 'lucide-react';
import {
  Search, BookOpen, ChevronDown, ChevronRight,
  Shield, Server, Cpu, Wifi, Key, Lock,
  Globe, Zap, FileCheck, X, Info, Award
} from 'lucide-react';


// ─── Types ──────────────────────────────────────────────────────────────────

interface GlossaryTerm {
  id: string;
  term: string;
  fullName: string;
  definition: string;
  version: string;
  updatedAt: string;
  masterPrompts: string[];
  category: TermCategory;
  pillar?: string;
  seuils?: string[];
  metriques?: string[];
  critical?: boolean;
}

interface Standard {
  id: string;
  code: string;
  domain: string;
  applicability: string;
  phase: string;
  link?: string;
}

interface MaturityScore {
  domain: string;
  current: number;
  projected: number;
  recommendation: string;
}

type TermCategory = 'BIOMETRIE' | 'SECURITE' | 'INFRASTRUCTURE' | 'GOUVERNANCE' | 'CRYPTO' | 'RESEAU';

// ─── Data ────────────────────────────────────────────────────────────────────

const TERMS: GlossaryTerm[] = [
  {
    id: 'abis', term: 'ABIS', fullName: 'Automated Biometric Identification System',
    definition: 'Moteur de déduplication biométrique 1:N fonctionnant sur cluster GPU certifié. Effectue la comparaison exhaustive contre les 15 millions d\'identités enregistrées.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-005', 'MP-006', 'MP-013'],
    category: 'BIOMETRIE', pillar: 'Pilier 1 — IDENTITÉ', critical: true,
    seuils: ['Score < 85% → Nouveau citoyen → NNI attribué', 'Score 85–95% → Tier-2 humain', 'Score ≥ 95% → Hard Match → DCPJ forensique'],
    metriques: ['Déduplication 1:N (15M) : < 30s', 'FAR : < 0.001%', 'FRR : < 1%']
  },
  {
    id: 'abac', term: 'ABAC', fullName: 'Attribute-Based Access Control',
    definition: 'Modèle de contrôle d\'accès basé sur les attributs contextuels (identité, rôle, localisation, appareil, heure). Implémenté via OPA/Rego en sidecar dans chaque microservice SNISID.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-001', 'MP-004'],
    category: 'SECURITE', pillar: 'Pilier 4 — APPLICATION'
  },
  {
    id: 'cockroachdb', term: 'CockroachDB', fullName: 'Base SQL Distribuée ACID Multi-Région',
    definition: 'Base SQL distribuée ACID, multi-région, compatible PostgreSQL. Garantit la cohérence des données citoyennes entre DC1 (Port-au-Prince) et DC2 (Cap-Haïtien).',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-007', 'MP-010'],
    category: 'INFRASTRUCTURE', pillar: 'Pilier 5 — DONNÉES'
  },
  {
    id: 'crdt', term: 'CRDT', fullName: 'Conflict-free Replicated Data Type',
    definition: 'Structure de données garantie de converger vers un état cohérent lors de la resynchronisation après une période hors-ligne. Utilisé dans les MEK pour résoudre automatiquement les conflits.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-008', 'MP-010'],
    category: 'INFRASTRUCTURE'
  },
  {
    id: 'ejbca', term: 'EJBCA', fullName: 'Enterprise JavaBeans Certificate Authority',
    definition: 'PKI open-source gouvernementale utilisée comme Issuing CA (niveaux 3A et 3B) pour l\'émission automatisée des certificats eID citoyens (~500K/an) et des certificats TLS microservices (24h via ACME).',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-009', 'MP-013'],
    category: 'CRYPTO'
  },
  {
    id: 'fido2', term: 'FIDO2', fullName: 'Fast Identity Online 2 — W3C WebAuthn',
    definition: 'Standard W3C d\'authentification sans mot de passe basé sur la cryptographie à clé publique. Implémenté via YubiKey pour les agents ONI, DGI, DCPJ (AAL3) et facteur fort opérateurs (AAL2).',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-001', 'MP-013'],
    category: 'SECURITE', pillar: 'AAL3 — HAUT'
  },
  {
    id: 'hsm', term: 'HSM', fullName: 'Hardware Security Module',
    definition: 'Module matériel dédié à la protection et à l\'exécution des opérations cryptographiques. Les clés privées ne quittent JAMAIS le HSM.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-009', 'MP-013'],
    category: 'CRYPTO', critical: true,
    metriques: ['Niveau requis: FIPS 140-2 L3 minimum', 'Root CA: FIPS 140-2 L4 (Thales Luna)', 'Quorum 5-of-9 pour Root CA']
  },
  {
    id: 'mek', term: 'MEK', fullName: 'Mobile Enrollment Kit',
    definition: 'Kit biométrique ruggedisé autonome pour l\'enrôlement terrain en zones rurales haïtiennes sans connectivité permanente. Solaire 200W, batterie LiFePO4, 72h d\'autonomie, IP67.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-008', 'MP-010', 'MP-013'],
    category: 'INFRASTRUCTURE',
    metriques: ['Autonomie: 72h sans soleil', 'Connectivité: 4G + Starlink backup', 'Stockage: NVMe 2TB LUKS + TPM 2.0', 'IP67 — Étanche']
  },
  {
    id: 'mtls', term: 'mTLS', fullName: 'Mutual TLS',
    definition: 'Extension du TLS où les deux parties (client ET serveur) s\'authentifient mutuellement via certificats X.509. Fondamental dans l\'architecture Zero Trust pour éliminer la confiance implicite.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-001', 'MP-003'],
    category: 'RESEAU', pillar: 'Pilier 3 — RÉSEAU'
  },
  {
    id: 'nats', term: 'NATS JetStream', fullName: 'Message Broker Léger Orienté Edge',
    definition: 'Message broker léger orienté edge computing avec persistence locale. Permet la file d\'attente de messages biométriques en zone hors-ligne et la resynchronisation automatique.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-008', 'MP-010'],
    category: 'INFRASTRUCTURE'
  },
  {
    id: 'nni', term: 'NNI', fullName: 'Numéro National d\'Identification',
    definition: 'Identifiant unique, pérenne et immuable attribué à chaque citoyen haïtien lors de l\'enrôlement biométrique. Attribué uniquement après déduplication ABIS 1:N réussie (score < 85%).',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-002', 'MP-005', 'MP-006', 'MP-013'],
    category: 'GOUVERNANCE', critical: true,
    metriques: ['Format: HTI-AAAA-DPT-NNNNN', 'Ex: HTI-2024-001-00001', 'Cadre légal requis avant Phase 2']
  },
  {
    id: 'oni', term: 'ONI', fullName: 'Office National d\'Identification',
    definition: 'Agence gouvernementale haïtienne responsable de l\'identité civile nationale. Opérateur principal du SNISID, supervise l\'attribution des NNI et l\'enrôlement des agents (AAL3).',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-002', 'MP-011', 'MP-013'],
    category: 'GOUVERNANCE'
  },
  {
    id: 'opa', term: 'OPA', fullName: 'Open Policy Agent',
    definition: 'Moteur de politique open-source (CNCF Graduated) utilisant le langage déclaratif Rego pour évaluer des décisions d\'autorisation ABAC. Déployé en sidecar dans chaque pod Kubernetes SNISID.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-001', 'MP-004'],
    category: 'SECURITE', pillar: 'Pilier 4 — APPLICATION'
  },
  {
    id: 'pad', term: 'PAD', fullName: 'Presentation Attack Detection',
    definition: 'Système de détection des tentatives de fraude biométrique (spoofing). Analyse les captures pour détecter faux doigts, photos, masques 3D, deepfakes vidéo en temps réel.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-005', 'MP-006'],
    category: 'BIOMETRIE', critical: true
  },
  {
    id: 'rke2', term: 'RKE2', fullName: 'Rancher Kubernetes Engine 2',
    definition: 'Distribution Kubernetes sécurisée par défaut, compatible FIPS 140-2. Utilisée pour les clusters de production SNISID (DC1 + DC2) avec hardening CIS Benchmark Level 2 intégré.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-003', 'MP-007'],
    category: 'INFRASTRUCTURE'
  },
  {
    id: 'rpo', term: 'RPO', fullName: 'Recovery Point Objective',
    definition: 'Perte de données maximale acceptable en cas d\'incident. Définit l\'intervalle maximum entre le dernier backup opérationnel et le moment de la panne.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-010', 'MP-013'],
    category: 'GOUVERNANCE', critical: true,
    metriques: ['SNISID : < 1 minute', 'Réplication synchrone CockroachDB DC1↔DC2']
  },
  {
    id: 'rto', term: 'RTO', fullName: 'Recovery Time Objective',
    definition: 'Durée maximale tolérée pour le rétablissement du service après un incident. Inclut la détection, le failover automatique et la reprise opérationnelle.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-010', 'MP-013'],
    category: 'GOUVERNANCE', critical: true,
    metriques: ['SNISID : < 15 minutes', 'Failover automatique DC1 → DC2 via ArgoCD']
  },
  {
    id: 'spiffe', term: 'SPIFFE', fullName: 'Secure Production Identity Framework for Everyone',
    definition: 'Framework CNCF pour l\'attribution d\'identités cryptographiques aux microservices via des SVIDs (SPIFFE Verifiable Identity Documents). Implémenté via SPIRE (agent/serveur).',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-001', 'MP-003'],
    category: 'SECURITE', pillar: 'Pilier 1 — IDENTITÉ'
  },
  {
    id: 'worm', term: 'WORM', fullName: 'Write-Once-Read-Many',
    definition: 'Principe de stockage immuable — une fois écrit, un enregistrement ne peut jamais être modifié ni supprimé. Obligatoire pour les logs d\'audit légaux SNISID (rétention 10 ans).',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-004', 'MP-013'],
    category: 'SECURITE',
    metriques: ['Hot: 30 jours (SSD Elasticsearch)', 'Warm: 1 an (Object Storage)', 'Cold/WORM: 10 ans (Air-gappé immuable)']
  },
  {
    id: 'xroad', term: 'X-Road', fullName: 'Standard Interopérabilité Estonien',
    definition: 'Protocole d\'échange de données inter-agences gouvernementales (Estonie), adopté par 30+ pays. Garantit confidentialité, intégrité et traçabilité de chaque échange inter-agences.',
    version: 'v1.0', updatedAt: '2026-05-24', masterPrompts: ['MP-003', 'MP-011', 'MP-013'],
    category: 'RESEAU', critical: true,
    metriques: ['Agences: ONI, DGI, DCPJ, ANH', 'Prérequis: Cadre légal avant Phase 3', '30+ pays déployés']
  },
];

const STANDARDS: Standard[] = [
  { id: 's1', code: 'NIST SP 800-63-3', domain: 'Digital Identity Guidelines', applicability: 'IAM complet AAL1/2/3', phase: 'Phase 1' },
  { id: 's2', code: 'NIST SP 800-207', domain: 'Zero Trust Architecture', applicability: 'Architecture réseau 7 piliers', phase: 'Phase 1' },
  { id: 's3', code: 'ISO/IEC 27001:2022', domain: 'Information Security Management', applicability: 'Certification complète', phase: 'Phase 5' },
  { id: 's4', code: 'ISO/IEC 27701:2019', domain: 'Privacy Information Management', applicability: 'Données biométriques citoyens', phase: 'Phase 2' },
  { id: 's5', code: 'ISO/IEC 19794-2', domain: 'Biometric Data Format Fingerprint', applicability: 'ABIS + MEK templates', phase: 'Phase 2' },
  { id: 's6', code: 'FIPS 140-2', domain: 'Cryptographic Modules Security', applicability: 'HSM L3 min + PKI + MEK', phase: 'Phase 1' },
  { id: 's7', code: 'X-Road Protocol', domain: 'Estonian Interoperability', applicability: 'Échanges inter-agences', phase: 'Phase 3' },
  { id: 's8', code: 'OWASP API Top 10', domain: 'API Security Guidelines', applicability: 'API Gateway Kong WAF', phase: 'Phase 1' },
  { id: 's9', code: 'MITRE ATT&CK', domain: 'Adversary Tactics & Techniques', applicability: 'SOC + SIEM playbooks', phase: 'Phase 2' },
  { id: 's10', code: 'CNCF Landscape', domain: 'Cloud Native Technologies', applicability: 'Stack technique complète', phase: 'Phase 1' },
  { id: 's11', code: 'SLSA Framework', domain: 'Software Supply Chain Security', applicability: 'DevSecOps SLSA L3→4', phase: 'Phase 2' },
];

const MATURITY: MaturityScore[] = [
  { domain: 'Architecture Globale & Microservices', current: 87, projected: 97, recommendation: 'DDD + Event Sourcing complet' },
  { domain: 'Cybersécurité & Zero Trust',           current: 85, projected: 96, recommendation: 'Anti-Prompt-Injection Hardening' },
  { domain: 'IAM National',                          current: 83, projected: 95, recommendation: 'SPIFFE/SPIRE full mesh' },
  { domain: 'Biométrie ABIS',                        current: 81, projected: 94, recommendation: 'Détection deepfakes + drift monitoring' },
  { domain: 'SOC/SIEM/SOAR',                         current: 79, projected: 93, recommendation: 'Playbooks Git + Threat Intel CARICOM' },
  { domain: 'PKI & HSM',                             current: 82, projected: 96, recommendation: 'CRLs via USSD/SMS (Haïti-first)' },
  { domain: 'Offline-First & Résilience',            current: 88, projected: 98, recommendation: 'Cartographie MEK temps réel' },
  { domain: 'DevSecOps/GitOps',                      current: 76, projected: 93, recommendation: 'SBOM + environnements éphémères' },
  { domain: 'Gouvernance des Données',               current: 74, projected: 90, recommendation: 'Matrice RACI + MoU internationaux' },
  { domain: 'Disaster Recovery PRA/PCA',             current: 80, projected: 96, recommendation: 'Tests PRA trimestriels validés' },
];

// ─── Helpers ─────────────────────────────────────────────────────────────────

const CATEGORY_CONFIG: Record<TermCategory, { label: string; labelHT: string; color: string; icon: LucideIcon }> = {
  BIOMETRIE:      { label: 'Biométrie',       labelHT: 'Byometrik',     color: 'text-purple-400 bg-purple-400/10 border-purple-400/30', icon: Cpu },
  SECURITE:       { label: 'Sécurité',        labelHT: 'Sekirite',      color: 'text-red-400 bg-red-400/10 border-red-400/30',          icon: Shield },
  INFRASTRUCTURE: { label: 'Infrastructure',  labelHT: 'Enfrastrikti',  color: 'text-cyan-400 bg-cyan-400/10 border-cyan-400/30',       icon: Server },
  GOUVERNANCE:    { label: 'Gouvernance',     labelHT: 'Gouvènans',     color: 'text-orange-400 bg-orange-400/10 border-orange-400/30', icon: Globe },
  CRYPTO:         { label: 'Cryptographie',   labelHT: 'Kriptografi',   color: 'text-yellow-400 bg-yellow-400/10 border-yellow-400/30', icon: Key },
  RESEAU:         { label: 'Réseau',          labelHT: 'Rezo',          color: 'text-blue-400 bg-blue-400/10 border-blue-400/30',       icon: Wifi },
};

const PHASE_COLOR: Record<string, string> = {
  'Phase 1': 'text-emerald-400 bg-emerald-400/10',
  'Phase 2': 'text-blue-400 bg-blue-400/10',
  'Phase 3': 'text-yellow-400 bg-yellow-400/10',
  'Phase 5': 'text-purple-400 bg-purple-400/10',
};

// ─── Subcomponents ────────────────────────────────────────────────────────────

const CategoryBadge: React.FC<{ cat: TermCategory; lang: string }> = ({ cat, lang }) => {
  const cfg = CATEGORY_CONFIG[cat];
  const Icon = cfg.icon;
  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border ${cfg.color}`}>
      <Icon size={11} />
      {lang === 'ht' ? cfg.labelHT : cfg.label}
    </span>
  );
};

const MaturityBar: React.FC<{ current: number; projected: number; label: string; recommendation: string }> = ({ current, projected, label, recommendation }) => {
  const gain = projected - current;
  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between text-xs">
        <span className="text-slate-300 font-medium truncate max-w-[55%]">{label}</span>
        <div className="flex items-center gap-2 flex-shrink-0">
          <span className="text-slate-500">{current}</span>
          <span className="text-slate-600">→</span>
          <span className="text-emerald-400 font-bold">{projected}</span>
          <span className="text-emerald-500 text-xs font-bold">+{gain}</span>
        </div>
      </div>
      <div className="h-2 bg-slate-800 rounded-full overflow-hidden relative">
        {/* Current */}
        <div className="h-full bg-[#1565C0]/60 rounded-full" style={{ width: `${current}%` }} />
        {/* Projected overlay */}
        <div
          className="h-full bg-emerald-500/50 rounded-full absolute top-0"
          style={{ left: `${current}%`, width: `${gain}%` }}
        />
      </div>
      <div className="text-xs text-slate-600 truncate">💡 {recommendation}</div>
    </div>
  );
};

const TermCard: React.FC<{ term: GlossaryTerm; lang: string }> = ({ term, lang }) => {
  const [expanded, setExpanded] = useState(false);
  const cfg = CATEGORY_CONFIG[term.category];
  const Icon = cfg.icon;

  return (
    <div
      className={`bg-slate-900/60 border rounded-xl overflow-hidden transition-all duration-200 ${term.critical ? 'border-orange-500/30' : 'border-slate-700/40'} hover:border-slate-600/60`}
    >
      {/* Header — always visible */}
      <button
        className="w-full flex items-start gap-4 p-4 text-left hover:bg-slate-800/30 transition-colors focus:outline-none focus:ring-2 focus:ring-[#1565C0] focus:ring-inset"
        onClick={() => setExpanded(!expanded)}
        aria-expanded={expanded}
      >
        <div className={`p-2 rounded-lg border flex-shrink-0 mt-0.5 ${cfg.color.split(' ').slice(1).join(' ')}`}>
          <Icon size={16} className={cfg.color.split(' ')[0]} />
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="font-bold text-white font-mono">{term.term}</span>
            {term.critical && (
              <span className="text-xs px-1.5 py-0.5 rounded bg-orange-500/15 text-orange-400 border border-orange-500/25 font-medium">
                CRITIQUE
              </span>
            )}
            <span className="text-xs text-slate-600">{term.version}</span>
          </div>
          <p className="text-xs text-slate-400 mt-0.5 italic">{term.fullName}</p>
          <p className="text-sm text-slate-300 mt-1.5 line-clamp-2">{term.definition}</p>
        </div>

        <div className="flex items-center gap-2 flex-shrink-0">
          <CategoryBadge cat={term.category} lang={lang} />
          <div className={`p-1 rounded text-slate-400 transition-transform ${expanded ? 'rotate-180' : ''}`}>
            <ChevronDown size={16} />
          </div>
        </div>
      </button>

      {/* Expanded content */}
      {expanded && (
        <div className="border-t border-slate-700/40 px-4 py-4 space-y-4 bg-slate-900/30">
          {/* Master Prompts */}
          <div>
            <div className="text-xs text-slate-500 font-medium mb-2">Master Prompts</div>
            <div className="flex flex-wrap gap-1.5">
              {term.masterPrompts.map(mp => (
                <span key={mp} className="px-2 py-0.5 rounded bg-[#1565C0]/20 text-blue-300 border border-[#1565C0]/30 text-xs font-mono">
                  {mp}
                </span>
              ))}
            </div>
          </div>

          {/* Pilier Zero Trust */}
          {term.pillar && (
            <div>
              <div className="text-xs text-slate-500 font-medium mb-1">Pilier Zero Trust</div>
              <span className="text-xs text-[#00BCD4]">{term.pillar}</span>
            </div>
          )}

          {/* Seuils */}
          {term.seuils && (
            <div>
              <div className="text-xs text-slate-500 font-medium mb-2">Seuils Opérationnels</div>
              <div className="space-y-1.5">
                {term.seuils.map((s, i) => (
                  <div key={i} className="flex items-start gap-2 text-xs text-slate-300">
                    <ChevronRight size={12} className="text-[#1565C0] mt-0.5 flex-shrink-0" />
                    {s}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Métriques */}
          {term.metriques && (
            <div>
              <div className="text-xs text-slate-500 font-medium mb-2">Métriques / Spécifications</div>
              <div className="grid grid-cols-1 gap-1">
                {term.metriques.map((m, i) => (
                  <div key={i} className="flex items-center gap-2 text-xs text-slate-300 bg-slate-800/40 rounded px-2 py-1">
                    <Zap size={10} className="text-yellow-400 flex-shrink-0" />
                    {m}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Version footer */}
          <div className="flex items-center justify-between text-xs text-slate-600 pt-2 border-t border-slate-700/30">
            <span>Version {term.version} — Mis à jour {new Date(term.updatedAt).toLocaleDateString('fr-HT')}</span>
            <span>Conseil Technique SNISID</span>
          </div>
        </div>
      )}
    </div>
  );
};

// ─── Main Component ───────────────────────────────────────────────────────────

type ActiveTab = 'glossary' | 'standards' | 'maturity';

export const GlossaryPage: React.FC = () => {
  const { i18n } = useTranslation();
  const lang = i18n.language;

  const [search, setSearch]           = useState('');
  const [catFilter, setCatFilter]     = useState<TermCategory | 'ALL'>('ALL');
  const [activeTab, setActiveTab]     = useState<ActiveTab>('glossary');

  const filteredTerms = useMemo(() => {
    let list = [...TERMS];
    if (search.trim()) {
      const q = search.toLowerCase();
      list = list.filter(t =>
        t.term.toLowerCase().includes(q) ||
        t.fullName.toLowerCase().includes(q) ||
        t.definition.toLowerCase().includes(q) ||
        t.masterPrompts.some(mp => mp.toLowerCase().includes(q))
      );
    }
    if (catFilter !== 'ALL') list = list.filter(t => t.category === catFilter);
    return list;
  }, [search, catFilter]);

  const globalAvg = Math.round(MATURITY.reduce((s, m) => s + m.current, 0) / MATURITY.length * 10) / 10;
  const projectedAvg = Math.round(MATURITY.reduce((s, m) => s + m.projected, 0) / MATURITY.length * 10) / 10;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white flex items-center gap-3">
            <BookOpen size={24} className="text-[#1565C0]" />
            {lang === 'ht' ? 'Glosè Teknik SNISID' : 'Glossaire Technique SNISID'}
          </h2>
          <p className="text-slate-400 text-sm mt-0.5">
            {lang === 'ht'
              ? 'v1.0 — MP-013 — Valide pa Konsèy Teknik SNISID'
              : 'v1.0 — MP-013 — Validé par le Conseil Technique SNISID'}
          </p>
        </div>
        <div className="text-right">
          <div className="text-xs text-slate-500">{lang === 'ht' ? 'Tèm ofisyèl' : 'Termes officiels'}</div>
          <div className="text-2xl font-bold text-white">{TERMS.length}</div>
          <div className="text-xs text-slate-600">+ {STANDARDS.length} {lang === 'ht' ? 'estanda' : 'standards'}</div>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 bg-slate-900/60 border border-slate-700/40 rounded-xl p-1">
        {([
          { id: 'glossary' as ActiveTab,  label: lang === 'ht' ? 'Glosè' : 'Glossaire',  icon: BookOpen },
          { id: 'standards' as ActiveTab, label: 'Standards',                              icon: FileCheck },
          { id: 'maturity' as ActiveTab,  label: lang === 'ht' ? 'Matirite' : 'Maturité', icon: Award },
        ] as const).map(tab => {
          const Icon = tab.icon;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex-1 flex items-center justify-center gap-2 py-2.5 rounded-lg text-sm font-medium transition-all focus:outline-none focus:ring-2 focus:ring-[#1565C0] ${
                activeTab === tab.id
                  ? 'bg-[#1565C0] text-white shadow-lg'
                  : 'text-slate-400 hover:text-white hover:bg-slate-800'
              }`}
            >
              <Icon size={15} />
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* ── GLOSSAIRE TAB ── */}
      {activeTab === 'glossary' && (
        <div className="space-y-4">
          {/* Search + Filter */}
          <div className="flex gap-3">
            <div className="relative flex-1">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
              <input
                id="glossary-search"
                type="search"
                placeholder={lang === 'ht' ? 'Chèche tèm, definisyon, MP-XXX...' : 'Rechercher terme, définition, MP-XXX...'}
                value={search}
                onChange={e => setSearch(e.target.value)}
                className="w-full bg-slate-800 border border-slate-700 rounded-lg pl-9 pr-4 py-2.5 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-[#1565C0]"
              />
            </div>
            <select
              value={catFilter}
              onChange={e => setCatFilter(e.target.value as TermCategory | 'ALL')}
              className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-2.5 text-sm text-slate-200 focus:outline-none focus:ring-2 focus:ring-[#1565C0] min-w-[160px]"
            >
              <option value="ALL">{lang === 'ht' ? 'Tout kategori' : 'Toutes catégories'}</option>
              {(Object.keys(CATEGORY_CONFIG) as TermCategory[]).map(cat => (
                <option key={cat} value={cat}>
                  {lang === 'ht' ? CATEGORY_CONFIG[cat].labelHT : CATEGORY_CONFIG[cat].label}
                </option>
              ))}
            </select>
            {(search || catFilter !== 'ALL') && (
              <button
                onClick={() => { setSearch(''); setCatFilter('ALL'); }}
                className="px-3 rounded-lg border border-slate-700 text-slate-400 hover:text-white hover:bg-slate-800 transition-colors focus:outline-none"
              >
                <X size={16} />
              </button>
            )}
          </div>

          {/* Category pills */}
          <div className="flex flex-wrap gap-2">
            {(Object.keys(CATEGORY_CONFIG) as TermCategory[]).map(cat => {
              const cfg = CATEGORY_CONFIG[cat];
              const count = TERMS.filter(t => t.category === cat).length;
              return (
                <button
                  key={cat}
                  onClick={() => setCatFilter(catFilter === cat ? 'ALL' : cat)}
                  className={`flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium border transition-colors focus:outline-none ${
                    catFilter === cat ? cfg.color : 'text-slate-500 border-slate-700 hover:border-slate-600'
                  }`}
                >
                  {lang === 'ht' ? cfg.labelHT : cfg.label}
                  <span className="opacity-60">{count}</span>
                </button>
              );
            })}
          </div>

          {/* Terms list */}
          <div className="space-y-3">
            {filteredTerms.length === 0 ? (
              <div className="py-16 text-center text-slate-500">
                <BookOpen size={40} className="mx-auto mb-4 opacity-30" />
                {lang === 'ht' ? 'Pa gen rezilta' : 'Aucun terme trouvé'}
              </div>
            ) : filteredTerms.map(term => (
              <TermCard key={term.id} term={term} lang={lang} />
            ))}
          </div>

          {/* Footer */}
          <div className="flex items-center justify-between text-xs text-slate-600 pt-2">
            <span>{filteredTerms.length}/{TERMS.length} {lang === 'ht' ? 'tèm' : 'termes'}</span>
            <div className="flex items-center gap-2">
              <Lock size={11} className="text-emerald-400" />
              <span>{lang === 'ht' ? 'v1.0 — SHA-256 verifikasyon aktif' : 'v1.0 — Vérification SHA-256 active'}</span>
            </div>
          </div>
        </div>
      )}

      {/* ── STANDARDS TAB ── */}
      {activeTab === 'standards' && (
        <div className="space-y-3">
          {STANDARDS.map(std => (
            <div key={std.id} className="bg-slate-900/60 border border-slate-700/40 rounded-xl p-4 flex items-start gap-4 hover:border-slate-600/60 transition-colors">
              <div className="p-2 rounded-lg bg-[#1565C0]/10 border border-[#1565C0]/20 flex-shrink-0">
                <FileCheck size={16} className="text-[#1565C0]" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <div className="font-bold text-white text-sm font-mono">{std.code}</div>
                    <div className="text-xs text-slate-400 mt-0.5">{std.domain}</div>
                    <div className="text-xs text-slate-300 mt-1.5">{std.applicability}</div>
                  </div>
                  <span className={`text-xs px-2 py-1 rounded font-medium flex-shrink-0 ${PHASE_COLOR[std.phase] || 'text-slate-400 bg-slate-700/30'}`}>
                    {std.phase}
                  </span>
                </div>
              </div>
            </div>
          ))}
          <div className="text-center text-xs text-slate-600 pt-2">
            {STANDARDS.length} standards applicables — {lang === 'ht' ? 'Validate ak konfòme SNISID' : 'Validés et conformes SNISID'}
          </div>
        </div>
      )}

      {/* ── MATURITÉ TAB ── */}
      {activeTab === 'maturity' && (
        <div className="space-y-6">
          {/* Global score */}
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-[#1E3A5F]/20 border border-slate-700/40 rounded-xl p-5 text-center">
              <div className="text-xs text-slate-400 mb-2">{lang === 'ht' ? 'Nòt Aktyèl' : 'Score Actuel'}</div>
              <div className="text-4xl font-bold text-[#1565C0]">{globalAvg}<span className="text-lg text-slate-500">/100</span></div>
              <div className="mt-2 h-2 bg-slate-800 rounded-full overflow-hidden">
                <div className="h-full bg-[#1565C0] rounded-full" style={{ width: `${globalAvg}%` }} />
              </div>
            </div>
            <div className="bg-emerald-900/10 border border-emerald-700/30 rounded-xl p-5 text-center">
              <div className="text-xs text-slate-400 mb-2">{lang === 'ht' ? 'Pwojeksyon M36' : 'Projection M36'}</div>
              <div className="text-4xl font-bold text-emerald-400">{projectedAvg}<span className="text-lg text-slate-500">/100</span></div>
              <div className="mt-2 text-xs text-emerald-500 font-medium">+{(projectedAvg - globalAvg).toFixed(1)} pts vs actuel</div>
            </div>
          </div>

          {/* Domain bars */}
          <div className="bg-slate-900/60 border border-slate-700/40 rounded-xl p-5 space-y-5">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-3 h-3 rounded bg-[#1565C0]/60" />
              <span className="text-xs text-slate-400">{lang === 'ht' ? 'Aktyèl' : 'Actuel'}</span>
              <div className="w-3 h-3 rounded bg-emerald-500/50 ml-2" />
              <span className="text-xs text-slate-400">{lang === 'ht' ? 'Projekte M36' : 'Projeté M36'}</span>
            </div>
            {MATURITY.map(m => (
              <MaturityBar
                key={m.domain}
                current={m.current}
                projected={m.projected}
                label={m.domain}
                recommendation={m.recommendation}
              />
            ))}
          </div>

          {/* Impact table */}
          <div className="bg-slate-900/60 border border-slate-700/40 rounded-xl overflow-hidden">
            <div className="px-4 py-3 border-b border-slate-700/40">
              <h3 className="text-sm font-medium text-slate-300 flex items-center gap-2">
                <Info size={15} className="text-[#00BCD4]" />
                {lang === 'ht' ? 'Enpak Rekòmandasyon v2.0' : 'Impact Recommandations v2.0'}
              </h3>
            </div>
            <div className="divide-y divide-slate-700/30">
              {[
                { domain: 'Sécurité des Prompts',     rec: 'Anti-Prompt-Injection Hardening',    gain: '+15 pts sécurité' },
                { domain: 'Auditabilité',             rec: 'Versioning + Hash SHA-256 des MPs',  gain: 'Traçabilité totale' },
                { domain: 'Décision Autonome',        rec: 'Arbres IF/THEN/ESCALATE',            gain: '-40% interventions' },
                { domain: 'Biométrie',                rec: 'Détection deepfakes + drift',        gain: 'FAR < 0.001% garanti' },
                { domain: 'SOC',                      rec: 'Playbooks Git + CARICOM Threat',     gain: 'MTTD objectif < 3min' },
                { domain: 'PKI',                      rec: 'CRLs via USSD/SMS (Haïti-first)',    gain: 'Couverture rurale 100%' },
                { domain: 'DevSecOps',                rec: 'SBOM + environnements éphémères',    gain: 'SLSA Level 4' },
                { domain: 'Résilience MEK',           rec: 'Cartographie MEK temps réel',        gain: 'Zéro perte terrain' },
                { domain: 'Gouvernance',              rec: 'Matrice RACI + MoUs int.',           gain: 'Décisions ×60% rapides' },
                { domain: 'UX/Accessibilité',         rec: 'Tests terrain haïtiens réels',       gain: 'Adoption citoyenne > 80%' },
              ].map(item => (
                <div key={item.domain} className="flex items-center gap-4 px-4 py-3 text-xs hover:bg-slate-800/20 transition-colors">
                  <span className="text-slate-400 w-36 flex-shrink-0">{item.domain}</span>
                  <span className="text-slate-300 flex-1 min-w-0 truncate">{item.rec}</span>
                  <span className="text-emerald-400 font-medium flex-shrink-0">{item.gain}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
