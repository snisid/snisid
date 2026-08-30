import React from 'react';
import { motion } from 'framer-motion';
import { Shield, Ban, Map, Activity, Users, Plane, Train, Ship, Car } from 'lucide-react';

const statCards = [
  { label: 'Screenings Today', value: '48,293', color: 'text-blue-400', icon: Activity },
  { label: 'No-Fly List', value: '4,872', color: 'text-red-400', icon: Ban },
  { label: 'Referrals Made', value: '89', color: 'text-orange-400', icon: Shield },
  { label: 'Clear Rate', value: '98.6%', color: 'text-green-400', icon: Users },
];

const screenings = [
  { type: 'Pass', count: 47612, pct: 98.6, color: 'text-green-400 bg-green-500/10 border-green-500/20' },
  { type: 'Fail', count: 592, pct: 1.2, color: 'text-red-400 bg-red-500/10 border-red-500/20' },
  { type: 'Referral', count: 89, pct: 0.2, color: 'text-orange-400 bg-orange-500/10 border-orange-500/20' },
];

const noFlyStats = [
  { category: 'Domestic', count: 3120, flag: 'red' },
  { category: 'International', count: 1428, flag: 'orange' },
  { category: 'Terrorist Watch', count: 324, flag: 'red' },
];

const zoneSecurity = [
  { zone: 'Terminal A', status: 'Secure', score: 96, color: 'text-green-400' },
  { zone: 'Terminal B', status: 'Secure', score: 93, color: 'text-green-400' },
  { zone: 'Cargo Area', status: 'Elevated', score: 78, color: 'text-yellow-400' },
  { zone: 'Airside', status: 'Secure', score: 97, color: 'text-green-400' },
  { zone: 'Perimeter', status: 'Patrol', score: 85, color: 'text-blue-400' },
];

const volumeTrend = [
  { hour: '00:00', count: 420 }, { hour: '04:00', count: 180 },
  { hour: '06:00', count: 2840 }, { hour: '08:00', count: 7200 },
  { hour: '10:00', count: 5800 }, { hour: '12:00', count: 6300 },
  { hour: '14:00', count: 5400 }, { hour: '16:00', count: 4800 },
  { hour: '18:00', count: 7200 }, { hour: '20:00', count: 4500 },
  { hour: '22:00', count: 1653 },
];

const maxVolume = Math.max(...volumeTrend.map(v => v.count));

const TransportSecurityDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Transportation Security</h2>
          <p className="text-gray-500 mt-1">Aviation and transit security screening operations</p>
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
            <Shield className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Recent Screenings</h3>
          </div>
          <div className="space-y-4">
            {screenings.map((s, i) => (
              <motion.div
                key={s.type}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.1 }}
                className={`p-5 rounded-2xl border ${s.color}`}
              >
                <div className="flex justify-between items-center mb-1">
                  <span className="font-bold">{s.type}</span>
                  <span className="text-xs font-black">{s.pct}%</span>
                </div>
                <div className="text-2xl font-black">{s.count.toLocaleString()}</div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Ban className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">No-Fly List Stats</h3>
          </div>
          <div className="space-y-4">
            {noFlyStats.map((n, i) => (
              <motion.div
                key={n.category}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className={`p-4 rounded-2xl border flex justify-between items-center ${
                  n.flag === 'red' ? 'border-red-500/20 bg-red-500/5' : 'border-orange-500/20 bg-orange-500/5'
                }`}
              >
                <span className="text-sm font-bold">{n.category}</span>
                <span className="text-lg font-black">{n.count.toLocaleString()}</span>
              </motion.div>
            ))}
            <div className="p-4 rounded-2xl bg-white/[0.02] border border-white/5">
              <div className="text-xs text-gray-500 mb-1">Total Unique Identities</div>
              <div className="text-2xl font-black">4,872</div>
            </div>
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Map className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Airport Zone Security</h3>
          </div>
          <div className="space-y-3">
            {zoneSecurity.map((z, i) => (
              <motion.div
                key={z.zone}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="flex items-center justify-between p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <span className="text-sm font-bold">{z.zone}</span>
                <div className="flex items-center gap-3">
                  <div className="h-2 w-24 bg-white/5 rounded-full overflow-hidden">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${z.score}%` }}
                      transition={{ duration: 0.6, delay: i * 0.05 }}
                      className={`h-full rounded-full ${
                        z.score >= 90 ? 'bg-green-500' :
                        z.score >= 80 ? 'bg-yellow-500' : 'bg-blue-500'
                      }`}
                    />
                  </div>
                  <span className={`text-xs font-bold ${z.color}`}>{z.status}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
        <div className="flex items-center gap-3 mb-6">
          <Activity className="w-5 h-5 text-gray-400" />
          <h3 className="text-xl font-bold">Screening Volume Trend (24h)</h3>
        </div>
        <div className="flex items-end gap-2 h-48">
          {volumeTrend.map((v, i) => (
            <div key={v.hour} className="flex-1 flex flex-col items-center gap-2">
              <motion.div
                initial={{ height: 0 }}
                animate={{ height: `${(v.count / maxVolume) * 180}px` }}
                transition={{ delay: i * 0.02, duration: 0.4 }}
                className="w-full bg-blue-500/20 rounded-t-lg relative"
              >
                <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] font-bold text-blue-400">
                  {v.count >= 1000 ? `${(v.count / 1000).toFixed(1)}k` : v.count}
                </div>
              </motion.div>
              <span className="text-[10px] text-gray-500 -rotate-45 origin-left">{v.hour}</span>
            </div>
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

export default TransportSecurityDashboard;
