import random
import json

def generate_fraud_dataset(num_nodes=1000):
    nodes = []
    edges = []
    
    # Generate nodes
    for i in range(num_nodes):
        nodes.append({
            "id": f"ID-{i}",
            "risk_score": random.uniform(0, 0.5),
            "label": 0 # Normal
        })
    
    # Create fraud clusters (Synthetic Identities)
    for c in range(5): # 5 clusters
        cluster_nodes = random.sample(range(num_nodes), 10)
        for i in cluster_nodes:
            nodes[i]["label"] = 1 # Fraud
            nodes[i]["risk_score"] = random.uniform(0.7, 1.0)
            
            # Connect cluster nodes to each other (clique-ish)
            for j in cluster_nodes:
                if i != j:
                    edges.append({"source": f"ID-{i}", "target": f"ID-{j}"})

    with open("data/fraud_nodes.json", "w") as f:
        json.dump(nodes, f)
    with open("data/fraud_edges.json", "w") as f:
        json.dump(edges, f)
    
    print(f"Generated {num_nodes} nodes and {len(edges)} edges for fraud dataset.")

if __name__ == "__main__":
    import os
    os.makedirs("data", exist_ok=True)
    generate_fraud_dataset()
