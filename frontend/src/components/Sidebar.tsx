import React from 'react';
import { NavLink } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Search, 
  ShieldAlert, 
  Activity, 
  Users, 
  Settings, 
  LogOut, 
  Shield, 
  Fingerprint,
  Map
} from 'lucide-react';
import { useStore } from '../store/useStore';

const Sidebar = () => {
  const { user } = useStore();

  return (
    <aside className="fixed left-0 top-0 h-full w-64 bg-[#0f1218] border-r border-white/5 z-50">
      <div className="p-6 flex items-center gap-3">
        <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-xl flex items-center justify-center shadow-lg shadow-blue-900/20">
          <Shield className="w-6 h-6 text-white" />
        </div>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-white">SNISID</h1>
          <p className="text-[10px] text-gray-500 uppercase tracking-widest font-semibold">National Intelligence</p>
        </div>
      </div>

      <nav className="mt-8 px-4 space-y-1">
        <NavItem to="/" icon={LayoutDashboard} label="Dashboard" />
        <NavItem to="/search" icon={Search} label="Identity Search" />
        <NavItem to="/fraud" icon={ShieldAlert} label="Fraud Analysis" />
        <NavItem to="/alerts" icon={Activity} label="SOC Alerts" />
        <NavItem to="/biometrics" icon={Fingerprint} label="Biometrics" />
        <NavItem to="/agencies" icon={Users} label="Agencies" />
        <NavItem to="/map" icon={Map} label="Geospatial" />
        <NavItem to="/settings" icon={Settings} label="Settings" />
      </nav>

      <div className="absolute bottom-8 left-0 w-full px-4">
        <div className="p-4 bg-white/5 rounded-2xl border border-white/5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-gray-700 to-gray-600 border border-white/10" />
            <div>
              <p className="text-sm font-medium">{user?.name}</p>
              <p className="text-[10px] text-gray-500">{user?.role}</p>
            </div>
          </div>
          <button className="w-full flex items-center justify-center gap-2 py-2 text-xs text-gray-400 hover:text-white transition-colors">
            <LogOut className="w-3.5 h-3.5" />
            Sign Out
          </button>
        </div>
      </div>
    </aside>
  );
};

const NavItem = ({ to, icon: Icon, label }: { to: string, icon: any, label: string }) => (
  <NavLink 
    to={to} 
    className={({ isActive }) => `
      flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all
      ${isActive 
        ? 'bg-blue-600/10 text-blue-400 border border-blue-500/20' 
        : 'text-gray-500 hover:text-gray-300 hover:bg-white/5 border border-transparent'}
    `}
  >
    <Icon className="w-5 h-5" />
    {label}
  </NavLink>
);

export default Sidebar;
