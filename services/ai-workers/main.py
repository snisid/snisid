from fastapi import FastAPI, File, UploadFile
import numpy as np

app = FastAPI(title="SNISID AI Worker")

# Mock models for orchestration
class MockModel:
    def extract(self, img): return [0.12] * 512
    def predict(self, img): return 0.85

arcface_model = MockModel()
deepfake_model = MockModel()

@app.post("/embed")
async def embed(image: bytes = File(...)):
    """
    Extracts biometric embedding vector from image.
    """
    vector = arcface_model.extract(image)
    return {"vector": vector, "model": "arcface_v2"}

@app.post("/deepfake")
async def detect(image: bytes = File(...)):
    """
    Detects deepfake probability in image frame.
    """
    score = deepfake_model.predict(image)
    return {
        "deepfake_score": float(score),
        "is_fake": score > 0.7,
        "recommendation": "ISOLATE" if score > 0.9 else "MONITOR"
    }

@app.get("/health")
def health():
    return {"status": "AI_WORKER_ONLINE"}
