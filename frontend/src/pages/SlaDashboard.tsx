import React from 'react';
import { Activity, CheckCircle, AlertTriangle, Clock, TrendingUp, Shield } from 'lucide-react';
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

const SlaDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">SLA Monitoring</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Service Level Agreement Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-emerald-500" />
          <span className="text-emerald-500 text-xs font-bold uppercase tracking-wider">98.6% Uptime</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Services Monitored" value="24" icon={<Activity className="w-5 h-5 text-blue-500" />} trend="All reporting" />
        <KpiCard title="SLA Met" value="23/24" icon={<CheckCircle className="w-5 h-5 text-emerald-500" />} trend="95.8% compliance" />
        <KpiCard title="SLA Breaches" value="1" icon={<AlertTriangle className="w-5 h-5 text-red-500" />} trend="This month" />
        <KpiCard title="Avg Response Time" value="124ms" icon={<Clock className="w-5 h-5 text-violet-500" />} trend="Target: <200ms" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><TrendingUp className="w-4 h-4 text-blue-500" /> Service Health Status</h3>
          <div className="space-y-2">
            {[
              { service: 'identity-api', sla: '99.95%', uptime: '99.99%', status: 'Healthy' },
              { service: 'fraud-engine', sla: '99.90%', uptime: '99.95%', status: 'Healthy' },
              { service: 'enrollment-svc', sla: '99.50%', uptime: '99.20%', status: 'Warning' },
              { service: 'biometric-svc', sla: '99.95%', uptime: '99.97%', status: 'Healthy' },
            ].map((s) => (
              <div key={s.service} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{s.service}</div>
                  <div className="text-xs text-gray-500">SLA: {s.sla} | Uptime: {s.uptime}</div>
                </div>
                <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${s.status === 'Healthy' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-amber-500/10 text-amber-500'}`}>{s.status}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Monthly SLA Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'Overall Uptime', val: '99.87%', color: 'bg-emerald-500' },
              { label: 'P99 Latency Compliance', val: '96.2%', color: 'bg-blue-500' },
              { label: 'Error Budget Remaining', val: '78%', color: 'bg-violet-500' },
              { label: 'MTTR', val: '14.2m', color: 'bg-cyan-500' },
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

export default SlaDashboard;
