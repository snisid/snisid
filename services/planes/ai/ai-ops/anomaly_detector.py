import time
import requests
import numpy as np
from sklearn.ensemble import IsolationForest
import mlflow

class AIOpsAnomalyDetector:
    def __init__(self, prometheus_url="http://prometheus:9090"):
        self.prometheus_url = prometheus_url
        self.model = IsolationForest(contamination=0.01)
        self.is_trained = False
        self.history = []

    def fetch_metrics(self):
        # Query Prometheus for total request rate
        query = 'rate(http_requests_total[5m])'
        response = requests.get(f"{self.prometheus_url}/api/v1/query", params={'query': query})
        data = response.json()
        
        values = []
        if data['status'] == 'success':
            for result in data['data']['result']:
                values.append(float(result['value'][1]))
        return values

    def run(self):
        print("AIOps Anomaly Detector started...")
        while True:
            metrics = self.fetch_metrics()
            if not metrics:
                time.sleep(60)
                continue

            if not self.is_trained:
                self.history.extend(metrics)
                if len(self.history) >= 100:
                    self._train()
                continue

            # Detect anomalies
            X = np.array(metrics).reshape(-1, 1)
            predictions = self.model.predict(X)
            
            for i, pred in enumerate(predictions):
                if pred == -1:
                    self._report_anomaly(metrics[i])
            
            time.sleep(60)

    def _train(self):
        print("Training AIOps model on metric history...")
        X = np.array(self.history).reshape(-1, 1)
        self.model.fit(X)
        self.is_trained = True
        mlflow.log_param("aio_model", "IsolationForest")

    def _report_anomaly(self, value):
        print(f"⚠️ AIOps ALERT: Anomaly detected in traffic rate: {value}")
        # In production, this would send an alert to the SOC Alert Feed via Kafka

if __name__ == "__main__":
    detector = AIOpsAnomalyDetector()
    detector.run()
