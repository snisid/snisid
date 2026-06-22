import React from 'react';
import { Briefcase, CheckCircle, XCircle, Clock, Users, FileSearch } from 'lucide-react';
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

const EmployeeVerifyDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">E-Verify Dashboard</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Employment Eligibility Verification</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-500/10 border border-blue-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
          <span className="text-blue-500 text-xs font-bold uppercase tracking-wider">Federated</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Verifications" value="3,892" icon={<Briefcase className="w-5 h-5 text-blue-500" />} trend="This month" />
        <KpiCard title="Employment Authorized" value="3,647" icon={<CheckCircle className="w-5 h-5 text-emerald-500" />} trend="93.7% authorized" />
        <KpiCard title="Tentative Non-Confirm" value="178" icon={<XCircle className="w-5 h-5 text-red-500" />} trend="4.6% TNC rate" />
        <KpiCard title="Pending SSA" value="67" icon={<Clock className="w-5 h-5 text-amber-500" />} trend="Awaiting response" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><Users className="w-4 h-4 text-blue-500" /> Recent Cases</h3>
          <div className="space-y-2">
            {[
              { employee: 'J. Chen', employer: 'TechCorp Inc', status: 'Authorized', ead: '2028-06-15' },
              { employee: 'M. Gonzalez', employer: 'BuildRight LLC', status: 'TNC', ead: 'Pending' },
              { employee: 'S. Patel', employer: 'HealthFirst', status: 'Authorized', ead: '2027-11-30' },
              { employee: 'L. Johnson', employer: 'EduStar Academy', status: 'SSA Review', ead: 'Pending' },
            ].map((c) => (
              <div key={c.employee} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{c.employee}</div>
                  <div className="text-xs text-gray-500">{c.employer}</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-[10px] text-gray-500">{c.ead}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${c.status === 'Authorized' ? 'bg-emerald-500/10 text-emerald-500' : c.status === 'TNC' ? 'bg-red-500/10 text-red-500' : 'bg-amber-500/10 text-amber-500'}`}>{c.status}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Case Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'Auto-Authorized Rate', val: '87%', color: 'bg-emerald-500' },
              { label: 'SSA Response Time', val: '2.4 days', color: 'bg-blue-500' },
              { label: 'TNC Contest Rate', val: '23%', color: 'bg-amber-500' },
              { label: 'Employer Compliance', val: '96%', color: 'bg-violet-500' },
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

export default EmployeeVerifyDashboard;
