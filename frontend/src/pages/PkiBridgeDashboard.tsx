import React from 'react';
import { Link2, Shield, CheckCircle, AlertTriangle, Activity, Layers } from 'lucide-react';
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

const PkiBridgeDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">PKI Bridge</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Cross-Domain PKI Federation</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-500/10 border border-blue-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
          <span className="text-blue-500 text-xs font-bold uppercase tracking-wider">Connected</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Trust Anchors" value="14" icon={<Shield className="w-5 h-5 text-blue-500" />} trend="2 pending" />
        <KpiCard title="Cross-Signed Certs" value="8,394" icon={<Link2 className="w-5 h-5 text-emerald-500" />} trend="+247 today" />
        <KpiCard title="Bridge Validation" value="99.8%" icon={<CheckCircle className="w-5 h-5 text-violet-500" />} trend="Path success rate" />
        <KpiCard title="Expired Paths" value="3" icon={<AlertTriangle className="w-5 h-5 text-amber-500" />} trend="Requires renewal" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><Layers className="w-4 h-4 text-blue-500" /> Trust Bridges</h3>
          <div className="space-y-2">
            {[
              { bridge: 'US Federal PKI ↔ SNISID', certs: '3,847', status: 'Active', lastSync: '2m ago' },
              { bridge: 'CAN SSP ↔ SNISID', certs: '2,103', status: 'Active', lastSync: '5m ago' },
              { bridge: 'EU eIDAS ↔ SNISID', certs: '1,892', status: 'Active', lastSync: '1m ago' },
              { bridge: 'NATO PKI ↔ SNISID', certs: '552', status: 'Degraded', lastSync: '1h ago' },
            ].map((b) => (
              <div key={b.bridge} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{b.bridge}</div>
                  <div className="text-xs text-gray-500">{b.certs} certificates</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-[10px] text-gray-500">{b.lastSync}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${b.status === 'Active' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-amber-500/10 text-amber-500'}`}>{b.status}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Bridge Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'Cert Path Success Rate', val: '99.8%', color: 'bg-emerald-500' },
              { label: 'CRL Cache Freshness', val: '98%', color: 'bg-blue-500' },
              { label: 'OCSP Responder Uptime', val: '99.99%', color: 'bg-violet-500' },
              { label: 'Cross-Cert Renewal Rate', val: '94%', color: 'bg-cyan-500' },
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

export default PkiBridgeDashboard;
