import React from 'react';
import { motion } from 'framer-motion';
import { Search, AlertTriangle, Users, FileText, Shield, Eye, Activity, ClipboardList } from 'lucide-react';

const statCards = [
  { label: 'Background Investigations', value: '1,247', color: 'text-blue-400', icon: Search },
  { label: 'Insider Threat Alerts', value: '38', color: 'text-red-400', icon: AlertTriangle },
  { label: 'Foreign Contacts Filed', value: '624', color: 'text-yellow-400', icon: Users },
  { label: 'Adjudication Queue', value: '156', color: 'text-orange-400', icon: ClipboardList },
];

const investigations = [
  { stage: 'Initial Review', count: 423, color: 'text-blue-400 bg-blue-500/10 border-blue-500/20' },
  { stage: 'Field Work', count: 287, color: 'text-yellow-400 bg-yellow-500/10 border-yellow-500/20' },
  { stage: 'Adjudication', count: 156, color: 'text-orange-400 bg-orange-500/10 border-orange-500/20' },
  { stage: 'Completed', count: 381, color: 'text-green-400 bg-green-500/10 border-green-500/20' },
];

const threatAlerts = [
  { id: 'IT-01', type: 'Unauthorized Access', user: 'Analyst J. Doe', severity: 'critical', time: '10m ago' },
  { id: 'IT-02', type: 'Data Exfiltration Attempt', user: 'Contractor Smith', severity: 'high', time: '35m ago' },
  { id: 'IT-03', type: 'Policy Violation', user: 'Admin K. Chen', severity: 'medium', time: '1h ago' },
  { id: 'IT-04', type: 'Suspicious Login', user: 'External Account', severity: 'high', time: '3h ago' },
  { id: 'IT-05', type: 'Routine Audit Flag', user: 'System', severity: 'low', time: '5h ago' },
];

const foreignContacts = [
  { name: 'Embassy Attaché Event', reporters: 12, contacts: 34, date: '2026-06-18', status: 'Filed' },
  { name: 'International Conference', reporters: 8, contacts: 27, date: '2026-06-14', status: 'Filed' },
  { name: 'Academic Exchange', reporters: 5, contacts: 18, date: '2026-06-10', status: 'Pending' },
  { name: 'Trade Delegation', reporters: 9, contacts: 42, date: '2026-06-07', status: 'Filed' },
];

const adjudication = [
  { id: 'ADJ-893', applicant: 'M. Laurent', type: 'Initial Clearance', risk: 'low', days: '12' },
  { id: 'ADJ-894', applicant: 'P. Blanc', type: 'Renewal', risk: 'medium', days: '8' },
  { id: 'ADJ-895', applicant: 'S. Dubois', type: 'Initial Clearance', risk: 'high', days: '5' },
  { id: 'ADJ-896', applicant: 'L. Petit', type: 'Elevation', risk: 'medium', days: '3' },
  { id: 'ADJ-897', applicant: 'C. Leroy', type: 'Renewal', risk: 'low', days: '1' },
];

const CounterintelDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Counterintelligence</h2>
          <p className="text-gray-500 mt-1">Insider threat detection, background investigations, and foreign contact reporting</p>
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
            <Search className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Background Investigations Pipeline</h3>
          </div>
          <div className="grid grid-cols-2 gap-4">
            {investigations.map((inv, i) => (
              <motion.div
                key={inv.stage}
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: i * 0.05 }}
                className={`p-5 rounded-2xl border ${inv.color}`}
              >
                <div className="text-3xl font-black mb-1">{inv.count}</div>
                <div className="text-sm font-bold">{inv.stage}</div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <AlertTriangle className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Insider Threat Alerts</h3>
          </div>
          <div className="space-y-3">
            {threatAlerts.map((t, i) => (
              <motion.div
                key={t.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${
                  t.severity === 'critical' ? 'border-red-500/20 bg-red-500/5' :
                  t.severity === 'high' ? 'border-orange-500/20 bg-orange-500/5' :
                  t.severity === 'medium' ? 'border-yellow-500/20 bg-yellow-500/5' :
                  'border-white/5 bg-white/[0.02]'
                }`}
              >
                <div className="flex justify-between items-center mb-2">
                  <span className="font-bold text-sm">{t.type}</span>
                  <span className={`text-[10px] font-black uppercase tracking-wider ${
                    t.severity === 'critical' ? 'text-red-500' :
                    t.severity === 'high' ? 'text-orange-500' :
                    t.severity === 'medium' ? 'text-yellow-500' : 'text-green-500'
                  }`}>{t.severity}</span>
                </div>
                <div className="flex justify-between text-xs text-gray-400">
                  <span>{t.user}</span>
                  <span>{t.time}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Users className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Foreign Contacts Reporting</h3>
          </div>
          <div className="space-y-3">
            {foreignContacts.map((fc, i) => (
              <motion.div
                key={fc.name}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="flex items-center justify-between p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div>
                  <div className="font-bold text-sm">{fc.name}</div>
                  <div className="text-xs text-gray-500">
                    {fc.reporters} reporters · {fc.contacts} contacts
                  </div>
                </div>
                <div className="text-right">
                  <div className="text-xs text-gray-500">{fc.date}</div>
                  <span className={`text-xs font-bold ${
                    fc.status === 'Filed' ? 'text-green-400' : 'text-yellow-400'
                  }`}>{fc.status}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <ClipboardList className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Adjudication Queue</h3>
          </div>
          <div className="space-y-3">
            {adjudication.map((a, i) => (
              <motion.div
                key={a.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="flex items-center justify-between p-4 rounded-2xl bg-white/[0.02] border border-white/5"
              >
                <div>
                  <div className="flex items-center gap-3 mb-1">
                    <span className="font-mono text-xs text-gray-500">{a.id}</span>
                    <span className={`text-[10px] font-bold uppercase ${
                      a.risk === 'high' ? 'text-red-500' :
                      a.risk === 'medium' ? 'text-yellow-500' : 'text-green-500'
                    }`}>{a.risk}</span>
                  </div>
                  <div className="font-bold text-sm">{a.applicant}</div>
                  <div className="text-xs text-gray-500">{a.type}</div>
                </div>
                <div className="text-right">
                  <div className="text-xs text-gray-500">Queue</div>
                  <div className="text-sm font-mono font-bold">{a.days}d</div>
                </div>
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

export default CounterintelDashboard;
