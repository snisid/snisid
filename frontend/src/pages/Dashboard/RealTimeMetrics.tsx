import React from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { useTranslation } from 'react-i18next';
import { useWebSocketConnection } from '../../services/websocket';
import { Activity, Users, ShieldAlert, CheckCircle } from 'lucide-react';

// Mock historical data for the chart since WS only gives current snapshot in this example
const generateMockHistory = () => {
  const data = [];
  for (let i = 0; i < 24; i++) {
    data.push({
      time: `${i}:00`,
      requests: Math.floor(Math.random() * 5000) + 1000,
      threats: Math.floor(Math.random() * 50)
    });
  }
  return data;
};

const mockChartData = generateMockHistory();

export const RealTimeMetrics: React.FC = () => {
  const { t } = useTranslation();
  const { isConnected, metrics } = useWebSocketConnection();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-white">{t('nav.dashboard')}</h2>
        <div className={`px-3 py-1 rounded-full text-xs font-medium flex items-center gap-2 border ${isConnected ? 'bg-green-500/10 text-green-500 border-green-500/20' : 'bg-red-500/10 text-red-500 border-red-500/20'}`}>
          <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`}></span>
          {isConnected ? 'Live Sync' : 'Disconnected'}
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard 
          title={t('dashboard.totalIdentities')} 
          value={metrics.totalIdentities.toLocaleString()} 
          icon={<Users size={24} className="text-blue-500" />} 
          trend="+12% cette semaine"
        />
        <KpiCard 
          title={t('dashboard.verifiedIdentities')} 
          value={metrics.verifiedIdentities.toLocaleString()} 
          icon={<CheckCircle size={24} className="text-green-500" />} 
          trend="98.5% taux de succès"
        />
        <KpiCard 
          title={t('dashboard.suspiciousActivities')} 
          value={metrics.suspiciousActivities.toString()} 
          icon={<Activity size={24} className="text-orange-500" />} 
          trend="Score risque moyen: 0.12"
        />
        <KpiCard 
          title={t('dashboard.recentAlerts')} 
          value={metrics.activeAlerts.toString()} 
          icon={<ShieldAlert size={24} className="text-red-500" />} 
          trend="-3% depuis hier"
        />
      </div>

      {/* Charts */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 shadow-xl">
        <h3 className="text-lg font-medium text-slate-200 mb-6">Trafic & Détections (Temps Réel)</h3>
        <div className="h-80 w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={mockChartData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
              <defs>
                <linearGradient id="colorRequests" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorThreats" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ef4444" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#ef4444" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <XAxis dataKey="time" stroke="#475569" fontSize={12} tickLine={false} axisLine={false} />
              <YAxis stroke="#475569" fontSize={12} tickLine={false} axisLine={false} tickFormatter={(value) => `${value}`} />
              <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
              <Tooltip 
                contentStyle={{ backgroundColor: '#0f172a', borderColor: '#1e293b', borderRadius: '0.5rem', color: '#f1f5f9' }}
                itemStyle={{ color: '#f1f5f9' }}
              />
              <Area type="monotone" dataKey="requests" stroke="#3b82f6" strokeWidth={2} fillOpacity={1} fill="url(#colorRequests)" />
              <Area type="monotone" dataKey="threats" stroke="#ef4444" strokeWidth={2} fillOpacity={1} fill="url(#colorThreats)" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
};

const KpiCard: React.FC<{ title: string; value: string; icon: React.ReactNode; trend: string }> = ({ title, value, icon, trend }) => (
  <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 shadow-lg flex flex-col gap-4">
    <div className="flex items-start justify-between">
      <h3 className="text-slate-400 font-medium text-sm">{title}</h3>
      <div className="p-2 bg-slate-800/50 rounded-lg">{icon}</div>
    </div>
    <div>
      <div className="text-3xl font-bold text-white mb-1">{value}</div>
      <div className="text-xs text-slate-500">{trend}</div>
    </div>
  </div>
);
