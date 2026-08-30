from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI(title="SNISID Anomaly Detection")


class AnomalyRequest(BaseModel):
    stream_id: str
    current_rate: float
    baseline_rate: float


@app.get("/healthz")
def healthz():
    return {"status": "ok"}


@app.post("/v1/analyze")
def analyze(req: AnomalyRequest):
    if req.baseline_rate <= 0:
        return {"stream_id": req.stream_id, "anomaly_score": 1.0, "anomalous": True}
    ratio = req.current_rate / req.baseline_rate
    score = min(1.0, abs(ratio - 1.0))
    return {"stream_id": req.stream_id, "anomaly_score": round(score, 4), "anomalous": score >= 0.4}
