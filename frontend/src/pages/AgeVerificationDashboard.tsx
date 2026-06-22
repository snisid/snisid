import React from 'react';
import { BadgeCheck, XCircle, Users, Clock, AlertTriangle, Shield } from 'lucide-react';
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

const AgeVerificationDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Age Verification</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Digital Age Assurance Platform</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
          <span className="text-emerald-500 text-xs font-bold uppercase tracking-wider">Compliant</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Verifications Today" value="8,234" icon={<BadgeCheck className="w-5 h-5 text-blue-500" />} trend="+12% vs yesterday" />
        <KpiCard title="Age Verified" value="7,891" icon={<Shield className="w-5 h-5 text-emerald-500" />} trend="95.8% pass rate" />
        <KpiCard title="Underage Flagged" value="343" icon={<XCircle className="w-5 h-5 text-red-500" />} trend="4.2% blocked" />
        <KpiCard title="Active Merchants" value="1,247" icon={<Users className="w-5 h-5 text-violet-500" />} trend="+28 this week" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><Clock className="w-4 h-4 text-amber-500" /> Recent Verifications</h3>
          <div className="space-y-2">
            {[
              { merchant: 'Liquor Store #147', method: 'QR Scan UL', result: 'Passed', age: '32', time: '30s ago' },
              { merchant: 'Vape Shop #23', method: 'BLE mDL', result: 'Passed', age: '24', time: '2m ago' },
              { merchant: 'Online CBD Store', method: 'ID Photo', result: 'Flagged', age: '17', time: '5m ago' },
              { merchant: 'Casino Entrance B', method: 'NFC', result: 'Passed', age: '45', time: '8m ago' },
            ].map((v) => (
              <div key={v.merchant} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{v.merchant}</div>
                  <div className="text-xs text-gray-500">{v.method} - Age {v.age}</div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-[10px] text-gray-500">{v.time}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${v.result === 'Passed' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-red-500/10 text-red-500'}`}>{v.result}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Compliance Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'KYC Pass Rate', val: '95.8%', color: 'bg-emerald-500' },
              { label: 'Response Time', val: '340ms', color: 'bg-blue-500' },
              { label: 'API Uptime', val: '99.97%', color: 'bg-violet-500' },
              { label: 'Fraud Attempts Blocked', val: '1,247', color: 'bg-amber-500' },
            ].map((m) => (
              <div key={m.label}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-400">{m.label}</span>
                  <span className="text-white font-bold">{m.val}</span>
                </div>
                <div className="h-1.5 bg-white/5 rounded-full overflow-hidden">
                  <div className={`h-full ${m.color} rounded-full`} style={{ width: '88%' }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AgeVerificationDashboard;
