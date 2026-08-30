# MP-28 — RESO-HT
## Analyse des Réseaux Criminels d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-28 | Code : RESO-HT
Dépendances      : GANG-HT, CHEF-HT, BLAN-HT, UCREF-INT, SIVC-HT, FIR-HT
Normes           : EUROPOL SIENA SNA, INTERPOL i-Link, Graph ML (ICCJ), Palantir methodology
Acteurs          : DCPJ BAC Analytique, Cellule Renseignement, DEA liaison, UNODC
```

---

## 1. CONTEXTE

L'analyse des réseaux criminels (SNA) révèle les connexions cachées entre gangs,
personnalités politiques corrompues, financiers et trafiquants. Ce module orchestre
le graphe criminologique central de SNISID — alimenté par tous les modules — et
fournit algorithmes de détection de communautés, scoring de centralité et prédiction
de liens via GNN PyTorch (déjà intégré dans la stack ML SNISID).

---

## 2. ARCHITECTURE

```
┌──────────────────────────────────────────────────────────────┐
│              RESO-HT — GRAPHE CRIMINOLOGIQUE CENTRAL          │
├──────────────────────────────────────────────────────────────┤
│  Kafka Consumers (alimentation continue depuis tous modules)  │
│  GANG-HT / CHEF-HT / FIR-HT / SIVC-HT / BLAN-HT / UCREF    │
├──────────────────────────────────────────────────────────────┤
│  Neo4j 5.x Cluster (3 noeuds — Port-au-Prince + Cap-Haïtien) │
│  Noeuds : Person, Gang, Vehicle, Location, Crime, Account    │
│  Aretes : MEMBER_OF, LED_BY, ASSOCIATED_WITH, TRANSACTS, ... │
├──────────────────────────────────────────────────────────────┤
│  Python ML Service (RESO-ML)                                  │
│  - Louvain community detection                               │
│  - PageRank / Betweenness centrality scoring                 │
│  - GNN Link Prediction (PyTorch Geometric)                   │
│  - Isolation Forest (anomaly detection)                      │
└──────────────────────────────────────────────────────────────┘
```

---

## 3. SCHÉMA NEO4J COMPLET

```cypher
// === NOEUDS ===
(:Person {
    snisid_id: String,
    name: String,
    aliases: [String],
    nationality: String,
    dob: Date,
    risk_score: Float,
    is_gang_member: Boolean,
    is_sanctioned: Boolean
})

(:Gang {
    gang_id: String,
    name: String,
    primary_activity: String,
    territory_dept: String,
    activity_level: String,
    member_count: Integer
})

(:Vehicle {
    plate_number: String,
    vin: String,
    make: String,
    model: String,
    alert_level: String,
    crime_categories: [String]
})

(:Location {
    dept_code: String,
    commune: String,
    lat: Float,
    lng: Float,
    is_gang_territory: Boolean,
    control_level: String
})

(:Crime {
    crime_id: String,
    category: String,
    date: DateTime,
    dept_code: String,
    case_reference: String,
    casualties: Integer
})

(:FinancialAccount {
    account_id: String,
    account_type: String,
    institution: String,
    country: String,
    risk_level: String,
    is_flagged: Boolean
})

// === RELATIONS CLES ===
// Organisation criminelle
(:Person)-[:MEMBER_OF {role: String, since: Date, confidence: Float}]->(:Gang)
(:Person)-[:LEADS {since: Date, confirmed: Boolean}]->(:Gang)
(:Gang)-[:ALLIED_WITH {type: String, since: Date, confidence: Float}]->(:Gang)
(:Gang)-[:RIVALS_WITH {since: Date}]->(:Gang)

// Criminalite
(:Person)-[:PARTICIPATED_IN {role: String, date: DateTime}]->(:Crime)
(:Gang)-[:INVOLVED_IN {role: String, date: DateTime}]->(:Crime)
(:Person)-[:DROVE {date: DateTime, confirmed: Boolean}]->(:Vehicle)
(:Vehicle)-[:LINKED_TO_CRIME {role: String, date: DateTime}]->(:Crime)
(:Crime)-[:OCCURRED_AT]->(:Location)

// Geographie
(:Person)-[:SIGHTED_AT {date: DateTime, source: String, confidence: Float}]->(:Location)
(:Gang)-[:CONTROLS {level: String, since: Date}]->(:Location)

// Finance
(:Person)-[:OWNS_ACCOUNT {opened: Date}]->(:FinancialAccount)
(:FinancialAccount)-[:TRANSACTS_WITH {
    amount: Float, date: DateTime, suspicious: Boolean
}]->(:FinancialAccount)

// Interconnexion personnes
(:Person)-[:ASSOCIATES_WITH {
    type: String, confidence: Float, source: String
}]->(:Person)
(:Person)-[:FAMILY_OF {relation: String}]->(:Person)
```

---

## 4. SERVICE PYTHON ML

```python
# services/reso-ml-svc/analysis/network_analyzer.py
from py2neo import Graph
import networkx as nx
from torch_geometric.data import Data
from torch_geometric.nn import GCNConv
import torch
import numpy as np
from typing import List, Dict, Tuple


class CriminalNetworkAnalyzer:
    def __init__(self, neo4j_uri: str, neo4j_password: str):
        self.graph = Graph(neo4j_uri, password=neo4j_password)

    def detect_communities(self) -> Dict:
        """Detecte les communautes criminelles via algorithme Louvain"""
        G = self._export_to_networkx()
        try:
            import community as community_louvain
            partition = community_louvain.best_partition(G)
            modularity = community_louvain.modularity(partition, G)
            communities = {}
            for node, comm in partition.items():
                communities.setdefault(comm, []).append(node)
            return {
                "num_communities": len(communities),
                "modularity": modularity,
                "communities": communities,
                "largest_community_size": max(len(v) for v in communities.values())
            }
        except ImportError:
            # Fallback: Label Propagation
            labels = nx.community.label_propagation_communities(G)
            return {"communities": [list(c) for c in labels]}

    def compute_key_actors(self, top_k: int = 20) -> List[Dict]:
        """Identifie les acteurs cles par centralite (brokers du reseau)"""
        G = self._export_to_networkx()
        pagerank = nx.pagerank(G, alpha=0.85)
        betweenness = nx.betweenness_centrality(G, normalized=True)
        degree = dict(G.degree())

        actors = []
        for node in G.nodes():
            actors.append({
                "node_id": node,
                "pagerank": pagerank.get(node, 0),
                "betweenness": betweenness.get(node, 0),
                "degree": degree.get(node, 0),
                "composite_score": (
                    pagerank.get(node, 0) * 0.4 +
                    betweenness.get(node, 0) * 0.4 +
                    (degree.get(node, 0) / max(degree.values())) * 0.2
                )
            })
        actors.sort(key=lambda x: x["composite_score"], reverse=True)
        return actors[:top_k]

    def find_shortest_path(self, from_id: str, to_id: str) -> List:
        """Chemin le plus court entre deux acteurs criminels"""
        G = self._export_to_networkx()
        try:
            path = nx.shortest_path(G, from_id, to_id)
            return [{"node": n, "degree": G.degree(n)} for n in path]
        except nx.NetworkXNoPath:
            return []

    def _export_to_networkx(self) -> nx.Graph:
        """Exporte le graphe Neo4j vers NetworkX pour analyse"""
        G = nx.Graph()
        result = self.graph.run("""
            MATCH (a:Person)-[r]-(b:Person)
            RETURN a.snisid_id AS src, b.snisid_id AS dst,
                   type(r) AS rel_type, r.confidence AS weight
            LIMIT 50000
        """)
        for record in result:
            G.add_edge(
                record["src"],
                record["dst"],
                rel_type=record["rel_type"],
                weight=record.get("weight", 1.0)
            )
        return G
```

---

## 5. API REST

| Méthode | Endpoint                               | Rôle         | Description                        |
|---------|----------------------------------------|--------------|------------------------------------|
| `GET`   | `/api/v1/reso/network/:person_id`      | DCPJ_INTEL   | Réseau direct d'un individu        |
| `GET`   | `/api/v1/reso/communities`             | DCPJ_INTEL   | Communautés criminelles détectées  |
| `GET`   | `/api/v1/reso/key-actors`              | DCPJ_INTEL   | Top acteurs-clés du réseau         |
| `GET`   | `/api/v1/reso/gang-overlap/:g1/:g2`    | DCPJ         | Membres communs deux gangs         |
| `POST`  | `/api/v1/reso/analyze/trigger`         | DCPJ_ADMIN   | Lancer analyse complète réseau     |
| `GET`   | `/api/v1/reso/path/:from_id/:to_id`    | DCPJ_INTEL   | Chemin entre deux acteurs          |
| `GET`   | `/api/v1/reso/centrality-scores`       | DCPJ_INTEL   | Scores centralité complets         |
| `GET`   | `/api/v1/reso/emerging-links`          | DCPJ_INTEL   | Liens prédits par GNN              |

---

## 6. KAFKA CONSUMERS — ALIMENTATION DU GRAPHE

```go
// Consommer tous les événements des modules sources
topics := []string{
    "gang.organization.created",
    "chef.member.created",
    "chef.status.changed",
    "sivc.alerts.created",
    "fir.conviction.added",
    "blan.transaction.suspicious",
    "dipe.person.found",
}
// Chaque événement est traduit en noeud/relation Neo4j
// via le GraphIngestionService
```

---

## 7. VARIABLES D'ENVIRONNEMENT

```dotenv
RESO_NEO4J_URI=bolt://neo4j:7687
RESO_NEO4J_USER=neo4j
RESO_NEO4J_PASSWORD=<VAULT:reso/neo4j_password>
RESO_ML_SERVICE_URL=http://reso-ml-svc:8199
RESO_KAFKA_BROKERS=kafka:9092
RESO_ANALYSIS_INTERVAL_HOURS=12
RESO_KEY_ACTOR_THRESHOLD=0.05
RESO_MAX_GRAPH_DEPTH=6
RESO_SERVICE_PORT=8101
RESO_ML_PORT=8199
```

---
*MP-28 — RESO-HT — Analyse Réseaux Criminels — SNISID — République d'Haïti*
