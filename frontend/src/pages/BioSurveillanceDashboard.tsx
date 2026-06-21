import React from 'react';
import { motion } from 'framer-motion';
import { Activity, Vaccine, Pill, AlertTriangle, Heart, Map, Thermometer, Droplets } from 'lucide-react';

const statCards = [
  { label: 'Active Alerts', value: '9', color: 'text-red-400', icon: Activity },
  { label: 'Vaccination Coverage', value: '74%', color: 'text-green-400', icon: Vaccine },
  { label: 'Facilities Stocked', value: '86%', color: 'text-blue-400', icon: Pill },
  { label: 'Regions Monitored', value: '42', color: 'text-cyan-400', icon: Map },
];

const alerts = [
  { disease: 'Cholera', region: 'Port-au-Prince', cases: 47, alert: 'red' },
  { disease: 'Dengue', region: 'Cap-Haïtien', cases: 124, alert: 'orange' },
  { disease: 'Malaria', region: 'Rural South', cases: 89, alert: 'yellow' },
  { disease: 'Tuberculosis', region: 'Urban Center', cases: 23, alert: 'orange' },
  { disease: 'COVID-19', region: 'National', cases: 312, alert: 'yellow' },
];

const campaigns = [
  { name: 'Polio Eradication', population: '2.1M', vaccinated: '1.8M', completion: 86 },
  { name: 'Measles Drive', population: '1.5M', vaccinated: '1.1M', completion: 73 },
  { name: 'Yellow Fever', population: '980K', vaccinated: '620K', completion: 63 },
];

const facilityStock = [
  { name: 'Hôpital Général', status: 'green', beds: 89, oxygen: 92, meds: 78 },
  { name: 'Clinic Sud-Est', status: 'yellow', beds: 45, oxygen: 32, meds: 56 },
  { name: 'Field Hospital A', status: 'red', beds: 12, oxygen: 8, meds: 22 },
  { name: 'Regional Med Ctr', status: 'green', beds: 120, oxygen: 105, meds: 95 },
];

const BioSurveillanceDashboard = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Health Security Center</h2>
          <p className="text-gray-500 mt-1">Biosurveillance and public health monitoring dashboard</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {statCards.map((card) => (
          <StatusCard key={card.label} {...card} />
        ))}
      </div>

      <div className="grid grid-cols-3 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <AlertTriangle className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Active Disease Alerts</h3>
          </div>
          <div className="space-y-3">
            {alerts.map((a, i) => (
              <motion.div
                key={`${a.disease}-${a.region}`}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${
                  a.alert === 'red' ? 'border-red-500/20 bg-red-500/5' :
                  a.alert === 'orange' ? 'border-orange-500/20 bg-orange-500/5' :
                  'border-yellow-500/20 bg-yellow-500/5'
                }`}
              >
                <div className="flex justify-between items-center mb-1">
                  <span className="font-bold">{a.disease}</span>
                  <div className={`w-2 h-2 rounded-full ${
                    a.alert === 'red' ? 'bg-red-500' :
                    a.alert === 'orange' ? 'bg-orange-500' : 'bg-yellow-500'
                  }`} />
                </div>
                <div className="text-sm text-gray-400">{a.region}</div>
                <div className="text-xs text-gray-500 mt-1">{a.cases} confirmed cases</div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Vaccine className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Vaccination Campaigns</h3>
          </div>
          <div className="space-y-5">
            {campaigns.map((c, i) => (
              <motion.div
                key={c.name}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
              >
                <div className="flex justify-between text-sm mb-2">
                  <span className="font-bold">{c.name}</span>
                  <span className="text-green-400 font-bold">{c.completion}%</span>
                </div>
                <div className="h-2 bg-white/5 rounded-full overflow-hidden mb-1">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${c.completion}%` }}
                    transition={{ duration: 0.8, delay: i * 0.1 }}
                    className="h-full rounded-full bg-green-500"
                  />
                </div>
                <div className="text-xs text-gray-500">
                  {c.vaccinated} / {c.population} population
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Pill className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Facility Stock Status</h3>
          </div>
          <div className="space-y-3">
            {facilityStock.map((f, i) => (
              <motion.div
                key={f.name}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${
                  f.status === 'green' ? 'border-green-500/20 bg-green-500/5' :
                  f.status === 'yellow' ? 'border-yellow-500/20 bg-yellow-500/5' :
                  'border-red-500/20 bg-red-500/5'
                }`}
              >
                <div className="flex justify-between items-center mb-2">
                  <span className="font-bold text-sm">{f.name}</span>
                  <div className={`w-2 h-2 rounded-full ${
                    f.status === 'green' ? 'bg-green-500' :
                    f.status === 'yellow' ? 'bg-yellow-500' : 'bg-red-500'
                  }`} />
                </div>
                <div className="grid grid-cols-3 gap-2 text-xs text-gray-400">
                  <span>Beds: {f.beds}%</span>
                  <span>O₂: {f.oxygen}%</span>
                  <span>Meds: {f.meds}%</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
        <div className="flex items-center gap-3 mb-6">
          <Heart className="w-5 h-5 text-gray-400" />
          <h3 className="text-xl font-bold">National Health Picture</h3>
        </div>
        <div className="grid grid-cols-4 gap-6">
          <div className="text-center p-4 rounded-2xl bg-blue-500/5 border border-blue-500/20">
            <Thermometer className="w-6 h-6 text-blue-400 mx-auto mb-2" />
            <div className="text-2xl font-black">3,847</div>
            <div className="text-xs text-gray-500">Total Cases (30d)</div>
          </div>
          <div className="text-center p-4 rounded-2xl bg-green-500/5 border border-green-500/20">
            <Heart className="w-6 h-6 text-green-400 mx-auto mb-2" />
            <div className="text-2xl font-black">2,104</div>
            <div className="text-xs text-gray-500">Recovered</div>
          </div>
          <div className="text-center p-4 rounded-2xl bg-red-500/5 border border-red-500/20">
            <Droplets className="w-6 h-6 text-red-400 mx-auto mb-2" />
            <div className="text-2xl font-black">47</div>
            <div className="text-xs text-gray-500">Active Outbreaks</div>
          </div>
          <div className="text-center p-4 rounded-2xl bg-yellow-500/5 border border-yellow-500/20">
            <Activity className="w-6 h-6 text-yellow-400 mx-auto mb-2" />
            <div className="text-2xl font-black">12</div>
            <div className="text-xs text-gray-500">Critical Facilities</div>
          </div>
        </div>
      </div>
    </div>
  );
};

const StatusCard = ({ label, value, color, icon: Icon }: { label: string; value: string; color: string; icon: React.ComponentType<any> }) => (
  <motion.div
    initial={{ opacity: 0, y: 20 }}
    animate={{ opacity: 1, y: 0 }}
    className="p-6 bg-[#0f1218] rounded-3xl border border-white/5"
  >
    <div className="flex justify-between items-start mb-2">
      <p className="text-xs font-bold uppercase tracking-widest text-gray-500">{label}</p>
      <Icon className={`w-5 h-5 ${color}`} />
    </div>
    <h3 className="text-3xl font-black tracking-tighter">{value}</h3>
  </motion.div>
);

export default BioSurveillanceDashboard;
