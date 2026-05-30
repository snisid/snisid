import json
import time
from kafka import KafkaConsumer, KafkaProducer

class FraudInvestigatorAgent:
    def __init__(self, broker="localhost:9092"):
        self.consumer = KafkaConsumer(
            "swarm.tasks",
            bootstrap_servers=[broker],
            value_deserializer=lambda v: json.loads(v.decode('utf-8'))
        )
        self.producer = KafkaProducer(
            bootstrap_servers=[broker],
            value_serializer=lambda v: json.dumps(v).encode('utf-8')
        )
        self.agent_id = "agent-fraud-01"

    def run(self):
        print(f"[{self.agent_id}] Online. Waiting for tasks...")
        for message in self.consumer:
            task = message.value
            if task["agentType"] == "FRAUD_INVESTIGATOR":
                self._handle_task(task)

    def _handle_task(self, task):
        print(f"[{self.agent_id}] Processing Task: {task['command']} on {task['context']}")
        
        # Simulate investigation
        time.Sleep(2)
        
        result = {
            "taskId": task["id"],
            "status": "completed",
            "result": f"Analysis of {task['context']} confirms 85% probability of synthetic fraud via GNN analysis."
        }
        
        self.producer.send("swarm.results", result)
        print(f"[{self.agent_id}] Result sent for {task['id']}")

if __name__ == "__main__":
    agent = FraudInvestigatorAgent()
    agent.run()
