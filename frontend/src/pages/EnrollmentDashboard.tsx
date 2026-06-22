import React from 'react';
import { UserPlus, CheckCircle, Clock, XCircle, Users, FileText } from 'lucide-react';
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

const EnrollmentDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Enrollment RA Workflow</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Registration Authority Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-500/10 border border-blue-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
          <span className="text-blue-500 text-xs font-bold uppercase tracking-wider">Processing</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Pending Enrollments" value="247" icon={<Clock className="w-5 h-5 text-amber-500" />} trend="+18 today" />
        <KpiCard title="Approved Today" value="156" icon={<CheckCircle className="w-5 h-5 text-emerald-500" />} trend="92% approval rate" />
        <KpiCard title="Rejected" value="12" icon={<XCircle className="w-5 h-5 text-red-500" />} trend="4.8% rejection rate" />
        <KpiCard title="Total RA Operators" value="34" icon={<Users className="w-5 h-5 text-blue-500" />} trend="6 online now" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><FileText className="w-4 h-4 text-blue-500" /> Pending Reviews</h3>
          <div className="space-y-2">
            {[
              { id: 'ENR-2026-08472', name: 'J. Smith', type: 'New Credential', submitted: '2h ago', priority: 'High' },
              { id: 'ENR-2026-08471', name: 'M. Rodriguez', type: 'Renewal', submitted: '3h ago', priority: 'Normal' },
              { id: 'ENR-2026-08470', name: 'A. Patel', type: 'Replacement', submitted: '5h ago', priority: 'Low' },
              { id: 'ENR-2026-08469', name: 'K. Williams', type: 'New Credential', submitted: '6h ago', priority: 'High' },
            ].map((e) => (
              <div key={e.id} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{e.name}</div>
                  <div className="text-xs text-gray-500">{e.id} - {e.type}</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-[10px] text-gray-500">{e.submitted}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${e.priority === 'High' ? 'bg-red-500/10 text-red-500' : e.priority === 'Normal' ? 'bg-amber-500/10 text-amber-500' : 'bg-gray-500/10 text-gray-400'}`}>{e.priority}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Enrollment Stats</h3>
          <div className="space-y-4">
            {[
              { label: 'Biometric Capture Rate', val: '97%', color: 'bg-emerald-500' },
              { label: 'Doc Verification Rate', val: '89%', color: 'bg-blue-500' },
              { label: 'Fraud Flag Rate', val: '2.1%', color: 'bg-amber-500' },
              { label: 'Avg Processing Time', val: '4.2m', color: 'bg-violet-500' },
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

export default EnrollmentDashboard;
