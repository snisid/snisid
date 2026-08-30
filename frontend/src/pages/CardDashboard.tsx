import React from 'react';
import { CreditCard, Package, CheckCircle, AlertTriangle, Truck, Printer } from 'lucide-react';
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

const CardDashboard = () => {
  return (
    <div className="min-h-screen bg-[#0a0c10] p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black tracking-tighter text-white">Card Personalization</h1>
          <p className="text-gray-500 uppercase tracking-[0.2em] text-[10px] font-bold mt-1">Inventory & Production Dashboard</p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-cyan-500/10 border border-cyan-500/20 rounded-full">
          <span className="w-2 h-2 rounded-full bg-cyan-500 animate-pulse" />
          <span className="text-cyan-500 text-xs font-bold uppercase tracking-wider">Production Active</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard title="Cards Produced" value="12,847" icon={<CreditCard className="w-5 h-5 text-blue-500" />} trend="+340 today" />
        <KpiCard title="In Inventory" value="5,230" icon={<Package className="w-5 h-5 text-emerald-500" />} trend="32% capacity" />
        <KpiCard title="Pending Print" value="892" icon={<Printer className="w-5 h-5 text-amber-500" />} trend="Est. 3.2 hours" />
        <KpiCard title="In Transit" value="2,100" icon={<Truck className="w-5 h-5 text-violet-500" />} trend="6 shipments" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2"><AlertTriangle className="w-4 h-4 text-amber-500" /> Quality Alerts</h3>
          <div className="space-y-2">
            {[
              { batch: 'BATCH-2026-06-21-A', issue: 'Chip encoding mismatch', cards: 12, severity: 'Critical' },
              { batch: 'BATCH-2026-06-21-B', issue: 'Lamination bubble', cards: 3, severity: 'Minor' },
              { batch: 'BATCH-2026-06-20-C', issue: 'Hologram offset', cards: 7, severity: 'Major' },
            ].map((a) => (
              <div key={a.batch} className="flex items-center justify-between p-3 bg-white/5 rounded-xl">
                <div>
                  <div className="text-sm font-medium text-white">{a.issue}</div>
                  <div className="text-xs text-gray-500">{a.batch} - {a.cards} cards affected</div>
                </div>
                <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${a.severity === 'Critical' ? 'bg-red-500/10 text-red-500' : a.severity === 'Major' ? 'bg-amber-500/10 text-amber-500' : 'bg-gray-500/10 text-gray-400'}`}>{a.severity}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-2xl p-6 shadow-xl">
          <h3 className="text-lg font-bold text-white mb-4">Production Metrics</h3>
          <div className="space-y-4">
            {[
              { label: 'Personalization Rate', val: '98.2%', color: 'bg-emerald-500' },
              { label: 'Quality Pass Rate', val: '99.1%', color: 'bg-blue-500' },
              { label: 'Printer Utilization', val: '76%', color: 'bg-violet-500' },
              { label: 'Chip Inventory', val: '42%', color: 'bg-cyan-500' },
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

export default CardDashboard;
