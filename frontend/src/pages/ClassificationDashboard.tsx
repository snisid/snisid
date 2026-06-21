import React from 'react';
import { motion } from 'framer-motion';
import { FileText, Shield, AlertTriangle, Clock, Lock, Unlock, Eye } from 'lucide-react';

const statCards = [
  { label: 'Classified Resources', value: '476,283', color: 'text-blue-400', icon: FileText },
  { label: 'Audit Events (30d)', value: '12,847', color: 'text-orange-400', icon: Eye },
  { label: 'Active Rules', value: '89', color: 'text-green-400', icon: Shield },
  { label: 'Declass Queue', value: '2,134', color: 'text-yellow-400', icon: Clock },
];

const levels = [
  { level: 'TS//SCI', count: 42300, color: 'text-red-400 bg-red-500/10 border-red-500/20' },
  { level: 'Top Secret', count: 89120, color: 'text-orange-400 bg-orange-500/10 border-orange-500/20' },
  { level: 'Secret', count: 187430, color: 'text-yellow-400 bg-yellow-500/10 border-yellow-500/20' },
  { level: 'Confidential', count: 124810, color: 'text-blue-400 bg-blue-500/10 border-blue-500/20' },
  { level: 'Unclassified', count: 32623, color: 'text-green-400 bg-green-500/10 border-green-500/20' },
];

const auditLog = [
  { id: 'AUD-4821', user: 'Analyst Martinez', action: 'Downgrade', resource: 'NIE 2026-03.pdf', from: 'TS//SCI', to: 'Top Secret', time: '12m ago', severity: 'info' },
  { id: 'AUD-4820', user: 'Admin Chen', action: 'Reclassify', resource: 'Intel Brief #47', from: 'Secret', to: 'Top Secret', time: '28m ago', severity: 'warning' },
  { id: 'AUD-4819', user: 'System', action: 'Auto-Downgrade', resource: 'Report #1122', from: 'Secret', to: 'Confidential', time: '45m ago', severity: 'info' },
  { id: 'AUD-4818', user: 'Analyst Okafor', action: 'Access Violation', resource: 'Database #4', from: '-', to: '-', time: '1h ago', severity: 'error' },
  { id: 'AUD-4817', user: 'Reviewer Singh', action: 'Classification', resource: 'New Product #89', from: '-', to: 'Secret', time: '2h ago', severity: 'info' },
];

const rules = [
  { rule: 'Duration-Based Downgrade', type: 'Automatic', target: 'All TS documents > 10yr', status: 'Active' },
  { rule: 'Originator Control', type: 'Mandatory', target: 'Foreign govt info', status: 'Active' },
  { rule: 'Compilation Rule', type: 'Advisory', target: 'Aggregated unclassified data', status: 'Active' },
  { rule: 'Critical Nuclear Design', type: 'Exemption', target: 'RD/FRD categories', status: 'Active' },
  { rule: 'Automatic Declass Review', type: 'Automatic', target: 'Records > 25yr', status: 'Pending' },
];

const declassQueue = [
  { id: 'DC-001', doc: 'Historical Brief Vol 3', age: '28yr', priority: 'high', analyst: 'Jones' },
  { id: 'DC-002', doc: 'Cold War Assessment', age: '32yr', priority: 'medium', analyst: 'Smith' },
  { id: 'DC-003', doc: 'Treaty Negotiations', age: '26yr', priority: 'high', analyst: 'Garcia' },
  { id: 'DC-004', doc: 'Intelligence Estimates 1998', age: '28yr', priority: 'low', analyst: 'Unassigned' },
  { id: 'DC-005', doc: 'Delegation Authority Memo', age: '30yr', priority: 'medium', analyst: 'Lee' },
];

const ClassificationDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Data Classification</h2>
          <p className="text-gray-500 mt-1">Classification management, auditing, and declassification oversight</p>
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
            <Lock className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Resources by Classification</h3>
          </div>
          <div className="space-y-4">
            {levels.map((l, i) => (
              <motion.div
                key={l.level}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${l.color}`}
              >
                <div className="flex justify-between items-center">
                  <span className="text-sm font-bold">{l.level}</span>
                  <span className="text-lg font-black">{l.count.toLocaleString()}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Eye className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Recent Audit Log</h3>
          </div>
          <div className="space-y-3">
            {auditLog.map((a, i) => (
              <motion.div
                key={a.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-3 rounded-2xl border ${
                  a.severity === 'error' ? 'border-red-500/20 bg-red-500/5' :
                  a.severity === 'warning' ? 'border-yellow-500/20 bg-yellow-500/5' :
                  'border-white/5 bg-white/[0.02]'
                }`}
              >
                <div className="flex justify-between items-center mb-1">
                  <span className="font-mono text-xs text-gray-500">{a.id}</span>
                  <span className="text-xs text-gray-500">{a.time}</span>
                </div>
                <div className="text-sm font-bold">{a.action}</div>
                <div className="text-xs text-gray-400">{a.user} — {a.resource}</div>
                {(a.from && a.from !== '-') && (
                  <div className="text-xs text-gray-500 mt-1">{a.from} → {a.to}</div>
                )}
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Shield className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Active Classification Rules</h3>
          </div>
          <div className="space-y-3">
            {rules.map((r, i) => (
              <motion.div
                key={r.rule}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div className="flex justify-between items-center mb-1">
                  <span className="font-bold text-sm">{r.rule}</span>
                  <span className="text-[10px] font-black uppercase text-green-400">{r.status}</span>
                </div>
                <div className="text-xs text-gray-400">{r.type} — {r.target}</div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
        <div className="flex items-center gap-3 mb-6">
          <Unlock className="w-5 h-5 text-gray-400" />
          <h3 className="text-xl font-bold">Downgrade / Declassification Queue</h3>
        </div>
        <div className="space-y-3">
          {declassQueue.map((d, i) => (
            <motion.div
              key={d.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.05 }}
              className={`p-4 rounded-2xl border flex items-center justify-between ${
                d.priority === 'high' ? 'border-orange-500/20 bg-orange-500/5' :
                d.priority === 'medium' ? 'border-yellow-500/20 bg-yellow-500/5' :
                'border-white/5 bg-white/[0.02]'
              }`}
            >
              <div className="flex items-center gap-4">
                <span className="font-mono text-xs text-gray-500">{d.id}</span>
                <div>
                  <div className="font-bold text-sm">{d.doc}</div>
                  <div className="text-xs text-gray-500">{d.age} · Analyst: {d.analyst}</div>
                </div>
              </div>
              <span className={`text-xs font-bold uppercase ${
                d.priority === 'high' ? 'text-orange-400' :
                d.priority === 'medium' ? 'text-yellow-400' : 'text-green-400'
              }`}>{d.priority}</span>
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

export default ClassificationDashboard;
