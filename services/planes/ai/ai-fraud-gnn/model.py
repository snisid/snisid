import torch
import torch.nn.functional as F
from torch_geometric.nn import GCNConv, GATConv

class FraudGNN(torch.nn.Module):
    def __init__(self, num_node_features, hidden_channels=64):
        super(FraudGNN, self).__init__()
        # Using GAT (Graph Attention Network) for better feature weighting in fraud detection
        self.conv1 = GATConv(num_node_features, hidden_channels, heads=4)
        self.conv2 = GATConv(hidden_channels * 4, hidden_channels, heads=1)
        self.fc = torch.nn.Linear(hidden_channels, 2) # Binary classification: Fraud vs Normal

    def forward(self, x, edge_index):
        # x: Node feature matrix [num_nodes, num_node_features]
        # edge_index: Graph connectivity matrix [2, num_edges]
        
        x = self.conv1(x, edge_index)
        x = F.elu(x)
        x = F.dropout(x, p=0.4, training=self.training)
        
        x = self.conv2(x, edge_index)
        x = F.elu(x)
        
        x = self.fc(x)
        return F.log_softmax(x, dim=1)
