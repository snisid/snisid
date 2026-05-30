import os
import logging
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import mlflow.pyfunc

# Configuration
MLFLOW_TRACKING_URI = os.getenv("MLFLOW_TRACKING_URI", "http://mlflow-server:5000")
MODEL_NAME = os.getenv("MODEL_NAME", "nexus-fraud-model")
MODEL_STAGE = os.getenv("MODEL_STAGE", "Production")

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("ai-runtime")

app = FastAPI(title="SNISID AI Inference Runtime")

# Load model from registry
try:
    mlflow.set_tracking_uri(MLFLOW_TRACKING_URI)
    model_uri = f"models:/{MODEL_NAME}/{MODEL_STAGE}"
    model = mlflow.pyfunc.load_model(model_uri)
    logger.info(f"Loaded model: {model_uri}")
except Exception as e:
    logger.warning(f"Could not load model from MLflow: {e}. Falling back to rule-based inference.")
    model = None

class InferenceRequest(BaseModel):
    id: str
    features: dict

@app.post("/v1/predict")
async def predict(req: InferenceRequest):
    if model:
        try:
            # Predict using MLflow model
            prediction = model.predict(req.features)
            return {"id": req.id, "prediction": prediction.tolist()}
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Inference error: {str(e)}")
    else:
        # Fallback rule-based logic for demo
        score = req.features.get("score", 0.0)
        risk = "HIGH" if score > 0.85 else "LOW"
        return {"id": req.id, "prediction": {"risk": risk}}

@app.get("/health")
async def health():
    return {"status": "up", "model_loaded": model is not None}
