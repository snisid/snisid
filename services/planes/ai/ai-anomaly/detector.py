import json
import numpy as np
from kafka import KafkaConsumer
from sklearn.ensemble import IsolationForest
import mlflow

class AnomalyDetector:
    def __init__(self, broker="localhost:9092"):
        self.consumer = KafkaConsumer(
            "identity.created",
            bootstrap_servers=[broker],
            value_deserializer=lambda v: json.loads(v.decode('utf-8'))
        )
        self.model = IsolationForest(contamination=0.05)
        self.is_trained = False
        self.buffer = []

    def process_stream(self):
        print("Anomaly Detector started, monitoring identity events...")
        for message in self.consumer:
            data = message.value
            # Extract features (e.g., hour of day, location hash, etc.)
            features = self._extract_features(data)
            
            if not self.is_trained:
                self.buffer.append(features)
                if len(self.buffer) >= 100:
                    self._train()
                continue
            
            prediction = self.model.predict([features])
            if prediction[0] == -1:
                self._report_anomaly(data)

    def _extract_features(self, data):
        # Dummy feature extraction
        return [len(data.get("firstName", "")), len(data.get("lastName", ""))]

    def _train(self):
        print("Training isolation forest on initial buffer...")
        self.model.fit(self.buffer)
        self.is_trained = True
        mlflow.log_param("model_type", "IsolationForest")
        mlflow.log_metric("train_size", len(self.buffer))

    def _report_anomaly(self, data):
        print(f"ANOMALY DETECTED: {data['identityId']}")
        # In production, this would publish to the 'alert.generated' Kafka topic

if __name__ == "__main__":
    detector = AnomalyDetector()
    detector.process_stream()
