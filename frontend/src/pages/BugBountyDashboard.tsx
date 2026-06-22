import React from 'react';
import { Bug, Trophy, Users, DollarSign, AlertTriangle, CheckCircle } from 'lucide-react';
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

const BugBountyDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Bug Bounty Program</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Vulnerability Disclosure & Rewards</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-amber-500/10 border border-amber-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-amber-500" />
          <span className="text-amber-500 text-xs font-bold uppercase tracking-wider">25 Reports Open</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Total Reports" value="847" icon={<Bug className="w-5 h-5 text-blue-500" />} trend="+12 this week" />
        <KpiCard title="Valid Findings" value="312" icon={<CheckCircle className="w-5 h-5 text-emerald-500" />} trend="36.8% valid rate" />
        <KpiCard title="Bounties Paid" value="$284,500" icon={<DollarSign className="w-5 h-5 text-amber-500" />} trend="$12,400 pending" />
        <KpiCard title="Active Researchers" value="89" icon={<Users className="w-5 h-5 text-violet-500" />} trend="7 new this month" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><AlertTriangle className="w-4 h-4 text-amber-500" /> Recent Submissions</h3>
          <div className="space-y-2">
            {[
              { researcher: '0xSecurityPro', vuln: 'SQL Injection - Audit API', severity: 'Critical', reward: '$15,000' },
              { researcher: 'CryptoCrack', vuln: 'JWT Key Confusion', severity: 'High', reward: '$8,000' },
              { researcher: 'WebWatcher', vuln: 'XSS in Enrollment Portal', severity: 'Medium', reward: '$2,500' },
              { researcher: 'NetNinja', vuln: 'Subdomain Takeover', severity: 'Low', reward: '$500' },
            ].map((s) => (
              <div key={s.researcher} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{s.vuln}</div>
                  <div className="text-xs text-gray-500">by {s.researcher}</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-amber-500 font-bold">{s.reward}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${s.severity === 'Critical' ? 'bg-red-500/10 text-red-500' : s.severity === 'High' ? 'bg-amber-500/10 text-amber-500' : s.severity === 'Medium' ? 'bg-yellow-500/10 text-yellow-500' : 'bg-gray-500/10 text-gray-400'}`}>{s.severity}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Program Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'Avg Time to Triage', val: '4.2h', color: 'bg-blue-500' },
              { label: 'Avg Time to Fix', val: '3.7 days', color: 'bg-violet-500' },
              { label: 'Critical Fix Rate', val: '100%', color: 'bg-emerald-500' },
              { label: 'Researcher Retention', val: '76%', color: 'bg-cyan-500' },
            ].map((m) => (
              <div key={m.label}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-400">{m.label}</span>
                  <span className="text-white font-bold">{m.val}</span>
                </div>
                <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                  <div className={`h-full ${m.color} rounded-full`} style={{ width: '80%' }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default BugBountyDashboard;
