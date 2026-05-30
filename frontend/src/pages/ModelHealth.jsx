import React, { useEffect, useState } from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";

export default function ModelHealth() {
  const [metrics, setMetrics] = useState({
    accuracy: 0.92,
    drift: 0.04,
    version: "v1.2.stable",
    history: [
      { time: "10:00", accuracy: 0.91, drift: 0.02 },
      { time: "11:00", accuracy: 0.93, drift: 0.03 },
      { time: "12:00", accuracy: 0.92, drift: 0.04 },
    ]
  });

  return (
    <div className="p-6 bg-slate-900 text-white rounded-xl">
      <h2 className="text-2xl font-bold mb-4">ML Intelligence Health</h2>
      <div className="grid grid-cols-3 gap-4 mb-8">
        <div className="bg-slate-800 p-4 rounded-lg border border-emerald-500/30">
          <p className="text-slate-400 text-sm">Accuracy</p>
          <p className="text-3xl font-mono text-emerald-400">{(metrics.accuracy * 100).toFixed(1)}%</p>
        </div>
        <div className="bg-slate-800 p-4 rounded-lg border border-amber-500/30">
          <p className="text-slate-400 text-sm">Drift Score (PSI)</p>
          <p className="text-3xl font-mono text-amber-400">{metrics.drift.toFixed(3)}</p>
        </div>
        <div className="bg-slate-800 p-4 rounded-lg border border-blue-500/30">
          <p className="text-slate-400 text-sm">Active Model</p>
          <p className="text-xl font-mono text-blue-400">{metrics.version}</p>
        </div>
      </div>

      <div className="bg-slate-800 p-6 rounded-lg h-64">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={metrics.history}>
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
            <XAxis dataKey="time" stroke="#94a3b8" />
            <YAxis stroke="#94a3b8" />
            <Tooltip contentStyle={{ backgroundColor: "#1e293b", border: "none" }} />
            <Line type="monotone" dataKey="accuracy" stroke="#10b981" strokeWidth={2} dot={{ fill: "#10b981" }} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
