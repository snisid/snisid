from fastapi import FastAPI, Body
import uvicorn
import torch
import torch.nn.functional as F
from torch_geometric.nn import GCNConv
from torch_geometric.data import Data

app = FastAPI(title="SNISID GNN Fraud Detection Service")

class GCN(torch.nn.Module):
    def __init__(self, num_features):
        super().__init__()
        self.conv1 = GCNConv(num_features, 16)
        self.conv2 = GCNConv(16, 2)

    def forward(self, data):
        x, edge_index = data.x, data.edge_index
        x = self.conv1(x, edge_index)
        x = F.relu(x)
        x = self.conv2(x, edge_index)
        return F.log_softmax(x, dim=1)

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/score")
async def score(data: dict = Body(...)):
    # Simple GNN scoring logic
    # In reality, this would fetch from Neo4j or take a subgraph
    return {
        "fraud_score": 0.12,
        "cluster_id": "C-99",
        "risk_level": "low"
    }

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8002)
