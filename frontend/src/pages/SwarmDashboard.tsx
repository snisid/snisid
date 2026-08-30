import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Terminal, Shield, Zap, Search, Activity, Command } from 'lucide-react';

const SwarmDashboard = () => {
  const [logs, setLogs] = useState([
    { id: 1, agent: 'Coordinator', msg: 'System goal received: Analyze Cluster F-501', time: '10:00:01' },
    { id: 2, agent: 'Fraud Investigator', msg: 'Task T1: Analyzing relationship graph for F-501', time: '10:00:05' },
    { id: 3, agent: 'Threat Hunter', msg: 'Task T2: Scanning ingress logs for IP overlap', time: '10:00:12' },
  ]);

  const agents = [
    { name: 'Coordinator', status: 'active', task: 'Task Orchestration', color: 'text-blue-400' },
    { name: 'Threat Hunter', status: 'idle', task: 'N/A', color: 'text-emerald-400' },
    { name: 'Fraud Investigator', status: 'active', task: 'GNN Analysis', color: 'text-indigo-400' },
    { name: 'Incident Responder', status: 'idle', task: 'N/A', color: 'text-red-400' },
  ];

  return (
    <div className="p-8 space-y-8 h-screen flex flex-col overflow-hidden">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-black tracking-tighter">AI Agent Swarm</h2>
          <p className="text-gray-500 mt-1">Autonomous orchestration of specialized intelligence agents</p>
        </div>
        <div className="flex gap-3">
          <div className="flex items-center gap-2 px-4 py-2 bg-emerald-500/10 border border-emerald-500/20 rounded-xl">
             <div className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse" />
             <span className="text-xs font-bold text-emerald-500 uppercase tracking-widest">Swarm Sync Active</span>
          </div>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {agents.map((agent) => (
          <div key={agent.name} className="p-6 bg-[#0f1218] border border-white/5 rounded-3xl relative overflow-hidden group">
            <div className={`absolute top-0 right-0 w-1 h-full ${agent.status === 'active' ? 'bg-blue-500' : 'bg-gray-800'}`} />
            <div className="flex justify-between items-start mb-4">
               <h4 className={`font-bold ${agent.color}`}>{agent.name}</h4>
               <span className={`text-[10px] font-black uppercase tracking-widest ${agent.status === 'active' ? 'text-blue-500' : 'text-gray-600'}`}>
                 {agent.status}
               </span>
            </div>
            <p className="text-xs text-gray-500 uppercase tracking-widest font-bold mb-1">Current Task</p>
            <p className="text-sm font-medium">{agent.task}</p>
          </div>
        ))}
      </div>

      <div className="flex-1 min-h-0 grid grid-cols-3 gap-8">
        {/* Swarm Log Feed */}
        <div className="col-span-2 bg-[#0a0c10] rounded-[2.5rem] border border-white/5 flex flex-col overflow-hidden">
          <div className="p-6 border-b border-white/5 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Terminal className="w-4 h-4 text-blue-400" />
              <h3 className="font-bold">Inter-Agent Messaging Stream</h3>
            </div>
          </div>
          <div className="flex-1 overflow-y-auto p-6 font-mono text-sm space-y-4 custom-scrollbar">
            {logs.map((log) => (
              <div key={log.id} className="flex gap-4 group">
                <span className="text-gray-600 shrink-0">[{log.time}]</span>
                <span className="text-blue-400 shrink-0 font-bold">{log.agent}:</span>
                <span className="text-gray-300">{log.msg}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Swarm Controls */}
        <div className="space-y-6">
          <div className="p-8 bg-blue-600/10 border border-blue-500/20 rounded-[2rem] space-y-4">
            <div className="flex items-center gap-2 text-blue-400">
               <Command className="w-5 h-5" />
               <h3 className="font-bold">Command Swarm</h3>
            </div>
            <textarea 
              placeholder="Inject high-level objective..." 
              className="w-full h-32 bg-black/20 border border-white/10 rounded-2xl p-4 text-sm focus:outline-none focus:border-blue-500/50 transition-all resize-none"
            />
            <button className="w-full py-4 bg-blue-600 hover:bg-blue-700 rounded-2xl font-bold text-sm shadow-lg shadow-blue-900/20 transition-all">
              Initiate Operation
            </button>
          </div>

          <div className="p-8 bg-[#0f1218] border border-white/5 rounded-[2rem] space-y-4">
             <h4 className="font-bold text-gray-500 uppercase tracking-widest text-xs">Swarm Statistics</h4>
             <StatRow label="Collective Memory" value="1.2 GB" />
             <StatRow label="Task Throughput" value="142/hr" />
             <StatRow label="Consensus Rate" value="99.4%" />
          </div>
        </div>
      </div>
    </div>
  );
};

const StatRow = ({ label, value }: { label: string, value: string }) => (
  <div className="flex justify-between items-center text-sm">
    <span className="text-gray-500">{label}</span>
    <span className="font-bold text-white">{value}</span>
  </div>
);

export default SwarmDashboard;
