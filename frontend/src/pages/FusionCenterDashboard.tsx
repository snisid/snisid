import React from 'react';
import { motion } from 'framer-motion';
import { Layers, FileText, Shield, Users, AlertTriangle, BookOpen, GitBranch, Activity } from 'lucide-react';

const statCards = [
  { label: 'Active Products', value: '476', color: 'text-blue-400', icon: FileText },
  { label: 'Threat Actors Tracked', value: '134', color: 'text-red-400', icon: Users },
  { label: 'Cross-Discipline Links', value: '89', color: 'text-green-400', icon: GitBranch },
  { label: 'NIES Published (YTD)', value: '23', color: 'text-purple-400', icon: BookOpen },
];

const classifications = [
  { level: 'TS//SCI', count: 89, color: 'text-red-400 bg-red-500/10 border-red-500/20' },
  { level: 'Top Secret', count: 142, color: 'text-orange-400 bg-orange-500/10 border-orange-500/20' },
  { level: 'Secret', count: 187, color: 'text-yellow-400 bg-yellow-500/10 border-yellow-500/20' },
  { level: 'Confidential', count: 58, color: 'text-blue-400 bg-blue-500/10 border-blue-500/20' },
];

const threatHeatMap = [
  { actor: 'Cyber Group Alpha', capability: 85, intent: 90, risk: 'critical' },
  { actor: 'Foreign Intel Service', capability: 95, intent: 75, risk: 'high' },
  { actor: 'Insider Network', capability: 60, intent: 85, risk: 'high' },
  { actor: 'Terrorist Cell Z', capability: 40, intent: 95, risk: 'medium' },
  { actor: 'Hacktivist Collective', capability: 55, intent: 50, risk: 'low' },
];

const correlations = [
  { id: 'CL-01', disciplines: 'SIGINT + HUMINT', finding: 'Identified comms pattern matches source report', confidence: 92, status: 'Confirmed' },
  { id: 'CL-02', disciplines: 'OSINT + IMINT', finding: 'Satellite imagery corroborates open-source claims', confidence: 87, status: 'Pending' },
  { id: 'CL-03', disciplines: 'HUMINT + CYBER', finding: 'Human source links to known C2 infrastructure', confidence: 78, status: 'Investigating' },
  { id: 'CL-04', disciplines: 'SIGINT + FININT', finding: 'Financial transactions correlate with intercepts', confidence: 95, status: 'Confirmed' },
];

const nies = [
  { title: 'NIE 2026-01: Regional Stability Assessment', date: '2026-01-15', classification: 'TS//SCI' },
  { title: 'NIE 2026-02: Cyber Threat Landscape', date: '2026-03-22', classification: 'Top Secret' },
  { title: 'NIE 2026-03: Foreign Military Modernization', date: '2026-05-10', classification: 'TS//SCI' },
  { title: 'NIE 2026-04: Economic Security Outlook', date: '2026-06-05', classification: 'Secret' },
];

const FusionCenterDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">All-Source Intelligence Fusion</h2>
          <p className="text-gray-500 mt-1">Cross-discipline analysis and national intelligence products</p>
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
            <Layers className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Products by Classification</h3>
          </div>
          <div className="space-y-4">
            {classifications.map((c, i) => (
              <motion.div
                key={c.level}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${c.color}`}
              >
                <div className="flex justify-between items-center">
                  <span className="text-sm font-bold">{c.level}</span>
                  <span className="text-2xl font-black">{c.count}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="col-span-2 bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Activity className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Threat Actor Heat Map (Capability × Intent)</h3>
          </div>
          <div className="space-y-3">
            {threatHeatMap.map((t, i) => (
              <motion.div
                key={t.actor}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-5 rounded-2xl border ${
                  t.risk === 'critical' ? 'border-red-500/20 bg-red-500/5' :
                  t.risk === 'high' ? 'border-orange-500/20 bg-orange-500/5' :
                  t.risk === 'medium' ? 'border-yellow-500/20 bg-yellow-500/5' :
                  'border-white/5 bg-white/[0.02]'
                }`}
              >
                <div className="flex justify-between items-center mb-3">
                  <span className="font-bold">{t.actor}</span>
                  <span className={`text-[10px] font-black uppercase tracking-wider ${
                    t.risk === 'critical' ? 'text-red-500' :
                    t.risk === 'high' ? 'text-orange-500' :
                    t.risk === 'medium' ? 'text-yellow-500' : 'text-green-500'
                  }`}>{t.risk}</span>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <div className="flex justify-between text-xs mb-1">
                      <span className="text-gray-400">Capability</span>
                      <span className="font-bold">{t.capability}%</span>
                    </div>
                    <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{ width: `${t.capability}%` }}
                        transition={{ duration: 0.6 }}
                        className="h-full rounded-full bg-red-500"
                      />
                    </div>
                  </div>
                  <div>
                    <div className="flex justify-between text-xs mb-1">
                      <span className="text-gray-400">Intent</span>
                      <span className="font-bold">{t.intent}%</span>
                    </div>
                    <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{ width: `${t.intent}%` }}
                        transition={{ duration: 0.6 }}
                        className="h-full rounded-full bg-orange-500"
                      />
                    </div>
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <GitBranch className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Cross-Discipline Correlations</h3>
          </div>
          <div className="space-y-3">
            {correlations.map((c, i) => (
              <motion.div
                key={c.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className="p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div className="flex justify-between items-center mb-2">
                  <span className="font-mono text-xs text-gray-500">{c.id}</span>
                  <span className={`text-[10px] font-bold uppercase ${
                    c.status === 'Confirmed' ? 'text-green-400' :
                    c.status === 'Pending' ? 'text-yellow-400' : 'text-orange-400'
                  }`}>{c.status}</span>
                </div>
                <div className="text-xs text-blue-400 font-bold mb-1">{c.disciplines}</div>
                <div className="text-sm mb-2">{c.finding}</div>
                <div className="flex items-center gap-2 text-xs text-gray-500">
                  <span>Confidence:</span>
                  <div className="h-1.5 w-20 bg-white/5 rounded-full overflow-hidden">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${c.confidence}%` }}
                      transition={{ duration: 0.6 }}
                      className="h-full rounded-full bg-green-500"
                    />
                  </div>
                  <span className="font-bold">{c.confidence}%</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <BookOpen className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Recent NIE Publications</h3>
          </div>
          <div className="space-y-3">
            {nies.map((n, i) => (
              <motion.div
                key={n.title}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className="p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div className="flex justify-between items-start mb-2">
                  <div className="font-bold text-sm flex-1 mr-4">{n.title}</div>
                  <span className={`text-[10px] font-black uppercase whitespace-nowrap ${
                    n.classification === 'TS//SCI' ? 'text-red-400' :
                    n.classification === 'Top Secret' ? 'text-orange-400' : 'text-yellow-400'
                  }`}>{n.classification}</span>
                </div>
                <div className="text-xs text-gray-500">{n.date}</div>
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

export default FusionCenterDashboard;
