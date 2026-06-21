import React from 'react';
import { motion } from 'framer-motion';
import { Radioactive, AlertTriangle, FlaskRound, Clock, Shield, Search, Ban, CheckCircle2 } from 'lucide-react';

const statCards = [
  { label: 'Registered Sources', value: '3,247', color: 'text-blue-400', icon: Radioactive },
  { label: 'Recent Alerts', value: '12', color: 'text-orange-400', icon: AlertTriangle },
  { label: 'Precursors Tracked', value: '891', color: 'text-green-400', icon: FlaskRound },
  { label: 'Unresponded Alerts', value: '2', color: 'text-red-400', icon: Clock },
];

const sourceStatus = [
  { status: 'Registered', count: 2984, color: 'text-green-400 bg-green-500/10 border-green-500/20' },
  { status: 'Lost', count: 87, color: 'text-red-400 bg-red-500/10 border-red-500/20' },
  { status: 'Stolen', count: 34, color: 'text-orange-400 bg-orange-500/10 border-orange-500/20' },
  { status: 'Recovered', count: 142, color: 'text-blue-400 bg-blue-500/10 border-blue-500/20' },
];

const recentAlerts = [
  { id: 'RA-01', location: 'Hospital Central', type: 'Gamma Spike', level: '3', time: '12m ago', status: 'Responding' },
  { id: 'RA-02', location: 'Industrial Zone B', type: 'Source Mismatch', level: '2', time: '45m ago', status: 'Investigating' },
  { id: 'RA-03', location: 'Port Authority', type: 'Container Alert', level: '1', time: '2h ago', status: 'Cleared' },
  { id: 'RA-04', location: 'Research Lab', type: 'Dose Exceeded', level: '2', time: '4h ago', status: 'Cleared' },
  { id: 'RA-05', location: 'Border Crossing', type: 'Vehicle Scan', level: '1', time: '6h ago', status: 'Resolved' },
];

const precursors = [
  { chemical: 'Sodium Hydroxide', facility: 'ChemCorp', status: 'Verified', qty: '2,400 kg' },
  { chemical: 'Sulfuric Acid', facility: 'InduChem', status: 'Pending', qty: '1,800 L' },
  { chemical: 'Hydrochloric Acid', facility: 'PharmaLab', status: 'Verified', qty: '950 L' },
  { chemical: 'Ammonium Nitrate', facility: 'AgriSupply', status: 'Discrepancy', qty: '4,200 kg' },
  { chemical: 'Potassium Permanganate', facility: 'WaterWorks', status: 'Verified', qty: '320 kg' },
];

const unresponded = [
  { location: 'Sector 7 — Waste Facility', type: 'Unaccounted Source', opened: '2026-06-18', level: 'high' },
  { location: 'Medical Depot #3', type: 'Missing Iridium-192', opened: '2026-06-19', level: 'critical' },
];

const RadiationSafetyDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Radiation & Chemical Safety</h2>
          <p className="text-gray-500 mt-1">Nuclear, radiological, and chemical materials tracking and alerting</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {statCards.map((card) => (
          <StatusCard key={card.label} {...card} />
        ))}
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Radioactive className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Source Status</h3>
          </div>
          <div className="grid grid-cols-2 gap-4">
            {sourceStatus.map((s, i) => (
              <motion.div
                key={s.status}
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: i * 0.05 }}
                className={`p-5 rounded-2xl border ${s.color}`}
              >
                <div className="text-3xl font-black mb-1">{s.count}</div>
                <div className="text-sm font-bold">{s.status}</div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <AlertTriangle className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Recent Radiation Alerts</h3>
          </div>
          <div className="space-y-3">
            {recentAlerts.map((a, i) => (
              <motion.div
                key={a.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className="flex items-center justify-between p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div>
                  <div className="flex items-center gap-3 mb-1">
                    <span className="font-mono text-xs text-gray-500">{a.id}</span>
                    <span className={`px-2 py-0.5 rounded text-[10px] font-black ${
                      a.level === '3' ? 'bg-red-500/20 text-red-400' :
                      a.level === '2' ? 'bg-orange-500/20 text-orange-400' : 'bg-yellow-500/20 text-yellow-400'
                    }`}>Level {a.level}</span>
                  </div>
                  <div className="font-bold text-sm">{a.type}</div>
                  <div className="text-xs text-gray-500">{a.location}</div>
                </div>
                <div className="text-right">
                  <div className="text-xs text-gray-500">{a.time}</div>
                  <span className={`text-[10px] font-bold uppercase ${
                    a.status === 'Cleared' || a.status === 'Resolved' ? 'text-green-400' :
                    a.status === 'Responding' ? 'text-orange-400' : 'text-yellow-400'
                  }`}>{a.status}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <FlaskRound className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Chemical Precursor Tracking</h3>
          </div>
          <div className="space-y-3">
            {precursors.map((p, i) => (
              <motion.div
                key={p.chemical}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="flex items-center justify-between p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div>
                  <div className="font-bold text-sm">{p.chemical}</div>
                  <div className="text-xs text-gray-500">{p.facility} · {p.qty}</div>
                </div>
                <span className={`text-xs font-bold uppercase ${
                  p.status === 'Verified' ? 'text-green-400' :
                  p.status === 'Pending' ? 'text-yellow-400' : 'text-red-400'
                }`}>{p.status}</span>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Clock className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Unresponded Alerts</h3>
          </div>
          {unresponded.length === 0 ? (
            <div className="flex items-center gap-3 p-6 rounded-2xl bg-green-500/5 border border-green-500/20">
              <CheckCircle2 className="w-6 h-6 text-green-400" />
              <span className="font-bold text-green-400">No unresponded alerts</span>
            </div>
          ) : (
            <div className="space-y-3">
              {unresponded.map((u, i) => (
                <motion.div
                  key={`${u.location}-${u.type}`}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: i * 0.1 }}
                  className={`p-5 rounded-2xl border ${
                    u.level === 'critical' ? 'border-red-500/20 bg-red-500/5' : 'border-orange-500/20 bg-orange-500/5'
                  }`}
                >
                  <div className="flex justify-between items-center mb-2">
                    <span className={`text-[10px] font-black uppercase tracking-wider ${
                      u.level === 'critical' ? 'text-red-500' : 'text-orange-500'
                    }`}>{u.level}</span>
                  </div>
                  <div className="font-bold mb-1">{u.type}</div>
                  <div className="text-sm text-gray-400">{u.location}</div>
                  <div className="text-xs text-gray-500 mt-2">Opened: {u.opened}</div>
                </motion.div>
              ))}
            </div>
          )}
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

export default RadiationSafetyDashboard;
