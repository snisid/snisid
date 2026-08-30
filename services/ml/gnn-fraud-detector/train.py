import torch
import torch.nn.functional as F
from torch_geometric.data import Data
from torch_geometric.nn import GCNConv
import numpy as np


class FraudGNN(torch.nn.Module):
    def __init__(self, in_channels=4, hidden_channels=64, out_channels=2):
        super().__init__()
        self.conv1 = GCNConv(in_channels, hidden_channels)
        self.conv2 = GCNConv(hidden_channels, hidden_channels)
        self.conv3 = GCNConv(hidden_channels, out_channels)

    def forward(self, x, edge_index):
        x = F.relu(self.conv1(x, edge_index))
        x = F.dropout(x, p=0.5, training=self.training)
        x = F.relu(self.conv2(x, edge_index))
        x = self.conv3(x, edge_index)
        return F.log_softmax(x, dim=1)


def generate_synthetic_fraud_graph(n_nodes=10000, fraud_rate=0.05):
    x = torch.randn(n_nodes, 4)
    n_fraud = int(n_nodes * fraud_rate)
    y = torch.zeros(n_nodes, dtype=torch.long)
    y[:n_fraud] = 1

    edges = []
    for i in range(0, n_fraud, 5):
        for j in range(i, min(i + 5, n_fraud)):
            for k in range(j + 1, min(i + 5, n_fraud)):
                edges.append([j, k])
                edges.append([k, j])

    if len(edges) == 0:
        edges = [[0, 1], [1, 0]]

    edge_index = torch.tensor(edges, dtype=torch.long).t().contiguous()

    n_train = int(n_nodes * 0.6)
    n_val = int(n_nodes * 0.2)
    train_mask = torch.zeros(n_nodes, dtype=torch.bool)
    val_mask = torch.zeros(n_nodes, dtype=torch.bool)
    test_mask = torch.zeros(n_nodes, dtype=torch.bool)
    train_mask[:n_train] = True
    val_mask[n_train:n_train + n_val] = True
    test_mask[n_train + n_val:] = True

    return Data(x=x, edge_index=edge_index, y=y,
                train_mask=train_mask, val_mask=val_mask, test_mask=test_mask)


def train_epoch(model, data, optimizer):
    model.train()
    optimizer.zero_grad()
    out = model(data.x, data.edge_index)
    loss = F.nll_loss(out[data.train_mask], data.y[data.train_mask])
    loss.backward()
    optimizer.step()
    return loss.item()


def evaluate(model, data):
    model.eval()
    with torch.no_grad():
        out = model(data.x, data.edge_index)
        pred = out.argmax(dim=1)
        acc = (pred[data.test_mask] == data.y[data.test_mask]).sum().float() / data.test_mask.sum()
    return acc.item()


if __name__ == "__main__":
    model = FraudGNN(in_channels=4, hidden_channels=64, out_channels=2)
    optimizer = torch.optim.Adam(model.parameters(), lr=0.01, weight_decay=5e-4)

    data = generate_synthetic_fraud_graph(n_nodes=10000, fraud_rate=0.05)

    for epoch in range(200):
        loss = train_epoch(model, data, optimizer)
        acc = evaluate(model, data)
        if epoch % 10 == 0:
            print(f"Epoch {epoch:03d} | Loss: {loss:.4f} | Acc: {acc:.4f}")

    torch.save(model.state_dict(), "models/gnn_fraud_v1.pt")
