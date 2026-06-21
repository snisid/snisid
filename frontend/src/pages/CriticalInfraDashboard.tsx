import React from 'react';
import { motion } from 'framer-motion';
import { Shield, AlertTriangle, Building2, Zap, Droplets, Wifi, Activity, Globe } from 'lucide-react';

const statCards = [
  { label: 'Total Assets', value: '4,832', color: 'text-blue-400', icon: Building2 },
  { label: 'Active Incidents', value: '17', color: 'text-red-400', icon: AlertTriangle },
  { label: 'Sectors Monitored', value: '16', color: 'text-green-400', icon: Globe },
  { label: 'Resilience Score', value: 'B+', color: 'text-yellow-400', icon: Shield },
];

const sectorAssets = [
  { sector: 'Energy', assets: 1240, color: 'bg-yellow-500' },
  { sector: 'Water', assets: 870, color: 'bg-blue-500' },
  { sector: 'Transportation', assets: 980, color: 'bg-green-500' },
  { sector: 'Communications', assets: 650, color: 'bg-purple-500' },
  { sector: 'Healthcare', assets: 540, color: 'bg-red-500' },
  { sector: 'Government', assets: 552, color: 'bg-orange-500' },
];

const maxSectorAssets = Math.max(...sectorAssets.map(s => s.assets));

const incidents = [
  { id: 'INC-01', sector: 'Energy', event: 'Power substation fire', status: 'Contained', severity: 'critical' },
  { id: 'INC-02', sector: 'Water', event: 'Treatment plant breach', status: 'Mitigating', severity: 'high' },
  { id: 'INC-03', sector: 'Transport', event: 'Bridge structural alert', status: 'Monitoring', severity: 'medium' },
  { id: 'INC-04', sector: 'Comms', event: 'Fiber optic cut', status: 'Repairing', severity: 'high' },
  { id: 'INC-05', sector: 'Healthcare', event: 'Generator failure', status: 'Resolved', severity: 'low' },
];

const riskAssessments = [
  { sector: 'Energy', risk: 'High', score: 72, color: 'text-red-400' },
  { sector: 'Water', risk: 'Medium', score: 58, color: 'text-yellow-400' },
  { sector: 'Transport', risk: 'Medium', score: 51, color: 'text-yellow-400' },
  { sector: 'Comms', risk: 'High', score: 68, color: 'text-red-400' },
  { sector: 'Healthcare', risk: 'Low', score: 34, color: 'text-green-400' },
  { sector: 'Government', risk: 'Medium', score: 45, color: 'text-yellow-400' },
];

const CriticalInfraDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Critical Infrastructure Protection</h2>
          <p className="text-gray-500 mt-1">National critical infrastructure security and resilience monitoring</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {statCards.map((card) => (
          <StatusCard key={card.label} {...card} />
        ))}
      </div>

      <div className="grid grid-cols-3 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Building2 className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Assets by Sector</h3>
          </div>
          <div className="space-y-4">
            {sectorAssets.map((s, i) => (
              <div key={s.sector}>
                <div className="flex justify-between text-sm mb-2">
                  <div className="flex items-center gap-2">
                    {s.sector === 'Energy' && <Zap className="w-4 h-4 text-yellow-500" />}
                    {s.sector === 'Water' && <Droplets className="w-4 h-4 text-blue-500" />}
                    {s.sector === 'Transportation' && <Building2 className="w-4 h-4 text-green-500" />}
                    {s.sector === 'Communications' && <Wifi className="w-4 h-4 text-purple-500" />}
                    {s.sector === 'Healthcare' && <Activity className="w-4 h-4 text-red-500" />}
                    {s.sector === 'Government' && <Shield className="w-4 h-4 text-orange-500" />}
                    <span className="text-gray-400">{s.sector}</span>
                  </div>
                  <span className="font-bold">{s.assets.toLocaleString()}</span>
                </div>
                <div className="h-2 bg-white/5 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${(s.assets / maxSectorAssets) * 100}%` }}
                    transition={{ duration: 0.6, delay: i * 0.05 }}
                    className={`h-full rounded-full ${s.color}`}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <AlertTriangle className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Active Incidents</h3>
          </div>
          <div className="space-y-3">
            {incidents.map((inc, i) => (
              <motion.div
                key={inc.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${
                  inc.severity === 'critical' ? 'border-red-500/20 bg-red-500/5' :
                  inc.severity === 'high' ? 'border-orange-500/20 bg-orange-500/5' :
                  inc.severity === 'medium' ? 'border-yellow-500/20 bg-yellow-500/5' :
                  'border-white/5 bg-white/[0.02]'
                }`}
              >
                <div className="flex justify-between items-center mb-2">
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-xs text-gray-500">{inc.id}</span>
                    <span className="text-xs font-bold text-blue-400">{inc.sector}</span>
                  </div>
                  <span className={`text-[10px] font-black uppercase ${
                    inc.status === 'Resolved' ? 'text-green-400' :
                    inc.status === 'Contained' ? 'text-blue-400' :
                    inc.status === 'Mitigating' ? 'text-orange-400' :
                    inc.status === 'Repairing' ? 'text-yellow-400' : 'text-gray-500'
                  }`}>{inc.status}</span>
                </div>
                <div className="text-sm">{inc.event}</div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Activity className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Sector Risk Assessments</h3>
          </div>
          <div className="space-y-3">
            {riskAssessments.map((r, i) => (
              <motion.div
                key={r.sector}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="flex items-center justify-between p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <span className="text-sm font-bold">{r.sector}</span>
                <div className="flex items-center gap-3">
                  <div className="h-2 w-16 bg-white/5 rounded-full overflow-hidden">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${r.score}%` }}
                      transition={{ duration: 0.6, delay: i * 0.05 }}
                      className={`h-full rounded-full ${
                        r.score >= 70 ? 'bg-red-500' :
                        r.score >= 50 ? 'bg-yellow-500' : 'bg-green-500'
                      }`}
                    />
                  </div>
                  <span className={`text-xs font-bold ${r.color}`}>{r.risk}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
        <div className="flex items-center gap-3 mb-6">
          <Globe className="w-5 h-5 text-gray-400" />
          <h3 className="text-xl font-bold">National Resilience Score</h3>
        </div>
        <div className="grid grid-cols-5 gap-6">
          <div className="col-span-1 flex flex-col items-center justify-center">
            <div className="text-6xl font-black text-yellow-400">B+</div>
            <div className="text-sm text-gray-500 mt-2">Overall Score</div>
          </div>
          <div className="col-span-4 grid grid-cols-4 gap-4">
            {[
              { metric: 'Prevention', score: 84, color: 'text-green-400' },
              { metric: 'Detection', score: 72, color: 'text-yellow-400' },
              { metric: 'Response', score: 78, color: 'text-yellow-400' },
              { metric: 'Recovery', score: 65, color: 'text-orange-400' },
            ].map((m, i) => (
              <motion.div
                key={m.metric}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="text-center p-5 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div className={`text-3xl font-black ${m.color}`}>{m.score}%</div>
                <div className="text-xs text-gray-500 mt-1 uppercase tracking-wider">{m.metric}</div>
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

export default CriticalInfraDashboard;
