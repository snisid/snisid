import React from 'react';
import { motion } from 'framer-motion';
import { Radio, Target, FileText, AlertTriangle, Shield, Activity, Clock, Zap } from 'lucide-react';

const statCards = [
  { label: 'Active Targets', value: '347', color: 'text-red-400', icon: Target },
  { label: 'Comms Intercepted (24h)', value: '12,843', color: 'text-blue-400', icon: Radio },
  { label: 'FISA Warrants Active', value: '89', color: 'text-orange-400', icon: FileText },
  { label: 'Emergency Auths', value: '3', color: 'text-yellow-400', icon: Zap },
];

const targetTypes = [
  { type: 'Diplomatic', count: 98, color: 'bg-red-500' },
  { type: 'Military', count: 124, color: 'bg-blue-500' },
  { type: 'Terrorist', count: 67, color: 'bg-orange-500' },
  { type: 'Economic', count: 42, color: 'bg-green-500' },
  { type: 'Cyber', count: 16, color: 'bg-purple-500' },
];

const recentComms = [
  { id: 'C-001', source: 'Embassy #12', type: 'Encrypted Voice', status: 'Decrypted', time: '2m ago', priority: 'high' },
  { id: 'C-002', source: 'Military Attaché', type: 'Satellite Uplink', status: 'Pending', time: '7m ago', priority: 'medium' },
  { id: 'C-003', source: 'Unknown Node', type: 'Data Burst', status: 'Partial', time: '15m ago', priority: 'high' },
  { id: 'C-004', source: 'Trade Mission', type: 'Email Traffic', status: 'Decrypted', time: '23m ago', priority: 'low' },
  { id: 'C-005', source: 'Consulate #4', type: 'Video Conference', status: 'Intercepted', time: '31m ago', priority: 'medium' },
];

const maxTargetCount = Math.max(...targetTypes.map(t => t.count));

const SigintDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">SIGINT Operations Center</h2>
          <p className="text-gray-500 mt-1">Signals intelligence collection and analysis dashboard</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {statCards.map((card) => (
          <StatusCard key={card.label} {...card} />
        ))}
      </div>

      <div className="grid grid-cols-3 gap-6">
        <div className="col-span-1 bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-8">
            <Target className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Target Type Breakdown</h3>
          </div>
          <div className="space-y-5">
            {targetTypes.map((t) => (
              <div key={t.type}>
                <div className="flex justify-between text-sm mb-2">
                  <span className="text-gray-400">{t.type}</span>
                  <span className="font-bold">{t.count}</span>
                </div>
                <div className="h-2 bg-white/5 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${(t.count / maxTargetCount) * 100}%` }}
                    transition={{ duration: 0.8, ease: 'easeOut' }}
                    className={`h-full rounded-full ${t.color}`}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="col-span-2 bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Activity className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Recent Intercepted Communications</h3>
          </div>
          <table className="w-full text-left">
            <thead>
              <tr className="text-xs uppercase tracking-widest text-gray-500 border-b border-white/5">
                <th className="pb-4 font-medium">ID</th>
                <th className="pb-4 font-medium">Source</th>
                <th className="pb-4 font-medium">Type</th>
                <th className="pb-4 font-medium">Status</th>
                <th className="pb-4 font-medium">Time</th>
              </tr>
            </thead>
            <tbody>
              {recentComms.map((c, i) => (
                <motion.tr
                  key={c.id}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: i * 0.05 }}
                  className="border-b border-white/5 hover:bg-white/[0.02] transition-colors"
                >
                  <td className="py-4 font-mono text-sm">{c.id}</td>
                  <td className="py-4 text-sm">{c.source}</td>
                  <td className="py-4 text-sm text-gray-300">{c.type}</td>
                  <td className="py-4">
                    <span className={`text-xs font-bold uppercase tracking-wider ${
                      c.status === 'Decrypted' ? 'text-green-400' :
                      c.status === 'Intercepted' ? 'text-blue-400' :
                      c.status === 'Partial' ? 'text-orange-400' : 'text-yellow-400'
                    }`}>{c.status}</span>
                  </td>
                  <td className="py-4 text-sm text-gray-500">{c.time}</td>
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="bg-[#0f1218] rounded-[2rem] border border-yellow-500/20 p-6 flex items-center gap-5">
        <div className="p-3 bg-yellow-500/10 rounded-2xl">
          <AlertTriangle className="w-6 h-6 text-yellow-500" />
        </div>
        <div>
          <h4 className="font-bold">Emergency Authorization Active</h4>
          <p className="text-sm text-gray-400 mt-1">FISA emergency authorization #EA-2026-089 active — retroactive filing due in 72 hours</p>
        </div>
        <div className="ml-auto flex items-center gap-2 text-yellow-500">
          <Clock className="w-4 h-4" />
          <span className="text-sm font-bold font-mono">67:42:11</span>
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

export default SigintDashboard;
