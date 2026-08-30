import torch
import pytest
import numpy as np
from model import ArcFaceModel

def test_model_output_shape():
    model = ArcFaceModel(embedding_size=512)
    model.eval()
    
    # Mock input batch: 1 image, 3 channels, 224x224
    dummy_input = torch.randn(1, 3, 224, 224)
    
    with torch.no_grad():
        output = model(dummy_input)
    
    assert output.shape == (1, 512)
    # Check if normalized
    norm = torch.norm(output, dim=1).item()
    assert pytest.approx(norm, 0.001) == 1.0

def test_cosine_similarity():
    # Verify that identical inputs yield similarity ~1.0
    v1 = np.array([1.0, 0.0, 0.0])
    v2 = np.array([1.0, 0.0, 0.0])
    
    sim = np.dot(v1, v2) / (np.linalg.norm(v1) * np.linalg.norm(v2))
    assert sim == 1.0
    
    # Verify orthogonal inputs yield similarity ~0.0
    v3 = np.array([0.0, 1.0, 0.0])
    sim_ortho = np.dot(v1, v3) / (np.linalg.norm(v1) * np.linalg.norm(v3))
    assert sim_ortho == 0.0
