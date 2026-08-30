import React from 'react';
import { motion } from 'framer-motion';
import { FileText, CheckCircle2, XCircle, Clock, AlertTriangle, Scale, Shield, Activity } from 'lucide-react';

const statCards = [
  { label: 'Warrants Filed', value: '234', color: 'text-blue-400', icon: FileText },
  { label: 'Approved', value: '218', color: 'text-green-400', icon: CheckCircle2 },
  { label: 'Denied', value: '16', color: 'text-red-400', icon: XCircle },
  { label: 'Emergency Auths', value: '7', color: 'text-orange-400', icon: AlertTriangle },
];

const warrantTimeline = [
  { month: 'Jan', filed: 38, approved: 35, denied: 3 },
  { month: 'Feb', filed: 32, approved: 30, denied: 2 },
  { month: 'Mar', filed: 42, approved: 39, denied: 3 },
  { month: 'Apr', filed: 36, approved: 34, denied: 2 },
  { month: 'May', filed: 44, approved: 41, denied: 3 },
  { month: 'Jun', filed: 42, approved: 39, denied: 3 },
];

const maxWarrant = Math.max(...warrantTimeline.flatMap(w => [w.filed, w.approved, w.denied]));

const emergencyAuths = [
  { id: 'EA-026', target: 'Foreign Entity X', date: '2026-06-18', status: 'Retroactive Filing Due', severity: 'high' },
  { id: 'EA-025', target: 'Communications Node', date: '2026-06-15', status: 'Filed — Pending Review', severity: 'medium' },
  { id: 'EA-024', target: 'Cyber Threat Actor', date: '2026-06-12', status: 'Approved Retroactively', severity: 'low' },
  { id: 'EA-023', target: 'Diplomatic Cover', date: '2026-06-10', status: 'Approved Retroactively', severity: 'low' },
  { id: 'EA-022', target: 'Financial Network', date: '2026-06-08', status: 'Filed — Pending Review', severity: 'medium' },
];

const activeWarrantsByType = [
  { type: 'National Security', count: 62, pct: 71, color: 'bg-red-500' },
  { type: 'Foreign Intelligence', count: 18, pct: 21, color: 'bg-blue-500' },
  { type: 'Counterintelligence', count: 7, pct: 8, color: 'bg-yellow-500' },
];

const complianceReports = [
  { report: 'Q1 2026 Compliance Audit', status: 'Certified', issues: 0, color: 'text-green-400' },
  { report: 'Monthly Minimization Report', status: 'Pending Review', issues: 3, color: 'text-yellow-400' },
  { report: 'Incidental Collection Review', status: 'Filed', issues: 1, color: 'text-blue-400' },
  { report: 'Annual FISA Report', status: 'In Progress', issues: 0, color: 'text-gray-400' },
];

const FisaCourtDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">FISA Court Oversight</h2>
          <p className="text-gray-500 mt-1">Foreign Intelligence Surveillance Court warrant management and compliance</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {statCards.map((card) => (
          <StatusCard key={card.label} {...card} />
        ))}
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Activity className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Warrants Filed / Approved / Denied</h3>
          </div>
          <div className="flex items-end gap-3 h-64">
            {warrantTimeline.map((w, i) => (
              <div key={w.month} className="flex-1 flex flex-col items-center gap-1">
                <motion.div
                  initial={{ height: 0 }}
                  animate={{ height: `${(w.approved / maxWarrant) * 200}px` }}
                  transition={{ delay: i * 0.05, duration: 0.5 }}
                  className="w-full bg-green-500/20 rounded-t-lg relative"
                >
                  <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] font-bold text-green-400">{w.approved}</div>
                </motion.div>
                <motion.div
                  initial={{ height: 0 }}
                  animate={{ height: `${(w.filed / maxWarrant) * 200}px` }}
                  transition={{ delay: i * 0.05, duration: 0.5 }}
                  className="w-full bg-blue-500/20 rounded-t-lg relative"
                >
                  <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] font-bold text-blue-400">{w.filed}</div>
                </motion.div>
                <motion.div
                  initial={{ height: 0 }}
                  animate={{ height: `${(w.denied / maxWarrant) * 200}px` }}
                  transition={{ delay: i * 0.05, duration: 0.5 }}
                  className="w-full bg-red-500/20 rounded-t-lg relative"
                >
                  <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] font-bold text-red-400">{w.denied}</div>
                </motion.div>
                <span className="text-xs text-gray-500 mt-1">{w.month}</span>
              </div>
            ))}
          </div>
          <div className="flex justify-center gap-6 mt-4 text-xs">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded bg-blue-500/40" />
              <span className="text-gray-400">Filed</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded bg-green-500/40" />
              <span className="text-gray-400">Approved</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded bg-red-500/40" />
              <span className="text-gray-400">Denied</span>
            </div>
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <AlertTriangle className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Emergency Authorizations</h3>
          </div>
          <div className="space-y-3">
            {emergencyAuths.map((ea, i) => (
              <motion.div
                key={ea.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${
                  ea.severity === 'high' ? 'border-red-500/20 bg-red-500/5' :
                  ea.severity === 'medium' ? 'border-yellow-500/20 bg-yellow-500/5' :
                  'border-green-500/20 bg-green-500/5'
                }`}
              >
                <div className="flex justify-between items-center mb-1">
                  <span className="font-mono text-xs text-gray-500">{ea.id}</span>
                  <span className={`text-[10px] font-black uppercase tracking-wider ${
                    ea.severity === 'high' ? 'text-red-500' :
                    ea.severity === 'medium' ? 'text-yellow-500' : 'text-green-500'
                  }`}>{ea.severity}</span>
                </div>
                <div className="font-bold text-sm">{ea.target}</div>
                <div className="flex justify-between text-xs text-gray-400 mt-1">
                  <span>{ea.date}</span>
                  <span>{ea.status}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Scale className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Active Warrants by Type</h3>
          </div>
          <div className="space-y-5">
            {activeWarrantsByType.map((w, i) => (
              <div key={w.type}>
                <div className="flex justify-between text-sm mb-2">
                  <span className="text-gray-400">{w.type}</span>
                  <div className="flex items-center gap-3">
                    <span className="font-bold">{w.count}</span>
                    <span className="text-xs text-gray-500">{w.pct}%</span>
                  </div>
                </div>
                <div className="h-3 bg-white/5 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${w.pct}%` }}
                    transition={{ duration: 0.6, delay: i * 0.1 }}
                    className={`h-full rounded-full ${w.color}`}
                  />
                </div>
              </div>
            ))}
            <div className="mt-6 p-4 rounded-2xl bg-white/[0.02] border border-white/5">
              <div className="text-xs text-gray-500 mb-1">Total Active Warrants</div>
              <div className="text-3xl font-black">{activeWarrantsByType.reduce((s, w) => s + w.count, 0)}</div>
            </div>
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Shield className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Compliance Reports</h3>
          </div>
          <div className="space-y-3">
            {complianceReports.map((cr, i) => (
              <motion.div
                key={cr.report}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className="flex items-center justify-between p-5 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div>
                  <div className="font-bold text-sm">{cr.report}</div>
                  <div className="text-xs text-gray-500 mt-1">
                    {cr.issues > 0 ? `${cr.issues} issue${cr.issues > 1 ? 's' : ''} found` : 'No issues'}
                  </div>
                </div>
                <span className={`text-xs font-bold uppercase ${cr.color}`}>{cr.status}</span>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

const StatusCard = ({ label, value, color, icon: Icon }: { label: string; value: string; color: string; icon: React.ComponentType<any> }) => (
  <motion.div
    initial={{ opacity: 0, y: 20 }}
    animate={{ opacity: 1, y: 0 }}
    className="p-6 bg-[#0f1218] rounded-3xl border border-white/5"
  >
    <div className="flex justify-between items-start mb-2">
      <p className="text-xs font-bold uppercase tracking-widest text-gray-500">{label}</p>
      <Icon className={`w-5 h-5 ${color}`} />
    </div>
    <h3 className="text-3xl font-black tracking-tighter">{value}</h3>
  </motion.div>
);

export default FisaCourtDashboard;
