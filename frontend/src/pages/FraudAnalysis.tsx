import { useMemo } from 'react';
import ForceGraph2D from 'react-force-graph-2d';
import { Filter, Share2 } from 'lucide-react';

const FraudAnalysis = () => {
  // Mock graph data
  const data = useMemo(() => {
    const nodes = [
      { id: 'I-101', group: 'identity', label: 'Jean Pierre', val: 10, color: '#3b82f6' },
      { id: 'I-102', group: 'identity', label: 'Marie G.', val: 8, color: '#3b82f6' },
      { id: 'I-103', group: 'identity', label: 'Louis X.', val: 8, color: '#3b82f6' },
      { id: 'F-501', group: 'fraud', label: 'SYNTHETIC CLUSTER', val: 20, color: '#ef4444' },
      { id: 'A-201', group: 'address', label: 'Rue Capois 12', val: 5, color: '#10b981' },
      { id: 'P-301', group: 'phone', label: '+509 3722 1100', val: 5, color: '#f59e0b' },
    ];
    const links = [
      { source: 'I-101', target: 'F-501' },
      { source: 'I-102', target: 'F-501' },
      { source: 'I-103', target: 'F-501' },
      { source: 'I-101', target: 'A-201' },
      { source: 'I-102', target: 'A-201' },
      { source: 'I-101', target: 'P-301' },
    ];
    return { nodes, links };
  }, []);

  return (
    <div className="p-8 space-y-8 h-screen flex flex-col">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Fraud Analysis Graph</h2>
          <p className="text-gray-500 mt-1">GNN-based identity cluster mapping and relationship discovery</p>
        </div>
        <div className="flex gap-3">
          <button className="flex items-center gap-2 px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-sm hover:bg-white/10 transition-all">
            <Filter className="w-4 h-4" /> Filter
          </button>
          <button className="flex items-center gap-2 px-4 py-2 bg-blue-600 rounded-xl text-sm font-medium hover:bg-blue-700 transition-all shadow-lg shadow-blue-900/20">
            <Share2 className="w-4 h-4" /> Export Report
          </button>
        </div>
      </header>

      <div className="flex-1 min-h-0 grid grid-cols-4 gap-8">
        {/* Graph Visualizer */}
        <div className="col-span-3 bg-[#0f1218] rounded-[2rem] border border-white/5 overflow-hidden relative">
          <div className="absolute top-6 left-6 z-10 flex gap-4">
            <Legend color="bg-blue-500" label="Identity" />
            <Legend color="bg-red-500" label="Fraud Cluster" />
            <Legend color="bg-emerald-500" label="Shared PII" />
          </div>
          <ForceGraph2D
            graphData={data}
            nodeLabel="label"
            nodeAutoColorBy="group"
            linkDirectionalParticles={2}
            linkDirectionalParticleSpeed={0.01}
            backgroundColor="#0f1218"
            nodeCanvasObject={(node: any, ctx, globalScale) => {
              const label = node.label;
              const fontSize = 12/globalScale;
              ctx.font = `${fontSize}px Inter`;
              ctx.textAlign = 'center';
              ctx.textBaseline = 'middle';
              ctx.fillStyle = node.color;
              ctx.beginPath();
              ctx.arc(node.x, node.y, node.val/2, 0, 2 * Math.PI, false);
              ctx.fill();
              ctx.fillStyle = 'white';
              ctx.fillText(label, node.x, node.y + node.val + 2);
            }}
          />
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6 overflow-y-auto pr-2 custom-scrollbar">
          <div className="p-6 bg-red-500/10 border border-red-500/20 rounded-3xl">
            <div className="flex items-center gap-2 text-red-400 mb-2">
              <ShieldAlert className="w-5 h-5" />
              <h3 className="font-bold">Critical Cluster</h3>
            </div>
            <p className="text-sm text-red-300/80 leading-relaxed">
              Detected 12 overlapping identity vectors sharing the same biometric hash and physical address.
            </p>
            <div className="mt-4 flex items-center justify-between text-xs font-bold text-red-400 uppercase tracking-tighter">
              <span>Risk Score</span>
              <span>98.4%</span>
            </div>
            <div className="w-full bg-red-900/30 h-1.5 rounded-full mt-1.5">
              <div className="bg-red-500 h-full w-[98%] rounded-full shadow-[0_0_10px_rgba(239,68,68,0.5)]" />
            </div>
          </div>

          <div className="p-6 bg-white/5 border border-white/5 rounded-3xl space-y-4">
            <h4 className="font-bold text-sm text-gray-400 uppercase tracking-widest">Node Properties</h4>
            <Property label="Type" value="Fraud Cluster" />
            <Property label="Node ID" value="F-501" />
            <Property label="Connections" value="48 nodes" />
            <Property label="First Seen" value="2024-03-12" />
            <Property label="Confidence" value="High" />
          </div>
        </div>
      </div>
    </div>
  );
};

const Legend = ({ color, label }: { color: string, label: string }) => (
  <div className="flex items-center gap-2 px-3 py-1.5 bg-[#0a0c10]/80 backdrop-blur-md border border-white/5 rounded-full">
    <div className={`w-2 h-2 rounded-full ${color}`} />
    <span className="text-[10px] font-bold uppercase tracking-tighter text-gray-400">{label}</span>
  </div>
);

const Property = ({ label, value }: { label: string, value: string }) => (
  <div className="flex justify-between items-center text-sm">
    <span className="text-gray-500">{label}</span>
    <span className="font-medium">{value}</span>
  </div>
);

const ShieldAlert = ({ className }: { className?: string }) => (
  <svg className={className} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.5 3.8 17 5 19 5a1 1 0 0 1 1 1z"/><path d="m12 8 3 4"/><path d="m15 8-3 4"/><circle cx="12" cy="16" r="1"/></svg>
);

export default FraudAnalysis;
