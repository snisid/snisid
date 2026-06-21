import React from 'react';
import { motion } from 'framer-motion';
import { Users, FileText, Calendar, AlertTriangle, Star, Shield, Activity, Clock } from 'lucide-react';

const statCards = [
  { label: 'Active Sources', value: '214', color: 'text-blue-400', icon: Users },
  { label: 'Reports (This Month)', value: '892', color: 'text-green-400', icon: FileText },
  { label: 'Debriefings Scheduled', value: '37', color: 'text-orange-400', icon: Calendar },
  { label: 'Source Alerts', value: '8', color: 'text-red-400', icon: AlertTriangle },
];

const credibilityMatrix = [
  { rating: 'A', label: 'Fully Reliable', count: 12, color: 'text-green-400 border-green-500/30 bg-green-500/10' },
  { rating: 'B', label: 'Usually Reliable', count: 34, color: 'text-blue-400 border-blue-500/30 bg-blue-500/10' },
  { rating: 'C', label: 'Fairly Reliable', count: 58, color: 'text-cyan-400 border-cyan-500/30 bg-cyan-500/10' },
  { rating: 'D', label: 'Not Usually Reliable', count: 47, color: 'text-yellow-400 border-yellow-500/30 bg-yellow-500/10' },
  { rating: 'E', label: 'Unreliable', count: 18, color: 'text-orange-400 border-orange-500/30 bg-orange-500/10' },
  { rating: 'F', label: 'Cannot Be Judged', count: 45, color: 'text-red-400 border-red-500/30 bg-red-500/10' },
];

const reportsByMonth = [
  { month: 'Jan', count: 620 }, { month: 'Feb', count: 780 }, { month: 'Mar', count: 650 },
  { month: 'Apr', count: 890 }, { month: 'May', count: 720 }, { month: 'Jun', count: 892 },
];

const maxReport = Math.max(...reportsByMonth.map(r => r.count));

const debriefings = [
  { source: 'CRESTVIEW', handler: 'Agent Martinez', time: '09:00', risk: 'high' },
  { source: 'GOLDENEYE', handler: 'Agent Chen', time: '10:30', risk: 'medium' },
  { source: 'NIGHTHAWK', handler: 'Agent Singh', time: '13:00', risk: 'low' },
  { source: 'SHADOWFAX', handler: 'Agent Okafor', time: '15:30', risk: 'critical' },
  { source: 'WATCHMAN', handler: 'Agent Dubois', time: '16:45', risk: 'medium' },
];

const HumintDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">HUMINT Source Management</h2>
          <p className="text-gray-500 mt-1">Human intelligence source oversight and operations dashboard</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {statCards.map((card) => (
          <StatusCard key={card.label} {...card} />
        ))}
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-8">
            <Star className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Source Credibility Matrix</h3>
          </div>
          <div className="grid grid-cols-3 gap-4">
            {credibilityMatrix.map((c) => (
              <div key={c.rating} className={`p-4 rounded-2xl border ${c.color}`}>
                <div className="text-2xl font-black mb-1">{c.rating}</div>
                <div className="text-xs text-gray-400 uppercase tracking-wider">{c.label}</div>
                <div className="text-lg font-bold mt-2">{c.count}</div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-8">
            <Activity className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Reports Submitted</h3>
          </div>
          <div className="flex items-end gap-3 h-48">
            {reportsByMonth.map((r, i) => (
              <div key={r.month} className="flex-1 flex flex-col items-center gap-2">
                <motion.div
                  initial={{ height: 0 }}
                  animate={{ height: `${(r.count / maxReport) * 100}%` }}
                  transition={{ delay: i * 0.1, duration: 0.6 }}
                  className="w-full bg-blue-500/20 rounded-t-lg relative"
                  style={{ height: `${(r.count / maxReport) * 180}px` }}
                >
                  <div className="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-bold text-blue-400">{r.count}</div>
                </motion.div>
                <span className="text-xs text-gray-500">{r.month}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
        <div className="flex items-center gap-3 mb-6">
          <Calendar className="w-5 h-5 text-gray-400" />
          <h3 className="text-xl font-bold">Debriefing Schedule</h3>
        </div>
        <div className="space-y-3">
          {debriefings.map((d, i) => (
            <motion.div
              key={d.source}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.05 }}
              className="flex items-center justify-between p-4 rounded-2xl bg-white/[0.02] border border-white/5"
            >
              <div className="flex items-center gap-4">
                <div className={`w-2 h-2 rounded-full ${
                  d.risk === 'critical' ? 'bg-red-500' :
                  d.risk === 'high' ? 'bg-orange-500' :
                  d.risk === 'medium' ? 'bg-yellow-500' : 'bg-green-500'
                }`} />
                <span className="font-mono text-sm font-bold">{d.source}</span>
                <span className="text-sm text-gray-400">{d.handler}</span>
              </div>
              <div className="flex items-center gap-3">
                <Clock className="w-4 h-4 text-gray-500" />
                <span className="text-sm font-mono">{d.time}</span>
              </div>
            </motion.div>
          ))}
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

export default HumintDashboard;
