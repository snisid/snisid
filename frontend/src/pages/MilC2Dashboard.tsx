import React from 'react';
import { motion } from 'framer-motion';
import { Map, Shield, FileText, Clock, Users, Activity, Target, Radio } from 'lucide-react';

const statCards = [
  { label: 'Active Operations', value: '7', color: 'text-red-400', icon: Map },
  { label: 'Units Deployed', value: '12,480', color: 'text-blue-400', icon: Shield },
  { label: 'SITREPs Today', value: '43', color: 'text-green-400', icon: FileText },
  { label: 'Mission Clock', value: 'D+127', color: 'text-orange-400', icon: Clock },
];

const operations = [
  { name: 'OP SENTINEL', status: 'Active', location: 'North Sector', personnel: '3,200', priority: 'critical' },
  { name: 'OP GUARDIAN', status: 'Active', location: 'Maritime Zone', personnel: '2,100', priority: 'high' },
  { name: 'OP WATCHTOWER', status: 'Standby', location: 'Eastern Border', personnel: '1,500', priority: 'medium' },
  { name: 'OP PHOENIX', status: 'Active', location: 'Urban Center', personnel: '850', priority: 'high' },
  { name: 'OP IRONCLAD', status: 'Reserve', location: 'HQ Battalion', personnel: '4,830', priority: 'low' },
];

const deployments = [
  { label: 'Deployed', value: 72, color: 'text-blue-400 bg-blue-500/10 border-blue-500/20' },
  { label: 'Standby', value: 18, color: 'text-yellow-400 bg-yellow-500/10 border-yellow-500/20' },
  { label: 'Reserve', value: 10, color: 'text-gray-400 bg-gray-500/10 border-gray-500/20' },
];

const sitreps = [
  { id: 'SR-047', op: 'SENTINEL', message: 'Patrol encountered small arms fire. No casualties. Area secured.', time: '14:32', priority: 'high' },
  { id: 'SR-046', op: 'GUARDIAN', message: 'Maritime interdiction — vessel boarded, 12 suspects detained.', time: '13:15', priority: 'high' },
  { id: 'SR-045', op: 'WATCHTOWER', message: 'Drone surveillance pattern completed. No anomalous activity detected.', time: '11:48', priority: 'low' },
  { id: 'SR-044', op: 'SENTINEL', message: 'Supply convoy arrived at FOB Alpha. Resupply complete.', time: '10:02', priority: 'medium' },
  { id: 'SR-043', op: 'PHOENIX', message: 'IED discovered and neutralized by EOD team.', time: '08:30', priority: 'critical' },
];

const timeline = [
  { time: 'D+0', event: 'OP SENTINEL — Initial deployment' },
  { time: 'D+14', event: 'OP GUARDIAN — Maritime zone established' },
  { time: 'D+45', event: 'OP WATCHTOWER — Border surveillance active' },
  { time: 'D+89', event: 'OP PHOENIX — Urban reinforcement' },
  { time: 'D+127', event: 'Current — All operations nominal' },
];

const MilC2Dashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Military Command & Control</h2>
          <p className="text-gray-500 mt-1">Operational command center — force disposition and tactical reporting</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {statCards.map((card) => (
          <StatusCard key={card.label} {...card} />
        ))}
      </div>

      <div className="grid grid-cols-3 gap-6">
        <div className="col-span-2 bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Map className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Active Operations</h3>
          </div>
          <div className="space-y-3">
            {operations.map((op, i) => (
              <motion.div
                key={op.name}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-5 rounded-2xl border flex items-center justify-between ${
                  op.priority === 'critical' ? 'border-red-500/20 bg-red-500/5' :
                  op.priority === 'high' ? 'border-orange-500/20 bg-orange-500/5' :
                  op.priority === 'medium' ? 'border-yellow-500/20 bg-yellow-500/5' :
                  'border-white/5 bg-white/[0.02]'
                }`}
              >
                <div>
                  <div className="flex items-center gap-3 mb-1">
                    <span className="font-bold">{op.name}</span>
                    <span className={`text-[10px] font-black uppercase tracking-widest ${
                      op.status === 'Active' ? 'text-green-400' :
                      op.status === 'Standby' ? 'text-yellow-400' : 'text-gray-500'
                    }`}>{op.status}</span>
                  </div>
                  <div className="text-sm text-gray-500">{op.location} · {op.personnel} personnel</div>
                </div>
                <Target className="w-5 h-5 text-gray-600" />
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-8">
            <Shield className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Deployment Status</h3>
          </div>
          <div className="flex flex-col items-center mb-8">
            <div className="relative w-40 h-40 mb-4">
              <svg viewBox="0 0 100 100" className="w-full h-full -rotate-90">
                {deployments.reduce((acc, d, i) => {
                  const prev = acc.length > 0 ? acc[acc.length - 1].endAngle : 0;
                  const angle = (d.value / 100) * 360;
                  const end = prev + angle;
                  const x1 = 50 + 38 * Math.cos((prev * Math.PI) / 180);
                  const y1 = 50 + 38 * Math.sin((prev * Math.PI) / 180);
                  const x2 = 50 + 38 * Math.cos((end * Math.PI) / 180);
                  const y2 = 50 + 38 * Math.sin((end * Math.PI) / 180);
                  acc.push({ d: `M 50 50 L ${x1} ${y1} A 38 38 0 ${angle > 180 ? 1 : 0} 1 ${x2} ${y2} Z`, color: d.label === 'Deployed' ? '#3b82f6' : d.label === 'Standby' ? '#eab308' : '#6b7280', endAngle: end });
                  return acc;
                }, [] as { d: string; color: string; endAngle: number }[]).map((s, i) => (
                  <motion.path
                    key={i}
                    d={s.d}
                    fill={s.color}
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 0.7 }}
                    transition={{ delay: i * 0.1 }}
                  />
                ))}
              </svg>
            </div>
            <div className="space-y-2 w-full">
              {deployments.map((d) => (
                <div key={d.label} className={`p-3 rounded-xl border text-center ${d.color}`}>
                  <div className="text-xs uppercase tracking-wider">{d.label}</div>
                  <div className="text-xl font-black">{d.value}%</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <FileText className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">SITREP Feed</h3>
          </div>
          <div className="space-y-4">
            {sitreps.map((s, i) => (
              <motion.div
                key={s.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${
                  s.priority === 'critical' ? 'border-red-500/20 bg-red-500/5' :
                  s.priority === 'high' ? 'border-orange-500/20 bg-orange-500/5' :
                  'border-white/5 bg-white/[0.02]'
                }`}
              >
                <div className="flex items-center gap-3 mb-2">
                  <span className="text-xs font-mono text-gray-500">{s.id}</span>
                  <span className="text-[10px] font-black text-blue-400 uppercase tracking-wider">{s.op}</span>
                  <span className="text-xs text-gray-500 ml-auto">{s.time}</span>
                </div>
                <p className="text-sm">{s.message}</p>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Activity className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Operation Timeline</h3>
          </div>
          <div className="relative pl-8 border-l border-white/10 space-y-6">
            {timeline.map((t, i) => (
              <motion.div
                key={t.time}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.1 }}
                className="relative"
              >
                <div className="absolute -left-[2.35rem] top-1 w-3 h-3 rounded-full bg-blue-500 border-2 border-[#0f1218]" />
                <div className="text-xs font-mono text-blue-400 mb-1">{t.time}</div>
                <div className="text-sm text-gray-300">{t.event}</div>
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

export default MilC2Dashboard;
