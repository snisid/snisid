import React from 'react';
import { FileSearch, Users, AlertTriangle, Clock, Shield, Activity } from 'lucide-react';
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

const AuditDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Audit Trail</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Immutable Audit Log Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-emerald-500" />
          <span className="text-emerald-500 text-xs font-bold uppercase tracking-wider">Integrity Verified</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Events Today" value="284,192" icon={<Activity className="w-5 h-5 text-blue-500" />} trend="+8% vs yesterday" />
        <KpiCard title="Anomalies Detected" value="23" icon={<AlertTriangle className="w-5 h-5 text-amber-500" />} trend="3 require review" />
        <KpiCard title="Audit Integrity" value="100%" icon={<Shield className="w-5 h-5 text-emerald-500" />} trend="All hashes valid" />
        <KpiCard title="Retention Days" value="2,847" icon={<Clock className="w-5 h-5 text-violet-500" />} trend="7.8 years stored" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><FileSearch className="w-4 h-4 text-blue-500" /> Recent Audit Events</h3>
          <div className="space-y-2">
            {[
              { action: 'Credential Issuance', actor: 'RA-Admin-07', resource: 'ID-2026-08472', status: 'Success' },
              { action: 'Role Escalation', actor: 'SysOp-JC', resource: 'Role:SuperAdmin', status: 'Alert' },
              { action: 'Data Export', actor: 'Auditor-KL', resource: 'Enrollment DB', status: 'Success' },
              { action: 'Failed Login', actor: 'Unknown', resource: 'Portal Admin', status: 'Review' },
            ].map((e) => (
              <div key={e.action} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{e.action}</div>
                  <div className="text-xs text-gray-500">{e.actor} - {e.resource}</div>
                </div>
                <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${e.status === 'Success' ? 'bg-emerald-500/10 text-emerald-500' : e.status === 'Alert' ? 'bg-red-500/10 text-red-500' : 'bg-amber-500/10 text-amber-500'}`}>{e.status}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Audit Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'Merkle Root Integrity', val: 'Validated', color: 'bg-emerald-500' },
              { label: 'Tamper Detection Latency', val: '< 1s', color: 'bg-blue-500' },
              { label: 'Event Indexing Rate', val: '98.5%', color: 'bg-violet-500' },
              { label: 'Compliance Coverage', val: '100%', color: 'bg-cyan-500' },
            ].map((m) => (
              <div key={m.label}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-400">{m.label}</span>
                  <span className="text-white font-bold">{m.val}</span>
                </div>
                <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                  <div className={`h-full ${m.color} rounded-full`} style={{ width: '95%' }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuditDashboard;
