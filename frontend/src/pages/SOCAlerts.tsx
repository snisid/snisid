import React from 'react';
import { useStore } from '../store/useStore';
import { motion } from 'framer-motion';
import { AlertTriangle, Bell, CheckCircle2, Info, Clock, ExternalLink } from 'lucide-react';

const SOCAlerts = () => {
  const { alerts } = useStore();

  return (
    <div className="p-8 space-y-8 h-screen flex flex-col">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">SOC Command Center</h2>
          <p className="text-gray-500 mt-1">Real-time security operations and incident response dashboard</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        <StatusCard label="Active Incidents" value="12" color="text-red-400" icon={AlertTriangle} />
        <StatusCard label="Pending Review" value="48" color="text-orange-400" icon={Clock} />
        <StatusCard label="Resolved Today" value="156" color="text-green-400" icon={CheckCircle2} />
        <StatusCard label="System Integrity" value="99.9%" color="text-blue-400" icon={Info} />
      </div>

      <div className="flex-1 min-h-0 bg-[#0f1218] rounded-[2rem] border border-white/5 overflow-hidden flex flex-col p-8">
        <div className="flex justify-between items-center mb-8">
          <h3 className="text-xl font-bold">Priority Alerts</h3>
          <div className="flex gap-2">
            <button className="px-3 py-1.5 bg-white/5 rounded-lg text-xs font-bold uppercase tracking-widest text-gray-400 hover:text-white transition-all">All</button>
            <button className="px-3 py-1.5 bg-red-500/10 rounded-lg text-xs font-bold uppercase tracking-widest text-red-400 border border-red-500/20">Critical</button>
          </div>
        </div>

        <div className="space-y-4 overflow-y-auto pr-2 custom-scrollbar flex-1">
          {alerts.map((alert, i) => (
            <motion.div 
              key={alert.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.1 }}
              className={`p-6 border rounded-3xl flex items-center justify-between group hover:scale-[1.01] transition-all cursor-pointer ${
                alert.severity === 'high' ? 'bg-red-500/5 border-red-500/20' : 
                alert.severity === 'medium' ? 'bg-orange-500/5 border-orange-500/20' : 'bg-blue-500/5 border-blue-500/20'
              }`}
            >
              <div className="flex items-center gap-6">
                <div className={`p-4 rounded-2xl ${
                  alert.severity === 'high' ? 'bg-red-500/10 text-red-500' : 
                  alert.severity === 'medium' ? 'bg-orange-500/10 text-orange-500' : 'bg-blue-500/10 text-blue-500'
                }`}>
                  <Bell className="w-6 h-6" />
                </div>
                <div>
                  <h4 className="font-bold text-lg">{alert.message}</h4>
                  <div className="flex items-center gap-3 mt-1">
                    <span className="text-xs font-mono text-gray-500">{alert.time}</span>
                    <span className="w-1 h-1 bg-gray-700 rounded-full" />
                    <span className={`text-[10px] font-black uppercase tracking-[0.2em] ${
                      alert.severity === 'high' ? 'text-red-500' : 
                      alert.severity === 'medium' ? 'text-orange-500' : 'text-blue-500'
                    }`}>{alert.type}</span>
                  </div>
                </div>
              </div>
              <button className="p-3 bg-white/5 rounded-2xl text-gray-500 hover:text-white transition-all opacity-0 group-hover:opacity-100">
                <ExternalLink className="w-5 h-5" />
              </button>
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  );
};

const StatusCard = ({ label, value, color, icon: Icon }: { label: string, value: string, color: string, icon: any }) => (
  <div className="p-6 bg-[#0f1218] rounded-3xl border border-white/5">
    <div className="flex justify-between items-start mb-2">
      <p className="text-xs font-bold uppercase tracking-widest text-gray-500">{label}</p>
      <Icon className={`w-5 h-5 ${color}`} />
    </div>
    <h3 className="text-3xl font-black tracking-tighter">{value}</h3>
  </div>
);

export default SOCAlerts;
