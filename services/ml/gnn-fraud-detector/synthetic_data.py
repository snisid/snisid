import numpy as np
import torch
from torch_geometric.data import Data


def generate_synthetic_fraud_graph(
    n_nodes: int = 10000, fraud_rate: float = 0.05
) -> Data:
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

    edge_index = (
        torch.tensor(edges, dtype=torch.long).t().contiguous()
        if edges
        else torch.zeros((2, 0), dtype=torch.long)
    )

    n_train = int(n_nodes * 0.6)
    n_val = int(n_nodes * 0.2)
    train_mask = torch.zeros(n_nodes, dtype=torch.bool)
    val_mask = torch.zeros(n_nodes, dtype=torch.bool)
    test_mask = torch.zeros(n_nodes, dtype=torch.bool)
    train_mask[:n_train] = True
    val_mask[n_train : n_train + n_val] = True
    test_mask[n_train + n_val :] = True

    return Data(
        x=x,
        edge_index=edge_index,
        y=y,
        train_mask=train_mask,
        val_mask=val_mask,
        test_mask=test_mask,
    )


if __name__ == "__main__":
    data = generate_synthetic_fraud_graph(n_nodes=1000, fraud_rate=0.05)
    print(f"Nodes: {data.x.shape[0]}")
    print(f"Features: {data.x.shape[1]}")
    print(f"Edges: {data.edge_index.shape[1]}")
    print(f"Fraud rate: {data.y.float().mean():.4f}")
