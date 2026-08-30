import React from 'react';
import { Award, CheckCircle, Clock, FileText, AlertTriangle, TrendingUp } from 'lucide-react';
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

const CertificationDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">NIST Certification</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Standards Compliance Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-amber-500/10 border border-amber-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-amber-500" />
          <span className="text-amber-500 text-xs font-bold uppercase tracking-wider">In Progress</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Active Certifications" value="24" icon={<Award className="w-5 h-5 text-blue-500" />} trend="3 pending renewal" />
        <KpiCard title="Compliant" value="21" icon={<CheckCircle className="w-5 h-5 text-emerald-500" />} trend="87.5% compliance" />
        <KpiCard title="In Review" value="3" icon={<Clock className="w-5 h-5 text-amber-500" />} trend="FIPS 140-3, SP 800-63" />
        <KpiCard title="Findings" value="7" icon={<AlertTriangle className="w-5 h-5 text-red-500" />} trend="2 critical" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><FileText className="w-4 h-4 text-blue-500" /> Certification Pipeline</h3>
          <div className="space-y-2">
            {[
              { standard: 'FIPS 140-3 Level 3', scope: 'HSM Module v2.1', status: 'In Review', due: '2026-08-15' },
              { standard: 'SP 800-63-4 IAL3', scope: 'Enrollment Platform', status: 'Scheduled', due: '2026-09-30' },
              { standard: 'FIPS 201-3 PIV', scope: 'Card Applet v4.2', status: 'Active', due: '2027-01-15' },
              { standard: 'ISO 24727-3', scope: 'Middleware API', status: 'Active', due: '2026-11-01' },
            ].map((c) => (
              <div key={c.standard} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{c.standard}</div>
                  <div className="text-xs text-gray-500">{c.scope}</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${c.status === 'Active' ? 'bg-emerald-500/10 text-emerald-500' : c.status === 'In Review' ? 'bg-amber-500/10 text-amber-500' : 'bg-gray-500/10 text-gray-400'}`}>{c.status}</span>
                  <span className="text-xs text-gray-500">{c.due}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Compliance Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'FIPS 140-3 Coverage', val: '94%', color: 'bg-emerald-500' },
              { label: 'SP 800-63 Compliance', val: '88%', color: 'bg-blue-500' },
              { label: 'Crypto Algorithm Audit', val: '100%', color: 'bg-violet-500' },
              { label: 'Physical Security', val: '96%', color: 'bg-cyan-500' },
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

export default CertificationDashboard;
