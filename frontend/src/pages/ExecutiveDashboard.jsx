import React, { useState, useEffect } from 'react';

const ExecutiveDashboard = () => {
  const [riskData, setRiskData] = useState([]);
  const [threats, setThreats] = useState([]);

  useEffect(() => {
    // WebSocket connection simulation
    const socket = new WebSocket('ws://api-gateway.nexus.svc/ws/executive');
    
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'THREAT_ALERT') {
        setThreats((prev) => [data.payload, ...prev].slice(0, 5));
      }
      if (data.type === 'RISK_UPDATE') {
        setRiskData(data.payload);
      }
    };

    return () => socket.close();
  }, []);

  return (
    <div className="command-center bg-gray-900 text-white min-h-screen p-8 font-mono">
      <header className="border-b border-cyan-500 pb-4 mb-8 flex justify-between">
        <h1 className="text-3xl font-bold tracking-tighter text-cyan-400">SNISID VERA :: COMMAND CENTER</h1>
        <div className="status-indicator flex items-center gap-2">
          <span className="animate-pulse w-3 h-3 bg-green-500 rounded-full"></span>
          <span>LIVE TELEMETRY: ACTIVE</span>
        </div>
      </header>

      <div className="grid grid-cols-12 gap-6">
        {/* Global Risk Map placeholder */}
        <section className="col-span-8 bg-black border border-gray-800 p-6 rounded-lg">
          <h2 className="text-xl mb-4 text-cyan-500">GLOBAL RISK INTENSITY MAP</h2>
          <div className="h-96 bg-gray-950 flex items-center justify-center border border-gray-900">
            [D3.js Heatmap Rendering Engine]
          </div>
        </section>

        {/* Threat Panel */}
        <section className="col-span-4 bg-black border border-gray-800 p-6 rounded-lg">
          <h2 className="text-xl mb-4 text-red-500">ACTIVE THREAT FEED</h2>
          <div className="space-y-4">
            {threats.map((t, i) => (
              <div key={i} className="border-l-2 border-red-500 pl-4 py-2 bg-gray-900">
                <p className="text-xs text-gray-500">{t.timestamp}</p>
                <p className="text-sm">{t.description}</p>
                <p className="text-xs font-bold text-red-400">RISK SCORE: {t.score}</p>
              </div>
            ))}
          </div>
        </section>
      </div>

      <footer className="mt-8 grid grid-cols-4 gap-4">
        <button className="bg-red-900 hover:bg-red-800 border border-red-500 p-4 font-bold">EMERGENCY LOCKDOWN</button>
        <button className="bg-cyan-900 hover:bg-cyan-800 border border-cyan-500 p-4 font-bold">RE-CALIBRATE AI</button>
        <button className="bg-gray-800 hover:bg-gray-700 p-4">VIEW AUDIT LEDGER</button>
        <button className="bg-gray-800 hover:bg-gray-700 p-4">EXPORT REPORT</button>
      </footer>
    </div>
  );
};

export default ExecutiveDashboard;
