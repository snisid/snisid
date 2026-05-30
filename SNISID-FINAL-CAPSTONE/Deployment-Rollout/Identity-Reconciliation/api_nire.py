from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict, Any, Tuple
from reconciliation_logic import NationalIdentityReconciliationEngine

app = FastAPI(
    title="SNISID NIRE API",
    description="API Nationale du Moteur de Réconciliation d'Identité (Phase 15)",
    version="1.0.0"
)

class ReconcileRequest(BaseModel):
    new_record: Dict[str, Any]
    existing_record: Dict[str, Any]

class ReconcileResponse(BaseModel):
    decision: str
    details: Dict[str, Any]

@app.post("/api/v1/reconcile", response_model=ReconcileResponse)
def reconcile_identities(request: ReconcileRequest):
    """
    Évalue une nouvelle identité par rapport à une identité existante 
    pour détecter les doublons ou tentatives d'usurpation.
    """
    try:
        decision, details = NationalIdentityReconciliationEngine.reconcile(
            request.new_record, 
            request.existing_record
        )
        return ReconcileResponse(decision=decision, details=details)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Erreur interne du NIRE: {str(e)}")

@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "National Identity Reconciliation Engine"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
