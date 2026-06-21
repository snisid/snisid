import React from 'react';
import { motion } from 'framer-motion';
import { Shield, Calendar, AlertTriangle, UserCheck, MapPin, Clock, Star, Users } from 'lucide-react';

const statCards = [
  { label: 'Protectees', value: '18', color: 'text-blue-400', icon: UserCheck },
  { label: 'Upcoming Movements', value: '34', color: 'text-orange-400', icon: Calendar },
  { label: 'Active Threats', value: '6', color: 'text-red-400', icon: AlertTriangle },
  { label: 'Security Clearances', value: '100%', color: 'text-green-400', icon: Shield },
];

const protectees = [
  { name: 'President', level: 'Tier 1', detail: '24/7 Detail', threats: 3, color: 'text-red-400 border-red-500/20 bg-red-500/10' },
  { name: 'Vice President', level: 'Tier 1', detail: '24/7 Detail', threats: 2, color: 'text-red-400 border-red-500/20 bg-red-500/10' },
  { name: 'Senate Leader', level: 'Tier 2', detail: 'Event Security', threats: 1, color: 'text-orange-400 border-orange-500/20 bg-orange-500/10' },
  { name: 'Chief Justice', level: 'Tier 2', detail: 'Residence + Travel', threats: 1, color: 'text-orange-400 border-orange-500/20 bg-orange-500/10' },
  { name: 'Foreign Minister', level: 'Tier 2', detail: 'Travel Security', threats: 2, color: 'text-orange-400 border-orange-500/20 bg-orange-500/10' },
  { name: 'Director NIS', level: 'Tier 3', detail: 'Threat-Based', threats: 0, color: 'text-yellow-400 border-yellow-500/20 bg-yellow-500/10' },
];

const movements = [
  { date: '2026-06-21', protectee: 'President', event: 'State Visit — Embassy', risk: 'high' },
  { date: '2026-06-22', protectee: 'Senate Leader', event: 'Parliament Session', risk: 'medium' },
  { date: '2026-06-23', protectee: 'Foreign Minister', event: 'International Summit', risk: 'critical' },
  { date: '2026-06-24', protectee: 'Vice President', event: 'Regional Tour', risk: 'high' },
  { date: '2026-06-25', protectee: 'Chief Justice', event: 'Judicial Conference', risk: 'low' },
];

const threats = [
  { id: 'TA-01', type: 'Credible Plot', source: 'HUMINT', status: 'Active Investigation', severity: 'critical' },
  { id: 'TA-02', type: 'Social Media', source: 'OSINT', status: 'Monitoring', severity: 'medium' },
  { id: 'TA-03', type: 'Suspicious Surveillance', source: 'SIGINT', status: 'Active Investigation', severity: 'high' },
  { id: 'TA-04', type: 'Letter Threat', source: 'Physical', status: 'Dismissed', severity: 'low' },
];

const securitySummary = [
  { metric: 'Residence Security', status: 'Optimal', color: 'text-green-400' },
  { metric: 'Travel Security', status: 'Enhanced', color: 'text-yellow-400' },
  { metric: 'Comm Systems', status: 'Secure', color: 'text-green-400' },
  { metric: 'Medical Readiness', status: 'Standby', color: 'text-blue-400' },
];

const ExecProtectionDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Executive Protection</h2>
          <p className="text-gray-500 mt-1">Dignitary protection and security operations center</p>
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
            <Star className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Protection Levels</h3>
          </div>
          <div className="space-y-3">
            {protectees.map((p, i) => (
              <motion.div
                key={p.name}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${p.color}`}
              >
                <div className="flex justify-between items-center mb-1">
                  <span className="font-bold">{p.name}</span>
                  <span className="text-xs font-black uppercase tracking-wider">{p.level}</span>
                </div>
                <div className="flex justify-between text-xs text-gray-400">
                  <span>{p.detail}</span>
                  <span>{p.threats} threat{p.threats !== 1 ? 's' : ''}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Calendar className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Upcoming Movements</h3>
          </div>
          <div className="space-y-3">
            {movements.map((m, i) => (
              <motion.div
                key={`${m.date}-${m.protectee}`}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border flex items-center justify-between ${
                  m.risk === 'critical' ? 'border-red-500/20 bg-red-500/5' :
                  m.risk === 'high' ? 'border-orange-500/20 bg-orange-500/5' :
                  m.risk === 'medium' ? 'border-yellow-500/20 bg-yellow-500/5' :
                  'border-white/5 bg-white/[0.02]'
                }`}
              >
                <div>
                  <div className="font-bold text-sm">{m.protectee}</div>
                  <div className="text-xs text-gray-500">{m.event}</div>
                </div>
                <div className="text-right">
                  <div className="text-xs font-mono text-gray-400">{m.date}</div>
                  <span className={`text-[10px] font-black uppercase ${
                    m.risk === 'critical' ? 'text-red-500' :
                    m.risk === 'high' ? 'text-orange-500' :
                    m.risk === 'medium' ? 'text-yellow-500' : 'text-green-500'
                  }`}>{m.risk}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <AlertTriangle className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Active Assessments</h3>
          </div>
          <div className="space-y-3">
            {threats.map((t, i) => (
              <motion.div
                key={t.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div className="flex justify-between items-center mb-1">
                  <span className="font-mono text-xs text-gray-500">{t.id}</span>
                  <span className={`text-[10px] font-black uppercase tracking-wider ${
                    t.severity === 'critical' ? 'text-red-500' :
                    t.severity === 'high' ? 'text-orange-500' :
                    t.severity === 'medium' ? 'text-yellow-500' : 'text-green-500'
                  }`}>{t.severity}</span>
                </div>
                <div className="font-bold text-sm">{t.type}</div>
                <div className="flex justify-between text-xs text-gray-400 mt-1">
                  <span>Source: {t.source}</span>
                  <span>{t.status}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
        <div className="flex items-center gap-3 mb-6">
          <Shield className="w-5 h-5 text-gray-400" />
          <h3 className="text-xl font-bold">Security Status Summary</h3>
        </div>
        <div className="grid grid-cols-4 gap-4">
          {securitySummary.map((s, i) => (
            <motion.div
              key={s.metric}
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: i * 0.05 }}
              className="p-5 rounded-2xl bg-white/[0.02] border border-white/5 text-center"
            >
              <div className="text-sm font-bold mb-2">{s.metric}</div>
              <div className={`${s.color} text-lg font-black`}>{s.status}</div>
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

export default ExecProtectionDashboard;
