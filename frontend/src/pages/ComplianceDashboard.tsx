import { 
  FileCheck, 
  ShieldCheck, 
  History, 
  Trash2, 
  Scale, 
  ArrowRight,
} from 'lucide-react';

const ComplianceDashboard = () => {
  const auditLogs = [
    { id: 'LOG-881', actor: 'OFFICER-01', action: 'Identity Read', target: 'HT-921-002', status: 'compliant' },
    { id: 'LOG-882', actor: 'SYSTEM-DAEMON', action: 'Data Retention', target: 'Archived 142 records', status: 'compliant' },
    { id: 'LOG-883', actor: 'AGENCY-PRP', action: 'Biometric Batch', target: '48 identities', status: 'warning' },
  ];

  return (
    <div className="p-8 space-y-8 h-screen flex flex-col overflow-hidden">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-black tracking-tighter text-white">Data Governance & Compliance</h2>
          <p className="text-gray-500 mt-1">GDPR-aligned consent management and national audit readiness</p>
        </div>
        <div className="flex gap-3">
          <button className="flex items-center gap-2 px-4 py-2 bg-blue-600 rounded-xl text-sm font-bold shadow-lg shadow-blue-900/20">
            <FileCheck className="w-4 h-4" /> Generate Audit Report
          </button>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        <GovernanceStat label="Compliance Score" value="99.2%" icon={ShieldCheck} color="text-emerald-400" />
        <GovernanceStat label="Consent Rate" value="84.1%" icon={Scale} color="text-blue-400" />
        <GovernanceStat label="Retention Violations" value="0" icon={Trash2} color="text-emerald-500" />
        <GovernanceStat label="Open Access Requests" value="12" icon={History} color="text-orange-400" />
      </div>

      <div className="flex-1 min-h-0 grid grid-cols-3 gap-8">
        {/* Real-time Audit Trail */}
        <div className="col-span-2 bg-[#0f1218] rounded-[2.5rem] border border-white/5 flex flex-col overflow-hidden">
          <div className="p-6 border-b border-white/5 flex items-center justify-between">
            <h3 className="font-bold">National Data Access Log</h3>
            <span className="text-[10px] font-black uppercase tracking-widest text-gray-500">Live Feed</span>
          </div>
          <div className="flex-1 overflow-x-auto">
             <table className="w-full text-left">
               <thead className="text-[10px] font-black uppercase tracking-widest text-gray-500 border-b border-white/5">
                 <tr>
                   <th className="px-6 py-4">ID</th>
                   <th className="px-6 py-4">Actor</th>
                   <th className="px-6 py-4">Action</th>
                   <th className="px-6 py-4">Target</th>
                   <th className="px-6 py-4 text-right">Status</th>
                 </tr>
               </thead>
               <tbody className="divide-y divide-white/[0.02] text-sm">
                 {auditLogs.map((log) => (
                   <tr key={log.id} className="hover:bg-white/[0.01] transition-all">
                     <td className="px-6 py-4 font-mono text-xs text-blue-400">{log.id}</td>
                     <td className="px-6 py-4 font-medium">{log.actor}</td>
                     <td className="px-6 py-4 text-gray-400">{log.action}</td>
                     <td className="px-6 py-4 text-gray-500">{log.target}</td>
                     <td className="px-6 py-4 text-right">
                        <span className={`px-2 py-1 rounded-full text-[10px] font-black uppercase tracking-tighter ${
                          log.status === 'compliant' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-orange-500/10 text-orange-500'
                        }`}>
                          {log.status}
                        </span>
                     </td>
                   </tr>
                 ))}
               </tbody>
             </table>
          </div>
        </div>

        {/* Governance Controls */}
        <div className="space-y-6">
           <ControlCard 
              title="Data Retention Policy" 
              status="Enforced" 
              desc="Automatic anonymization of identity records older than 10 years." 
              action="Configure"
           />
           <ControlCard 
              title="Data Subject Access" 
              status="Active" 
              desc="Automated handling of citizen data disclosure requests." 
              action="Manage"
           />
           <div className="p-8 bg-blue-600/10 border border-blue-500/20 rounded-[2rem] space-y-4">
              <h4 className="font-bold">Compliance Readiness</h4>
              <p className="text-xs text-gray-400 leading-relaxed">System is currently aligned with 14/15 ID4D Principles and 100% of National Data Protection laws.</p>
              <button className="flex items-center gap-2 text-xs font-bold text-blue-400 hover:text-blue-300 transition-all">
                View Full Audit History <ArrowRight className="w-3 h-3" />
              </button>
           </div>
        </div>
      </div>
    </div>
  );
};

const GovernanceStat = ({ label, value, icon: Icon, color }: { label: string, value: string, icon: any, color: string }) => (
  <div className="p-6 bg-[#0f1218] border border-white/5 rounded-3xl">
    <div className="flex justify-between items-start mb-2">
       <span className="text-[10px] font-black uppercase tracking-widest text-gray-500">{label}</span>
       <Icon className={`w-5 h-5 ${color}`} />
    </div>
    <div className="text-3xl font-black text-white tabular-nums tracking-tighter">{value}</div>
  </div>
);

const ControlCard = ({ title, status, desc, action }: { title: string, status: string, desc: string, action: string }) => (
  <div className="p-6 bg-[#0f1218] border border-white/5 rounded-[2rem] space-y-3">
    <div className="flex justify-between items-center">
       <h4 className="font-bold">{title}</h4>
       <span className="text-[10px] font-black uppercase tracking-widest text-emerald-500 bg-emerald-500/10 px-2 py-0.5 rounded-full">{status}</span>
    </div>
    <p className="text-xs text-gray-500 leading-relaxed">{desc}</p>
    <button className="text-[10px] font-black uppercase tracking-widest text-blue-400 hover:text-white transition-all">{action}</button>
  </div>
);

export default ComplianceDashboard;
