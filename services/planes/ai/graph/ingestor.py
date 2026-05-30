# ingestor.py
import json
from kafka import KafkaConsumer
from neo4j import GraphDatabase

class GraphIngestor:
    """
    Sovereign Graph Intelligence.
    Ingests identity events into Neo4j to detect suspicious relationship clusters.
    """
    def __init__(self, uri="bolt://neo4j:7687", user="neo4j", password="password"):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))
        print("NEXUS-GRAPH: Identity Ingestor Operational.")

    def process_events(self):
        consumer = KafkaConsumer(
            "identity-events",
            bootstrap_servers=["kafka:9092"],
            value_deserializer=lambda m: json.loads(m.decode("utf-8"))
        )

        for msg in consumer:
            data = msg.value
            with self.driver.session() as session:
                session.execute_write(self._merge_citizen, data)
                print(f"NEXUS-GRAPH: Integrated citizen {data['id']} into global graph.")

    @staticmethod
    def _merge_citizen(tx, data):
        # Build relationship graph
        tx.run(
            "MERGE (c:Citizen {id: $id}) "
            "SET c.name = $name, c.risk = $risk",
            id=data["id"], name=data.get("name"), risk=data.get("risk", 0.0)
        )
        if "address_id" in data:
            tx.run(
                "MATCH (c:Citizen {id: $id}) "
                "MERGE (a:Address {id: $addr_id}) "
                "MERGE (c)-[:LIVES_AT]->(a)",
                id=data["id"], addr_id=data["address_id"]
            )

if __name__ == "__main__":
    ingestor = GraphIngestor()
    # ingestor.process_events() # Running in background
