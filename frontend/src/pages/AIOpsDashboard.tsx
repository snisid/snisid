import React from 'react';
import { motion } from 'framer-motion';
import { 
  ResponsiveContainer, 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  AreaChart, 
  Area 
} from 'recharts';
import { Activity, Shield, AlertCircle, Zap, TrendingUp, BarChart3 } from 'lucide-react';

const data = [
  { name: '10:00', load: 400, predicted: 420, health: 98 },
  { name: '10:05', load: 300, predicted: 310, health: 99 },
  { name: '10:10', load: 600, predicted: 580, health: 97 },
  { name: '10:15', load: 800, predicted: 850, health: 95 },
  { name: '10:20', load: 500, predicted: 520, health: 98 },
  { name: '10:25', load: 400, predicted: 410, health: 99 },
];

const AIOpsDashboard = () => {
  return (
    <div className="p-8 space-y-8 h-screen overflow-y-auto custom-scrollbar pb-24">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-black tracking-tighter text-white">Observability & AIOps Intelligence</h2>
          <p className="text-gray-500 mt-1">Predictive infrastructure monitoring and automated health scoring</p>
        </div>
        <div className="flex gap-3">
          <div className="px-4 py-2 bg-blue-500/10 border border-blue-500/20 rounded-xl flex items-center gap-2">
             <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse" />
             <span className="text-[10px] font-black uppercase tracking-widest text-blue-500">Tempo Tracing Active</span>
          </div>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        <HealthMetric label="Global Health Score" value="98.4%" color="text-emerald-400" icon={Activity} />
        <HealthMetric label="Inference Latency" value="12ms" color="text-blue-400" icon={Zap} />
        <HealthMetric label="Active Anomaly Clusters" value="0" color="text-emerald-500" icon={Shield} />
        <HealthMetric label="Predicted Load Spike" value="Low" color="text-emerald-400" icon={TrendingUp} />
      </div>

      <div className="grid grid-cols-2 gap-8">
        {/* Predictive Traffic Forecast */}
        <div className="p-8 bg-[#0f1218] rounded-[2.5rem] border border-white/5">
          <div className="flex justify-between items-center mb-8">
            <h3 className="text-xl font-bold">Predictive Traffic Scaling</h3>
            <span className="text-[10px] font-black uppercase tracking-widest text-gray-500">Holt-Winters Forecast</span>
          </div>
          <div className="h-[300px]">
             <ResponsiveContainer width="100%" height="100%">
               <AreaChart data={data}>
                 <defs>
                   <linearGradient id="colorLoad" x1="0" y1="0" x2="0" y2="1">
                     <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
                     <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                   </linearGradient>
                 </defs>
                 <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#1f2937" />
                 <XAxis dataKey="name" stroke="#4b5563" fontSize={11} tickLine={false} axisLine={false} />
                 <YAxis stroke="#4b5563" fontSize={11} tickLine={false} axisLine={false} />
                 <Tooltip contentStyle={{ backgroundColor: '#111827', border: '1px solid #374151', borderRadius: '16px' }} />
                 <Area type="monotone" dataKey="load" stroke="#3b82f6" strokeWidth={4} fill="url(#colorLoad)" name="Actual Load" />
                 <Line type="monotone" dataKey="predicted" stroke="#6366f1" strokeDasharray="5 5" name="Predicted" />
               </AreaChart>
             </ResponsiveContainer>
          </div>
        </div>

        {/* Distributed Tracing Graph Placeholder */}
        <div className="p-8 bg-[#0f1218] rounded-[2.5rem] border border-white/5">
           <div className="flex justify-between items-center mb-8">
             <h3 className="text-xl font-bold">Span Latency (Tracing)</h3>
             <BarChart3 className="w-5 h-5 text-gray-600" />
           </div>
           <div className="space-y-6">
              <TraceRow label="gateway -> identity-api" value="4ms" width="w-[20%]" />
              <TraceRow label="identity-api -> postgres" value="12ms" width="w-[45%]" />
              <TraceRow label="identity-api -> kafka" value="8ms" width="w-[30%]" />
              <TraceRow label="fraud-engine -> neo4j" value="56ms" width="w-[85%]" />
           </div>
        </div>
      </div>
    </div>
  );
};

const HealthMetric = ({ label, value, color, icon: Icon }: { label: string, value: string, color: string, icon: any }) => (
  <div className="p-6 bg-[#0f1218] rounded-3xl border border-white/5">
    <div className="flex justify-between items-start mb-2">
      <p className="text-[10px] font-black uppercase tracking-widest text-gray-500">{label}</p>
      <Icon className={`w-5 h-5 ${color}`} />
    </div>
    <h3 className="text-3xl font-black tracking-tighter text-white">{value}</h3>
  </div>
);

const TraceRow = ({ label, value, width }: { label: string, value: string, width: string }) => (
  <div className="space-y-2">
    <div className="flex justify-between text-xs">
      <span className="text-gray-400 font-mono">{label}</span>
      <span className="text-blue-400 font-bold">{value}</span>
    </div>
    <div className="w-full bg-white/5 h-1.5 rounded-full overflow-hidden">
      <div className={`bg-blue-500 h-full ${width} rounded-full`} />
    </div>
  </div>
);

export default AIOpsDashboard;
