import React from 'react';
import { Server, Shield, Clock, AlertTriangle, RefreshCw, Database } from 'lucide-react';
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

const DrDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Disaster Recovery</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">DR Orchestration & Failover Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-emerald-500" />
          <span className="text-emerald-500 text-xs font-bold uppercase tracking-wider">All Sites Healthy</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="RPO" value="4.2s" icon={<Clock className="w-5 h-5 text-blue-500" />} trend="Target: 60s" />
        <KpiCard title="RTO" value="2.1m" icon={<Server className="w-5 h-5 text-emerald-500" />} trend="Target: 15m" />
        <KpiCard title="Failovers This Year" value="3" icon={<RefreshCw className="w-5 h-5 text-amber-500" />} trend="All successful" />
        <KpiCard title="Active DR Sites" value="3" icon={<Shield className="w-5 h-5 text-violet-500" />} trend="Primary + 2 replicas" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {[
          { site: 'Primary (US-East-1)', status: 'Active', cpu: '42%', mem: '58%', lag: '0ms' },
          { site: 'DR Replica (US-West-2)', status: 'Standby', cpu: '12%', mem: '24%', lag: '1.2s' },
          { site: 'DR Replica (EU-West-1)', status: 'Standby', cpu: '8%', mem: '18%', lag: '3.7s' },
        ].map((s) => (
          <div key={s.site} className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-bold text-white">{s.site}</h3>
              <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${s.status === 'Active' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-gray-500/10 text-gray-400'}`}>{s.status}</span>
            </div>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between"><span className="text-gray-400">CPU</span><span className="text-white">{s.cpu}</span></div>
              <div className="flex justify-between"><span className="text-gray-400">Memory</span><span className="text-white">{s.mem}</span></div>
              <div className="flex justify-between"><span className="text-gray-400">Replication Lag</span><span className="text-white">{s.lag}</span></div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default DrDashboard;
