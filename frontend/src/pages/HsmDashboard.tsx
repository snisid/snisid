import React from 'react';
import { Key, Shield, Activity, AlertTriangle, RefreshCw } from 'lucide-react';
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

const HsmDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">HSM Key Management</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Hardware Security Module Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
          <span className="text-emerald-500 text-xs font-bold uppercase tracking-wider">HSM Online</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Active Keys" value="1,284" icon={<Key className="w-5 h-5 text-blue-500" />} trend="+24 this week" />
        <KpiCard title="HSM Slots" value="8/12" icon={<Shield className="w-5 h-5 text-emerald-500" />} trend="4 available" />
        <KpiCard title="Key Ops/sec" value="2,450" icon={<Activity className="w-5 h-5 text-violet-500" />} trend="Peak: 3,200" />
        <KpiCard title="Expiring Keys" value="3" icon={<AlertTriangle className="w-5 h-5 text-amber-500" />} trend="Within 30 days" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><RefreshCw className="w-4 h-4 text-blue-500" /> Key Rotation Timeline</h3>
          <div className="space-y-3">
            {[
              { key: 'snisid-root-ca-2026', status: 'Active', rot: '2026-08-15', type: 'RSA-4096' },
              { key: 'hsm-signing-key-a1', status: 'Active', rot: '2026-09-01', type: 'ECDSA-P384' },
              { key: 'enclave-key-prod-3', status: 'Rotating', rot: '2026-07-01', type: 'AES-256' },
              { key: 'backup-hsm-key-v2', status: 'Standby', rot: '2026-12-01', type: 'RSA-2048' },
            ].map((k) => (
              <div key={k.key} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{k.key}</div>
                  <div className="text-xs text-gray-500 mt-0.5">{k.type}</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${k.status === 'Active' ? 'bg-emerald-500/10 text-emerald-500' : k.status === 'Rotating' ? 'bg-amber-500/10 text-amber-500' : 'bg-gray-500/10 text-gray-400'}`}>{k.status}</span>
                  <span className="text-xs text-gray-500">Next: {k.rot}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">HSM Health</h3>
          <div className="space-y-4">
            {[
              { label: 'Session Pool', val: '92%', color: 'bg-emerald-500' },
              { label: 'CPU Load', val: '34%', color: 'bg-blue-500' },
              { label: 'Memory', val: '56%', color: 'bg-violet-500' },
              { label: 'Network I/O', val: '28%', color: 'bg-cyan-500' },
            ].map((m) => (
              <div key={m.label}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-400">{m.label}</span>
                  <span className="text-white font-bold">{m.val}</span>
                </div>
                <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                  <div className={`h-full ${m.color} rounded-full`} style={{ width: m.val }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default HsmDashboard;
