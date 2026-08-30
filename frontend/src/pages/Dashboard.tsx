import React, { useState } from 'react';
import { 
  Shield, 
  Users, 
  AlertTriangle, 
  ArrowUpRight,
  Zap,
  Globe
} from 'lucide-react';
import { 
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip
} from 'recharts';
import { motion, AnimatePresence } from 'framer-motion';
import { useStore } from '../store/useStore';

const chartData = [
  { name: '00:00', value: 400 },
  { name: '04:00', value: 300 },
  { name: '08:00', value: 900 },
  { name: '12:00', value: 1200 },
  { name: '16:00', value: 1500 },
  { name: '20:00', value: 800 },
  { name: '23:59', value: 600 },
];

const Dashboard = () => {
  const { alerts } = useStore();

  return (
    <div className="p-8 space-y-8 h-screen overflow-y-auto custom-scrollbar">
      {/* Header Area */}
      <header className="flex justify-between items-start">
        <div>
          <h2 className="text-3xl font-black tracking-tighter">National Intelligence Dashboard</h2>
          <div className="flex items-center gap-2 mt-1">
            <span className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse" />
            <span className="text-xs font-bold uppercase tracking-widest text-emerald-500/80">Live Surveillance Active</span>
          </div>
        </div>
        <div className="flex gap-4">
          <div className="px-4 py-2 bg-white/5 border border-white/10 rounded-2xl flex items-center gap-3">
            <Globe className="w-4 h-4 text-blue-400" />
            <span className="text-xs font-bold">Node: Port-au-Prince Alpha</span>
          </div>
        </div>
      </header>

      {/* Hero Stats */}
      <div className="grid grid-cols-4 gap-6">
        <StatCard label="Total Identities" value="24.8M" change="+2.4%" icon={Users} color="text-blue-400" />
        <StatCard label="Biometric Verifications" value="142K" change="+12.1%" icon={Shield} color="text-emerald-400" />
        <StatCard label="Anomaly detections" value="1.2K" change="-5.2%" icon={AlertTriangle} color="text-red-400" />
        <StatCard label="Network Latency" value="14ms" change="Stable" icon={Zap} color="text-indigo-400" />
      </div>

      <div className="grid grid-cols-3 gap-8">
        {/* Verification Traffic */}
        <div className="col-span-2 p-8 bg-[#0f1218] rounded-[2.5rem] border border-white/5">
          <div className="flex justify-between items-center mb-8">
            <div>
              <h3 className="text-xl font-bold tracking-tight">Identity Throughput</h3>
              <p className="text-sm text-gray-500">Real-time verification volume across national infrastructure</p>
            </div>
            <button className="p-3 bg-white/5 rounded-2xl hover:bg-white/10 transition-all">
              <ArrowUpRight className="w-5 h-5 text-gray-400" />
            </button>
          </div>
          <div className="h-[320px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={chartData}>
                <defs>
                  <linearGradient id="colorValue" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#1f2937" />
                <XAxis dataKey="name" stroke="#4b5563" fontSize={11} tickLine={false} axisLine={false} />
                <YAxis stroke="#4b5563" fontSize={11} tickLine={false} axisLine={false} />
                <Tooltip 
                  contentStyle={{ backgroundColor: '#111827', border: '1px solid #374151', borderRadius: '16px', boxShadow: '0 20px 25px -5px rgba(0,0,0,0.5)' }}
                  itemStyle={{ color: '#fff' }}
                />
                <Area type="monotone" dataKey="value" stroke="#3b82f6" strokeWidth={4} fillOpacity={1} fill="url(#colorValue)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Live Event Feed */}
        <div className="p-8 bg-[#0f1218] rounded-[2.5rem] border border-white/5 flex flex-col h-full">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-lg font-bold">Live Events</h3>
            <div className="w-2 h-2 bg-red-500 rounded-full animate-ping" />
          </div>
          <div className="space-y-4 flex-1 overflow-y-auto pr-2 custom-scrollbar">
            <AnimatePresence>
              {alerts.map((alert) => (
                <motion.div 
                  key={alert.id}
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  className="p-4 bg-white/[0.03] border border-white/5 rounded-2xl hover:bg-white/[0.06] transition-all cursor-pointer group"
                >
                  <div className="flex gap-3">
                    <div className={`w-1.5 h-1.5 rounded-full mt-2 flex-shrink-0 ${
                      alert.severity === 'high' ? 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.6)]' : 'bg-blue-500'
                    }`} />
                    <div>
                      <p className="text-sm font-medium leading-relaxed group-hover:text-blue-300 transition-colors">{alert.message}</p>
                      <div className="flex items-center gap-2 mt-2">
                        <span className="text-[10px] font-mono text-gray-500">{alert.time}</span>
                        <span className="text-[10px] font-black uppercase tracking-widest text-blue-500/80">{alert.type}</span>
                      </div>
                    </div>
                  </div>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>
        </div>
      </div>
    </div>
  );
};

const StatCard = ({ label, value, change, icon: Icon, color }: { label: string, value: string, change: string, icon: any, color: string }) => (
  <motion.div 
    whileHover={{ y: -4 }}
    className="p-6 bg-[#0f1218] rounded-[2rem] border border-white/5 hover:border-white/10 transition-all relative overflow-hidden group"
  >
    <div className={`absolute -right-4 -bottom-4 w-24 h-24 opacity-[0.03] group-hover:opacity-[0.06] transition-opacity ${color}`}>
      <Icon className="w-full h-full" />
    </div>
    <div className="flex justify-between items-start mb-4">
      <div className={`p-3 rounded-2xl bg-white/5 ${color}`}>
        <Icon className="w-5 h-5" />
      </div>
      <span className={`text-[10px] font-bold px-2 py-1 rounded-lg ${change.startsWith('+') ? 'bg-emerald-500/10 text-emerald-500' : 'bg-red-500/10 text-red-500'}`}>
        {change}
      </span>
    </div>
    <p className="text-xs font-bold uppercase tracking-widest text-gray-500">{label}</p>
    <h3 className="text-2xl font-black mt-1 tracking-tighter">{value}</h3>
  </motion.div>
);

export default Dashboard;
