import React, { useState } from 'react';
import { Shield, Lock, User, ArrowRight, Fingerprint } from 'lucide-react';
import { motion } from 'framer-motion';

const Login = () => {
  const [isBiometric, setIsBiometric] = useState(false);

  return (
    <div className="min-h-screen bg-[#0a0c10] flex items-center justify-center p-4 selection:bg-blue-500/30">
      {/* Background Glow */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-blue-600/10 blur-[120px] rounded-full" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-indigo-600/10 blur-[120px] rounded-full" />
      </div>

      <motion.div 
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="w-full max-w-md relative z-10"
      >
        <div className="text-center mb-10">
          <div className="inline-flex w-16 h-16 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-2xl items-center justify-center shadow-2xl shadow-blue-900/40 mb-6">
            <Shield className="w-10 h-10 text-white" />
          </div>
          <h1 className="text-4xl font-black tracking-tighter text-white">SNISID</h1>
          <p className="text-gray-500 uppercase tracking-[0.3em] text-[10px] font-bold mt-2">Secure National Identity System</p>
        </div>

        <div className="bg-[#0f1218] border border-white/5 rounded-[2.5rem] p-10 shadow-2xl shadow-black/50">
          <div className="flex bg-black/20 p-1 rounded-2xl mb-8">
            <button 
              onClick={() => setIsBiometric(false)}
              className={`flex-1 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all ${!isBiometric ? 'bg-white/5 text-white' : 'text-gray-500'}`}
            >
              Standard
            </button>
            <button 
              onClick={() => setIsBiometric(true)}
              className={`flex-1 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all ${isBiometric ? 'bg-white/5 text-white' : 'text-gray-500'}`}
            >
              Biometric
            </button>
          </div>

          {!isBiometric ? (
            <div className="space-y-6">
              <div className="space-y-2">
                <label className="text-[10px] uppercase tracking-widest font-bold text-gray-500 ml-4">Credential ID</label>
                <div className="relative">
                  <User className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-600" />
                  <input 
                    type="text" 
                    placeholder="Enter agency ID" 
                    className="w-full bg-black/20 border border-white/5 rounded-2xl py-4 pl-12 pr-4 focus:outline-none focus:border-blue-500/50 transition-all text-sm"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-[10px] uppercase tracking-widest font-bold text-gray-500 ml-4">Access Key</label>
                <div className="relative">
                  <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-600" />
                  <input 
                    type="password" 
                    placeholder="••••••••" 
                    className="w-full bg-black/20 border border-white/5 rounded-2xl py-4 pl-12 pr-4 focus:outline-none focus:border-blue-500/50 transition-all text-sm"
                  />
                </div>
              </div>

              <button className="w-full bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 py-4 rounded-2xl font-bold text-sm shadow-lg shadow-blue-900/20 flex items-center justify-center gap-2 group transition-all">
                Authorize Access
                <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </button>
            </div>
          ) : (
            <div className="py-8 text-center space-y-6">
              <motion.div 
                animate={{ scale: [1, 1.05, 1] }}
                transition={{ repeat: Infinity, duration: 2 }}
                className="inline-flex w-24 h-24 bg-blue-500/10 rounded-full items-center justify-center border-2 border-blue-500/20"
              >
                <Fingerprint className="w-12 h-12 text-blue-500" />
              </motion.div>
              <div className="space-y-2">
                <h3 className="font-bold text-white">Biometric Scan Required</h3>
                <p className="text-sm text-gray-500">Please place your finger on the scanner or look at the camera for 3D facial verification.</p>
              </div>
              <button className="w-full bg-white/5 border border-white/5 hover:bg-white/10 py-4 rounded-2xl font-bold text-sm transition-all">
                Cancel
              </button>
            </div>
          )}
        </div>

        <p className="text-center mt-8 text-[10px] text-gray-600 uppercase tracking-widest font-bold">
          © 2026 SNISID Platform • Restricted Government Access
        </p>
      </motion.div>
    </div>
  );
};

export default Login;
