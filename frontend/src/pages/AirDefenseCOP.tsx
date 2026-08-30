import React from 'react';
import { motion } from 'framer-motion';
import { Radar, Crosshair, Ban, Clock, Shield, AlertTriangle, Plane, Navigation } from 'lucide-react';

const statCards = [
  { label: 'Radar Contacts', value: '156', color: 'text-green-400', icon: Radar },
  { label: 'Active Threats', value: '4', color: 'text-red-400', icon: Crosshair },
  { label: 'No-Fly Violations', value: '23', color: 'text-orange-400', icon: Ban },
  { label: 'Incidents (24h)', value: '8', color: 'text-yellow-400', icon: Clock },
];

const contactTypes = [
  { label: 'Commercial', value: 78, color: 'bg-blue-500' },
  { label: 'Military', value: 42, color: 'bg-green-500' },
  { label: 'Private', value: 23, color: 'bg-yellow-500' },
  { label: 'Unknown', value: 13, color: 'bg-red-500' },
];

const totalContacts = contactTypes.reduce((s, c) => s + c.value, 0);

const activeThreats = [
  { id: 'T-01', type: 'Unauthorized UAV', altitude: '3,400 ft', heading: 'NE', speed: '45 kts', severity: 'high' },
  { id: 'T-02', type: 'Squawk 7500', altitude: '12,000 ft', heading: 'S', speed: '320 kts', severity: 'critical' },
  { id: 'T-03', type: 'Unknown Track', altitude: '8,200 ft', heading: 'W', speed: '180 kts', severity: 'medium' },
  { id: 'T-04', type: 'Airspace Violation', altitude: '2,100 ft', heading: 'N', speed: '90 kts', severity: 'high' },
];

const incidents = [
  { time: '02:34', type: 'Near Miss', location: 'Sector 7', status: 'Investigating' },
  { time: '04:15', type: 'Drone Incursion', location: 'Air Base Alpha', status: 'Resolved' },
  { time: '07:42', type: 'Radio Failure', location: 'Flight Corridor B', status: 'Monitoring' },
  { time: '11:03', type: 'Bird Strike', location: 'Runway 2R', status: 'Resolved' },
  { time: '14:28', type: 'Unauthorized Landing', location: 'Private Strip #4', status: 'Pending' },
];

const contacts = [
  { id: 'AC-101', type: 'Boeing 737', alt: '35,000 ft', speed: '450 kts', heading: '090°', friendly: true },
  { id: 'AC-102', type: 'C-130', alt: '24,000 ft', speed: '280 kts', heading: '270°', friendly: true },
  { id: 'AC-103', type: 'Unknown', alt: '8,000 ft', speed: '120 kts', heading: '045°', friendly: false },
  { id: 'AC-104', type: 'Gulfstream V', alt: '41,000 ft', speed: '490 kts', heading: '180°', friendly: true },
];

const AirDefenseCOP = () => {
  return (
    <div className="p-8 space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Air Defense — Common Operating Picture</h2>
          <p className="text-gray-500 mt-1">Real-time air surveillance and threat management</p>
        </div>
      </header>

      <div className="grid grid-cols-4 gap-6">
        {statCards.map((card) => (
          <StatusCard key={card.label} {...card} />
        ))}
      </div>

      <div className="grid grid-cols-3 gap-6">
        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-8">
            <Radar className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Radar Contacts by Type</h3>
          </div>
          <div className="relative w-48 h-48 mx-auto mb-6">
            <svg viewBox="0 0 100 100" className="w-full h-full -rotate-90">
              {contactTypes.reduce((acc, c, i) => {
                const prevAngle = acc.length > 0 ? acc[acc.length - 1].endAngle : 0;
                const angle = (c.value / totalContacts) * 360;
                const endAngle = prevAngle + angle;
                const x1 = 50 + 40 * Math.cos((prevAngle * Math.PI) / 180);
                const y1 = 50 + 40 * Math.sin((prevAngle * Math.PI) / 180);
                const x2 = 50 + 40 * Math.cos((endAngle * Math.PI) / 180);
                const y2 = 50 + 40 * Math.sin((endAngle * Math.PI) / 180);
                const largeArc = angle > 180 ? 1 : 0;
                acc.push({ prevAngle, endAngle, path: `M 50 50 L ${x1} ${y1} A 40 40 0 ${largeArc} 1 ${x2} ${y2} Z`, color: c.color });
                return acc;
              }, [] as { prevAngle: number; endAngle: number; path: string; color: string }[]).map((slice, i) => (
                <motion.path
                  key={i}
                  d={slice.path}
                  fill={slice.color.replace('bg-', '')}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 0.8 }}
                  transition={{ delay: i * 0.1 }}
                />
              ))}
            </svg>
          </div>
          <div className="space-y-2">
            {contactTypes.map((c) => (
              <div key={c.label} className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-2">
                  <div className={`w-3 h-3 rounded-full ${c.color}`} />
                  <span className="text-gray-400">{c.label}</span>
                </div>
                <span className="font-bold">{c.value}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Crosshair className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Active Threats</h3>
          </div>
          <div className="space-y-3">
            {activeThreats.map((t, i) => (
              <motion.div
                key={t.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.05 }}
                className={`p-4 rounded-2xl border ${
                  t.severity === 'critical' ? 'border-red-500/20 bg-red-500/5' :
                  t.severity === 'high' ? 'border-orange-500/20 bg-orange-500/5' :
                  'border-yellow-500/20 bg-yellow-500/5'
                }`}
              >
                <div className="flex justify-between items-start mb-2">
                  <span className="font-bold text-sm">{t.type}</span>
                  <span className={`text-[10px] font-black uppercase tracking-widest ${
                    t.severity === 'critical' ? 'text-red-500' :
                    t.severity === 'high' ? 'text-orange-500' : 'text-yellow-500'
                  }`}>{t.severity}</span>
                </div>
                <div className="flex gap-4 text-xs text-gray-400">
                  <span>ALT {t.altitude}</span>
                  <span>HDG {t.heading}</span>
                  <span>{t.speed}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
          <div className="flex items-center gap-3 mb-6">
            <Clock className="w-5 h-5 text-gray-400" />
            <h3 className="text-xl font-bold">Recent Incidents</h3>
          </div>
          <div className="space-y-3">
            {incidents.map((inc, i) => (
              <motion.div
                key={`${inc.time}-${inc.type}`}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/5"
              >
                <div>
                  <div className="text-sm font-bold">{inc.type}</div>
                  <div className="text-xs text-gray-500">{inc.location}</div>
                </div>
                <div className="text-right">
                  <div className="text-xs font-mono text-gray-400">{inc.time}</div>
                  <span className={`text-[10px] font-bold uppercase ${
                    inc.status === 'Resolved' ? 'text-green-400' :
                    inc.status === 'Investigating' ? 'text-orange-400' : 'text-yellow-400'
                  }`}>{inc.status}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      <div className="bg-[#0f1218] rounded-[2rem] border border-white/5 p-8">
        <div className="flex items-center gap-3 mb-6">
          <Navigation className="w-5 h-5 text-gray-400" />
          <h3 className="text-xl font-bold">Current Air Picture</h3>
        </div>
        <div className="grid grid-cols-4 gap-4">
          {contacts.map((ac, i) => (
            <motion.div
              key={ac.id}
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: i * 0.05 }}
              className="p-5 rounded-2xl bg-white/[0.02] border border-white/5"
            >
              <div className="flex items-center gap-3 mb-3">
                <Plane className={`w-5 h-5 ${ac.friendly ? 'text-blue-400' : 'text-red-400'}`} />
                <span className="font-bold text-sm">{ac.id}</span>
              </div>
              <div className="text-lg font-bold mb-1">{ac.type}</div>
              <div className="grid grid-cols-2 gap-1 text-xs text-gray-500">
                <span>ALT {ac.alt}</span>
                <span>{ac.speed}</span>
                <span>HDG {ac.heading}</span>
              </div>
            </motion.div>
          ))}
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

export default AirDefenseCOP;
