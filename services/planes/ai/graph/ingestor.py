import json
import os
from kafka import KafkaConsumer
from neo4j import GraphDatabase


class GraphIngestor:
    def __init__(self, uri=None, user=None, password=None):
        self.uri = uri or os.getenv("NEO4J_URI", "bolt://neo4j:7687")
        self.user = user or os.getenv("NEO4J_USER", "neo4j")
        self.password = password or os.getenv("NEO4J_PASSWORD", "")
        self.driver = GraphDatabase.driver(self.uri, auth=(self.user, self.password))

    def process_events(self):
        consumer = KafkaConsumer(
            "identity-events",
            bootstrap_servers=os.getenv("KAFKA_BROKERS", "kafka:9092"),
            value_deserializer=lambda m: json.loads(m.decode("utf-8"))
        )

        for msg in consumer:
            data = msg.value
            with self.driver.session() as session:
                session.execute_write(self._merge_citizen, data)

    @staticmethod
    def _merge_citizen(tx, data):
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
