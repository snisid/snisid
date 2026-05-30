import { useMemo } from 'react';
import ForceGraph3D from 'react-force-graph-3d';
import { Box, Cpu, Database, Activity, Play, RotateCcw } from 'lucide-react';

const DigitalTwin = () => {
  // Mock infrastructure graph
  const data = useMemo(() => {
    const nodes = [
      { id: 'Node-1', type: 'host', label: 'EKS-Worker-01', color: '#1f2937' },
      { id: 'Node-2', type: 'host', label: 'EKS-Worker-02', color: '#1f2937' },
      { id: 'Pod-1', type: 'service', label: 'identity-api', color: '#3b82f6' },
      { id: 'Pod-2', type: 'service', label: 'fraud-engine', color: '#3b82f6' },
      { id: 'DB-1', type: 'db', label: 'PostgreSQL', color: '#10b981' },
      { id: 'DB-2', type: 'db', label: 'Neo4j', color: '#10b981' },
    ];
    const links = [
      { source: 'Pod-1', target: 'Node-1' },
      { source: 'Pod-2', target: 'Node-1' },
      { source: 'Pod-1', target: 'DB-1' },
      { source: 'Pod-2', target: 'DB-2' },
      { source: 'Node-1', target: 'Node-2' },
    ];
    return { nodes, links };
  }, []);

  return (
    <div className="p-8 space-y-8 h-screen flex flex-col">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-black tracking-tighter">Digital Twin Infrastructure</h2>
          <p className="text-gray-500 mt-1">Real-time infrastructure mirroring and predictive failure simulation</p>
        </div>
        <div className="flex gap-3">
          <button className="flex items-center gap-2 px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-sm hover:bg-white/10 transition-all">
            <RotateCcw className="w-4 h-4" /> Reset Twin
          </button>
          <button className="flex items-center gap-2 px-4 py-2 bg-indigo-600 rounded-xl text-sm font-medium hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-900/20">
            <Play className="w-4 h-4" /> Run Simulation
          </button>
        </div>
      </header>

      <div className="flex-1 min-h-0 grid grid-cols-4 gap-8">
        {/* 3D Graph Visualizer */}
        <div className="col-span-3 bg-[#0a0c10] rounded-[2.5rem] border border-white/5 overflow-hidden relative shadow-inner shadow-black/50">
          <div className="absolute top-8 left-8 z-10 space-y-4">
             <div className="flex gap-4">
              <Legend color="bg-gray-800" label="EKS Nodes" />
              <Legend color="bg-blue-500" label="Microservices" />
              <Legend color="bg-emerald-500" label="Data Tier" />
            </div>
            <div className="p-4 bg-black/40 backdrop-blur-xl border border-white/5 rounded-2xl">
              <div className="flex items-center gap-2 text-[10px] font-bold text-emerald-400 uppercase tracking-widest mb-2">
                <Activity className="w-3 h-3" /> Real-time Sync
              </div>
              <div className="text-2xl font-black tabular-nums tracking-tighter">14ms <span className="text-sm font-normal text-gray-500">Latency</span></div>
            </div>
          </div>
          
          <ForceGraph3D
            graphData={data}
            backgroundColor="#0a0c10"
            nodeLabel="label"
            nodeColor="color"
            linkOpacity={0.2}
            linkWidth={1}
            showNavInfo={false}
          />
        </div>

        {/* Control Panel */}
        <div className="space-y-6 overflow-y-auto pr-2 custom-scrollbar">
          <h3 className="text-lg font-bold">Simulation Scenarios</h3>
          <ScenarioCard 
            title="Node Outage" 
            desc="Simulate the failure of a primary EKS worker node." 
            icon={Cpu}
          />
          <ScenarioCard 
            title="Kafka Saturation" 
            desc="Predict impact of 100x event volume spike." 
            icon={Box}
          />
          <ScenarioCard 
            title="DB Partition" 
            desc="Analyze system behavior during network partition." 
            icon={Database}
          />
        </div>
      </div>
    </div>
  );
};

const Legend = ({ color, label }: { color: string, label: string }) => (
  <div className="flex items-center gap-2 px-3 py-1.5 bg-[#0f1218]/80 backdrop-blur-md border border-white/5 rounded-full">
    <div className={`w-2 h-2 rounded-full ${color}`} />
    <span className="text-[10px] font-bold uppercase tracking-tighter text-gray-400">{label}</span>
  </div>
);

const ScenarioCard = ({ title, desc, icon: Icon }: { title: string, desc: string, icon: any }) => (
  <div className="p-6 bg-[#0f1218] border border-white/5 rounded-3xl hover:bg-white/[0.02] transition-all group cursor-pointer">
    <div className="flex justify-between items-start mb-4">
      <div className="p-3 bg-white/5 rounded-2xl group-hover:bg-indigo-500/10 group-hover:text-indigo-400 transition-all">
        <Icon className="w-5 h-5" />
      </div>
      <button className="text-[10px] font-bold text-gray-500 uppercase tracking-widest hover:text-white">Run</button>
    </div>
    <h4 className="font-bold mb-1">{title}</h4>
    <p className="text-xs text-gray-500 leading-relaxed">{desc}</p>
  </div>
);

export default DigitalTwin;
