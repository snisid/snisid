import React from 'react';
import { Headphones, Ticket, CheckCircle, Clock, AlertTriangle, Users } from 'lucide-react';
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

const ServiceDeskDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Service Desk</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">IT Support & Incident Management</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-500/10 border border-blue-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
          <span className="text-blue-500 text-xs font-bold uppercase tracking-wider">5 Agents Online</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Open Tickets" value="127" icon={<Ticket className="w-5 h-5 text-blue-500" />} trend="+23 today" />
        <KpiCard title="Resolved Today" value="89" icon={<CheckCircle className="w-5 h-5 text-emerald-500" />} trend="78% resolution rate" />
        <KpiCard title="SLA Breached" value="7" icon={<AlertTriangle className="w-5 h-5 text-red-500" />} trend="5.5% breach rate" />
        <KpiCard title="Active Agents" value="12" icon={<Headphones className="w-5 h-5 text-violet-500" />} trend="5 online, 7 away" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><Clock className="w-4 h-4 text-amber-500" /> Pending Tickets</h3>
          <div className="space-y-2">
            {[
              { ticket: 'INC-2026-08421', subject: 'VPN Access Issue', priority: 'Critical', age: '45m', agent: 'Unassigned' },
              { ticket: 'INC-2026-08420', subject: 'Password Reset Request', priority: 'High', age: '2h', agent: 'J. Davis' },
              { ticket: 'INC-2026-08419', subject: 'Biometric Scanner Calibration', priority: 'Medium', age: '4h', agent: 'M. Lee' },
              { ticket: 'INC-2026-08418', subject: 'Email Delivery Delay', priority: 'Low', age: '1d', agent: 'S. Wilson' },
            ].map((t) => (
              <div key={t.ticket} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{t.subject}</div>
                  <div className="text-xs text-gray-500">{t.ticket} - {t.agent}</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-[10px] text-gray-500">{t.age}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${t.priority === 'Critical' ? 'bg-red-500/10 text-red-500' : t.priority === 'High' ? 'bg-amber-500/10 text-amber-500' : t.priority === 'Medium' ? 'bg-yellow-500/10 text-yellow-500' : 'bg-gray-500/10 text-gray-400'}`}>{t.priority}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Service Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'Avg Response Time', val: '12m', color: 'bg-blue-500' },
              { label: 'Avg Resolution Time', val: '4.2h', color: 'bg-violet-500' },
              { label: 'First Contact Resolution', val: '72%', color: 'bg-emerald-500' },
              { label: 'CSAT Score', val: '4.6/5.0', color: 'bg-cyan-500' },
            ].map((m) => (
              <div key={m.label}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-400">{m.label}</span>
                  <span className="text-white font-bold">{m.val}</span>
                </div>
                <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                  <div className={`h-full ${m.color} rounded-full`} style={{ width: '85%' }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ServiceDeskDashboard;
