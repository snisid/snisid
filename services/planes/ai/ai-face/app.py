import torch
import cv2
import numpy as np
import mlflow
from fastapi import FastAPI, File, UploadFile, HTTPException
from model import ArcFaceModel
from torchvision import transforms
from PIL import Image
import io

app = FastAPI(title="SNISID Face Recognition")

# Load model
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
model = ArcFaceModel().to(device)
model.eval()

# Tracking
mlflow.set_experiment("Identity_Verification")

preprocess = transforms.Compose([
    transforms.Resize((224, 224)),
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
])

@app.post("/extract-embeddings")
async def extract_embeddings(file: UploadFile = File(...)):
    try:
        content = await file.read()
        image = Image.open(io.BytesIO(content)).convert('RGB')
        tensor = preprocess(image).unsqueeze(0).to(device)
        
        with torch.no_grad():
            embeddings = model(tensor)
            
        with mlflow.start_run(nested=True):
            mlflow.log_param("device", str(device))
            mlflow.log_metric("embedding_norm", torch.norm(embeddings).item())

        return {"embeddings": embeddings.cpu().numpy().tolist()[0]}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/verify")
async def verify(file1: UploadFile = File(...), file2: UploadFile = File(...)):
    # Simple cosine similarity verification
    emb1 = await extract_embeddings(file1)
    emb2 = await extract_embeddings(file2)
    
    v1 = np.array(emb1["embeddings"])
    v2 = np.array(emb2["embeddings"])
    
    similarity = np.dot(v1, v2) / (np.norm(v1) * np.norm(v2))
    
    return {
        "verified": bool(similarity > 0.65),
        "similarity": float(similarity)
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
