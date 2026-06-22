import React from 'react';
import { Eye, CheckCircle, AlertTriangle, Users, FileText, Activity } from 'lucide-react';
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

const AccessibilityDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">WCAG Accessibility</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Digital Accessibility Compliance Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-emerald-500" />
          <span className="text-emerald-500 text-xs font-bold uppercase tracking-wider">WCAG 2.2 AA</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Pages Audited" value="342" icon={<FileText className="w-5 h-5 text-blue-500" />} trend="100% coverage" />
        <KpiCard title="WCAG Pass Rate" value="96.4%" icon={<CheckCircle className="w-5 h-5 text-emerald-500" />} trend="+2.1% this month" />
        <KpiCard title="Violations" value="47" icon={<AlertTriangle className="w-5 h-5 text-amber-500" />} trend="12 critical" />
        <KpiCard title="Screen Reader Users" value="1,247" icon={<Users className="w-5 h-5 text-violet-500" />} trend="+18% YoY" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><Eye className="w-4 h-4 text-blue-500" /> Violation Breakdown</h3>
          <div className="space-y-2">
            {[
              { guideline: '1.1.1 Non-text Content', level: 'A', count: 8, status: 'In Progress' },
              { guideline: '1.4.3 Contrast Minimum', level: 'AA', count: 14, status: 'Fixing' },
              { guideline: '2.1.1 Keyboard', level: 'A', count: 5, status: 'Open' },
              { guideline: '2.4.4 Link Purpose', level: 'A', count: 12, status: 'In Progress' },
              { guideline: '4.1.2 Name, Role, Value', level: 'A', count: 8, status: 'Open' },
            ].map((v) => (
              <div key={v.guideline} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{v.guideline}</div>
                  <div className="text-xs text-gray-500">Level {v.level} - {v.count} occurrences</div>
                </div>
                <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${v.status === 'Fixed' ? 'bg-emerald-500/10 text-emerald-500' : v.status === 'Fixing' ? 'bg-amber-500/10 text-amber-500' : 'bg-red-500/10 text-red-500'}`}>{v.status}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Compliance Score</h3>
          <div className="space-y-4">
            {[
              { label: 'WCAG 2.2 Level A', val: '97%', color: 'bg-emerald-500' },
              { label: 'WCAG 2.2 Level AA', val: '94%', color: 'bg-blue-500' },
              { label: 'Section 508 Compliance', val: '98%', color: 'bg-violet-500' },
              { label: 'EN 301 549', val: '95%', color: 'bg-cyan-500' },
            ].map((m) => (
              <div key={m.label}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-400">{m.label}</span>
                  <span className="text-white font-bold">{m.val}</span>
                </div>
                <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                  <div className={`h-full ${m.color} rounded-full`} style={{ width: m.val }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AccessibilityDashboard;
