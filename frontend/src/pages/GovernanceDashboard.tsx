import React from 'react';
import { Scale, CheckCircle, AlertTriangle, FileText, Users, Activity } from 'lucide-react';
import { motion } from 'framer-motion';

const KpiCard: React.FC<{ title: string; value: string; icon: React.ReactNode; trend: string }> = ({ title, value, icon, trend }) => (
  <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl flex flex-col gap-4">
    <div className="flex items-start justify-between">
      <h3 className="text-gray-400 font-medium text-sm uppercase tracking-wider">{title}</h3>
      <div className="p-2 bg-white/5 rounded-xl">{icon}</div>
    </div>
    <div>
      <div className="text-3xl font-bold text-white mb-1">{value}</div>
      <div className="text-xs text-gray-500">{trend}</div>
    </div>
  </motion.div>
);

const GovernanceDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Governance</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Policy & Compliance Oversight Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-500/10 border border-blue-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
          <span className="text-blue-500 text-xs font-bold uppercase tracking-wider">All Policies Active</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Active Policies" value="84" icon={<Scale className="w-5 h-5 text-blue-500" />} trend="+3 this quarter" />
        <KpiCard title="Compliance Score" value="94.2%" icon={<CheckCircle className="w-5 h-5 text-emerald-500" />} trend="+1.8% QoQ" />
        <KpiCard title="Open Controls" value="12" icon={<AlertTriangle className="w-5 h-5 text-amber-500" />} trend="2 overdue" />
        <KpiCard title="Audits Scheduled" value="6" icon={<FileText className="w-5 h-5 text-violet-500" />} trend="This quarter" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><Activity className="w-4 h-4 text-blue-500" /> Policy Compliance</h3>
          <div className="space-y-2">
            {[
              { policy: 'Data Privacy (GDPR)', owner: 'DPO Office', status: 'Compliant', score: '98%' },
              { policy: 'Access Control (NIST 800-53)', owner: 'Security', status: 'Compliant', score: '96%' },
              { policy: 'Records Retention', owner: 'Legal', status: 'In Progress', score: '72%' },
              { policy: 'Third-Party Risk', owner: 'Procurement', status: 'Compliant', score: '91%' },
              { policy: 'Business Continuity', owner: 'Operations', status: 'Compliant', score: '95%' },
            ].map((p) => (
              <div key={p.policy} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{p.policy}</div>
                  <div className="text-xs text-gray-500">Owner: {p.owner}</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-gray-500">{p.score}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${p.status === 'Compliant' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-amber-500/10 text-amber-500'}`}>{p.status}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Governance Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'Policy Attestation Rate', val: '94%', color: 'bg-blue-500' },
              { label: 'Training Completion', val: '88%', color: 'bg-violet-500' },
              { label: 'Risk Acceptances', val: '12 active', color: 'bg-amber-500' },
              { label: 'Control Effectiveness', val: '96%', color: 'bg-emerald-500' },
            ].map((m) => (
              <div key={m.label}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-400">{m.label}</span>
                  <span className="text-white font-bold">{m.val}</span>
                </div>
                <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                  <div className={`h-full ${m.color} rounded-full`} style={{ width: '90%' }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default GovernanceDashboard;
