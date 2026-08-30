import torch
from neo4j import GraphDatabase
from torch_geometric.data import Data
import numpy as np

class Neo4jIngestor:
    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))

    def close(self):
        self.driver.close()

    def get_graph_data(self):
        with self.driver.session() as session:
            # Fetch nodes (Identities)
            nodes_result = session.run("MATCH (n:Identity) RETURN n.id as id, n.risk_score as score")
            node_map = {}
            node_features = []
            for i, record in enumerate(nodes_result):
                node_map[record["id"]] = i
                # Simple feature vector [risk_score, is_verified (dummy)]
                node_features.append([record["score"] or 0.0, 1.0])
            
            # Fetch edges (Relationships)
            edges_result = session.run("MATCH (a:Identity)-[r]->(b:Identity) RETURN a.id as source, b.id as target")
            edge_index = []
            for record in edges_result:
                if record["source"] in node_map and record["target"] in node_map:
                    edge_index.append([node_map[record["source"]], node_map[record["target"]]])
            
            x = torch.tensor(node_features, dtype=torch.float)
            edge_index = torch.tensor(edge_index, dtype=torch.long).t().contiguous()
            
            return Data(x=x, edge_index=edge_index)
