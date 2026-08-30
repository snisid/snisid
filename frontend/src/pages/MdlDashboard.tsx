import React from 'react';
import { Smartphone, Download, QrCode, CheckCircle, Users, Activity } from 'lucide-react';
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

const MdlDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Mobile Driver's License</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">ISO 18013-5 mDL Platform</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-500/10 border border-blue-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
          <span className="text-blue-500 text-xs font-bold uppercase tracking-wider">Active</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="mDL Issued" value="45,230" icon={<Smartphone className="w-5 h-5 text-blue-500" />} trend="+1,240 this week" />
        <KpiCard title="Active Devices" value="38,912" icon={<Download className="w-5 h-5 text-emerald-500" />} trend="86% activation" />
        <KpiCard title="Verifications" value="12,478" icon={<QrCode className="w-5 h-5 text-violet-500" />} trend="Today" />
        <KpiCard title="Trusted Readers" value="1,892" icon={<CheckCircle className="w-5 h-5 text-cyan-500" />} trend="+47 this month" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><Activity className="w-4 h-4 text-blue-500" /> Verification Activity</h3>
          <div className="space-y-2">
            {[
              { location: 'DMV Downtown', verifications: 847, status: 'Online', rate: '99.2%' },
              { location: 'Airport Security Check B', verifications: 623, status: 'Online', rate: '98.7%' },
              { location: 'Police Traffic Unit', verifications: 412, status: 'Online', rate: '100%' },
              { location: 'Age-Restricted Store #42', verifications: 289, status: 'Degraded', rate: '87.3%' },
            ].map((v) => (
              <div key={v.location} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{v.location}</div>
                  <div className="text-xs text-gray-500">{v.verifications} verifications</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-gray-500">{v.rate}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${v.status === 'Online' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-amber-500/10 text-amber-500'}`}>{v.status}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Platform Health</h3>
          <div className="space-y-4">
            {[
              { label: 'mDL Wallet Downloads', val: '45.2K', color: 'bg-blue-500' },
              { label: 'QR Scan Success Rate', val: '99.3%', color: 'bg-emerald-500' },
              { label: 'Response Time (P95)', val: '187ms', color: 'bg-violet-500' },
              { label: 'Data Sync Latency', val: '2.1s', color: 'bg-cyan-500' },
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

export default MdlDashboard;
