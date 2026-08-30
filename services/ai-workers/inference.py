from fastapi import FastAPI, UploadFile, File
import numpy as np
import time

app = FastAPI(title="SNISID-AI-Worker")

@app.post("/v1/biometric/embed")
async def extract_embedding(file: UploadFile = File(...)):
    print("AI-WORKER: Extracting ArcFace embedding...")
    # Simulated 512D vector
    embedding = np.random.randn(512).tolist()
    return {"id": "face_001", "vector": embedding}

@app.post("/v1/security/deepfake")
async def detect_deepfake(file: UploadFile = File(...)):
    print("AI-WORKER: Running deepfake artifact detection...")
    score = np.random.uniform(0, 1)
    return {
        "deepfake_score": score,
        "is_fake": score > 0.7,
        "engine": "transformer-v2"
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5000)
