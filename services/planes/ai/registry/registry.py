import logging
import time
from typing import Dict, List
from pydantic import BaseModel

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("model-registry")

class ModelMetadata(BaseModel):
    name: str
    version: str
    stage: str  # Production, Staging, Archive
    performance_metrics: Dict[str, float]
    last_updated: float

class ModelRegistry:
    def __init__(self):
        self.models: Dict[str, List[ModelMetadata]] = {}

    def register_model(self, metadata: ModelMetadata):
        if metadata.name not in self.models:
            self.models[metadata.name] = []
        self.models[metadata.name].append(metadata)
        logger.info(f"Registered model {metadata.name} v{metadata.version}")

    def detect_drift(self, model_name: str, current_metrics: Dict[str, float]) -> bool:
        if model_name not in self.models or not self.models[model_name]:
            return False
        
        baseline = self.models[model_name][-1].performance_metrics
        for metric, value in current_metrics.items():
            if metric in baseline:
                # Simple drift detection: 20% degradation
                if value < baseline[metric] * 0.8:
                    logger.warning(f"Drift detected for {model_name} in {metric}: {value} vs {baseline[metric]}")
                    return True
        return False

# Example usage
if __name__ == "__main__":
    registry = ModelRegistry()
    registry.register_model(ModelMetadata(
        name="fraud-detection",
        version="1.0.0",
        stage="Production",
        performance_metrics={"accuracy": 0.95, "f1": 0.92},
        last_updated=time.time()
    ))
    
    drifted = registry.detect_drift("fraud-detection", {"accuracy": 0.70, "f1": 0.65})
    print(f"Drift detected: {drifted}")
