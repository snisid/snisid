import React, { useState } from 'react';
import { Search, Filter, MoreVertical, ShieldCheck, AlertCircle, Fingerprint, Eye } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const IdentitySearch = () => {
  const [query, setQuery] = useState('');
  
  const mockIdentities = [
    { id: 'HT-921-002', name: 'Alix Desrosiers', dob: '1988-04-12', gender: 'M', status: 'verified', risk: 'low' },
    { id: 'HT-452-991', name: 'Marie Claudette', dob: '1992-11-23', gender: 'F', status: 'pending', risk: 'medium' },
    { id: 'HT-112-404', name: 'Unknown Entity', dob: 'N/A', gender: 'U', status: 'suspicious', risk: 'high' },
  ];

  return (
    <div className="p-8 space-y-8 h-screen flex flex-col">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Identity Search</h2>
          <p className="text-gray-500 mt-1">Query the national identity database with real-time biometric matching</p>
        </div>
      </header>

      <div className="flex gap-4">
        <div className="flex-1 relative">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
          <input 
            type="text" 
            placeholder="Search by ID, name, or biometric hash..." 
            className="w-full bg-[#0f1218] border border-white/10 rounded-2xl py-4 pl-12 pr-4 focus:outline-none focus:border-blue-500/50 transition-all shadow-inner shadow-black/20"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <button className="px-6 bg-[#0f1218] border border-white/10 rounded-2xl flex items-center gap-2 hover:bg-white/5 transition-all">
          <Filter className="w-4 h-4" /> Filters
        </button>
      </div>

      <div className="flex-1 min-h-0 bg-[#0f1218] rounded-[2rem] border border-white/5 overflow-hidden flex flex-col">
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-white/5 text-[10px] uppercase tracking-widest text-gray-500 font-bold">
                <th className="px-8 py-6">Identifier</th>
                <th className="px-8 py-6">Full Name</th>
                <th className="px-8 py-6">DOB / Gender</th>
                <th className="px-8 py-6">Status</th>
                <th className="px-8 py-6 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/[0.02]">
              {mockIdentities.map((item, i) => (
                <motion.tr 
                  key={item.id}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: i * 0.05 }}
                  className="hover:bg-white/[0.01] transition-all group cursor-pointer"
                >
                  <td className="px-8 py-6">
                    <span className="font-mono text-xs text-blue-400 bg-blue-400/5 px-2 py-1 rounded border border-blue-400/10">
                      {item.id}
                    </span>
                  </td>
                  <td className="px-8 py-6 font-medium">{item.name}</td>
                  <td className="px-8 py-6 text-sm text-gray-400">
                    {item.dob} <span className="mx-2 text-gray-700">|</span> {item.gender}
                  </td>
                  <td className="px-8 py-6">
                    <StatusBadge status={item.status} />
                  </td>
                  <td className="px-8 py-6 text-right">
                    <div className="flex justify-end gap-2">
                      <button className="p-2 hover:bg-white/5 rounded-lg text-gray-500 hover:text-white transition-all">
                        <Eye className="w-4 h-4" />
                      </button>
                      <button className="p-2 hover:bg-white/5 rounded-lg text-gray-500 hover:text-white transition-all">
                        <Fingerprint className="w-4 h-4" />
                      </button>
                      <button className="p-2 hover:bg-white/5 rounded-lg text-gray-500 hover:text-white transition-all">
                        <MoreVertical className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

const StatusBadge = ({ status }: { status: string }) => {
  const configs: any = {
    verified: { icon: ShieldCheck, color: 'text-green-400', bg: 'bg-green-400/10', border: 'border-green-400/20' },
    pending: { icon: AlertCircle, color: 'text-orange-400', bg: 'bg-orange-400/10', border: 'border-orange-400/20' },
    suspicious: { icon: AlertCircle, color: 'text-red-400', bg: 'bg-red-400/10', border: 'border-red-400/20' },
  };
  const config = configs[status] || configs.pending;
  const Icon = config.icon;

  return (
    <div className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-tighter ${config.bg} ${config.color} border ${config.border}`}>
      <Icon className="w-3 h-3" />
      {status}
    </div>
  );
};

export default IdentitySearch;
