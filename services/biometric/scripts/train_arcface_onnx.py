#!/usr/bin/env python3
"""Train ArcFace model and export to ONNX for NPU deployment.

Usage:
    python scripts/train_arcface_onnx.py \
        --data-dir /path/to/faces \
        --epochs 50 \
        --batch-size 64 \
        --embedding-dim 512 \
        --num-classes 1000 \
        --output ./models/arcface.onnx

This script:
  1. Loads or creates synthetic training data (replace with real dataset).
  2. Trains ArcFace with ArcMarginProduct loss.
  3. Exports the backbone to ONNX.
  4. Validates ONNX output matches PyTorch output.
  5. Optimizes the ONNX graph for Qualcomm AI 100.
  6. Prints NPU deployment instructions.
"""

import argparse
import logging
import os
import sys
import time
from pathlib import Path

import numpy as np
import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
from torch.utils.data import DataLoader, Dataset, TensorDataset

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("train_arcface_onnx")


class ArcFaceModel(nn.Module):
    """ResNet50-based ArcFace model for face recognition."""

    def __init__(self, embedding_size=512):
        super().__init__()
        from torchvision.models import resnet50
        self.backbone = resnet50(weights=None)
        in_features = self.backbone.fc.in_features
        self.backbone.fc = nn.Linear(in_features, embedding_size)
        self.bn = nn.BatchNorm1d(embedding_size)

    def forward(self, x):
        embeddings = self.backbone(x)
        embeddings = self.bn(embeddings)
        return F.normalize(embeddings)


class ArcMarginProduct(nn.Module):
    """ArcFace additive angular margin loss."""

    def __init__(self, in_features, out_features, s=30.0, m=0.50):
        super().__init__()
        self.in_features = in_features
        self.out_features = out_features
        self.s = s
        self.m = m
        self.weight = nn.Parameter(torch.FloatTensor(out_features, in_features))
        nn.init.xavier_uniform_(self.weight)

    def forward(self, input_ids, labels):
        cosine = F.linear(F.normalize(input_ids), F.normalize(self.weight))
        sine = torch.sqrt(1.0 - torch.pow(cosine, 2) + 1e-8)
        phi = cosine * torch.cos(torch.tensor(self.m)) - sine * torch.sin(torch.tensor(self.m))
        one_hot = torch.zeros_like(cosine)
        one_hot.scatter_(1, labels.view(-1, 1).long(), 1)
        output = (one_hot * phi) + ((1.0 - one_hot) * cosine)
        output *= self.s
        return output


def generate_synthetic_dataset(
    num_classes: int,
    samples_per_class: int,
    image_size: int = 112,
    embedding_dim: int = 512,
) -> tuple[Dataset, int]:
    """Generate synthetic face data for demonstration.

    Replace this with a real dataset loader in production.
    """
    total = num_classes * samples_per_class
    images = torch.randn(total, 3, image_size, image_size)
    labels = torch.arange(num_classes, dtype=torch.long).repeat_interleave(samples_per_class)
    logger.info(
        "Generated synthetic dataset: %d samples, %d classes, %dx%d",
        total,
        num_classes,
        image_size,
        image_size,
    )
    return TensorDataset(images, labels), num_classes


def train_epoch(
    model: nn.Module,
    margin: ArcMarginProduct,
    loader: DataLoader,
    optimizer: optim.Optimizer,
    device: torch.device,
) -> float:
    model.train()
    total_loss = 0.0

    for images, labels in loader:
        images, labels = images.to(device), labels.to(device)

        embeddings = model(images)
        output = margin(embeddings, labels)
        loss = F.cross_entropy(output, labels)

        optimizer.zero_grad()
        loss.backward()
        optimizer.step()

        total_loss += loss.item()

    return total_loss / len(loader)


@torch.no_grad()
def validate(
    model: nn.Module,
    margin: ArcMarginProduct,
    loader: DataLoader,
    device: torch.device,
) -> float:
    model.eval()
    correct = 0
    total = 0

    for images, labels in loader:
        images, labels = images.to(device), labels.to(device)
        embeddings = model(images)
        output = margin(embeddings, labels)
        preds = output.argmax(dim=1)
        correct += (preds == labels).sum().item()
        total += labels.size(0)

    return correct / total


def main():
    parser = argparse.ArgumentParser(description="Train ArcFace and export to ONNX")
    parser.add_argument("--data-dir", type=str, default=None, help="Path to face dataset (optional, uses synthetic if not set)")
    parser.add_argument("--epochs", type=int, default=10, help="Number of training epochs")
    parser.add_argument("--batch-size", type=int, default=64, help="Batch size")
    parser.add_argument("--embedding-dim", type=int, default=512, help="Embedding dimension")
    parser.add_argument("--num-classes", type=int, default=100, help="Number of identities (classes)")
    parser.add_argument("--samples-per-class", type=int, default=20, help="Samples per identity")
    parser.add_argument("--learning-rate", type=float, default=0.01, help="Initial learning rate")
    parser.add_argument("--output", type=str, default="./models/arcface.onnx", help="Output ONNX path")
    parser.add_argument("--checkpoint", type=str, default="./models/arcface_checkpoint.pt", help="Checkpoint path")
    parser.add_argument("--device", type=str, default="auto", help="Device: auto, cpu, or cuda")
    args = parser.parse_args()

    if args.device == "auto":
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    else:
        device = torch.device(args.device)

    logger.info("Using device: %s", device)
    logger.info("Training config: epochs=%d, batch=%d, lr=%.4f, embedding=%d",
                args.epochs, args.batch_size, args.learning_rate, args.embedding_dim)

    # 1. Load data
    if args.data_dir and os.path.isdir(args.data_dir):
        logger.info("Loading real dataset from %s ...", args.data_dir)
        # Stub: replace with ImageFolder or custom dataset
        dataset, num_classes = generate_synthetic_dataset(
            args.num_classes, args.samples_per_class, embedding_dim=args.embedding_dim
        )
    else:
        logger.info("No dataset provided; using synthetic data for demonstration.")
        dataset, num_classes = generate_synthetic_dataset(
            args.num_classes, args.samples_per_class, embedding_dim=args.embedding_dim
        )

    train_size = int(0.8 * len(dataset))
    val_size = len(dataset) - train_size
    train_ds, val_ds = torch.utils.data.random_split(dataset, [train_size, val_size])
    train_loader = DataLoader(train_ds, batch_size=args.batch_size, shuffle=True, num_workers=0)
    val_loader = DataLoader(val_ds, batch_size=args.batch_size, shuffle=False, num_workers=0)

    # 2. Build model
    model = ArcFaceModel(embedding_size=args.embedding_dim).to(device)
    margin = ArcMarginProduct(args.embedding_dim, num_classes).to(device)
    optimizer = optim.SGD(
        list(model.parameters()) + list(margin.parameters()),
        lr=args.learning_rate,
        momentum=0.9,
        weight_decay=5e-4,
    )
    scheduler = optim.lr_scheduler.CosineAnnealingLR(optimizer, T_max=args.epochs)

    # 3. Train
    best_acc = 0.0
    output_dir = Path(args.output).parent
    output_dir.mkdir(parents=True, exist_ok=True)

    for epoch in range(1, args.epochs + 1):
        start = time.perf_counter()
        train_loss = train_epoch(model, margin, train_loader, optimizer, device)
        val_acc = validate(model, margin, val_loader, device)
        scheduler.step()
        elapsed = time.perf_counter() - start

        logger.info(
            "Epoch %3d/%d | loss=%.4f | val_acc=%.4f | lr=%.6f | %.2fs",
            epoch, args.epochs, train_loss, val_acc,
            scheduler.get_last_lr()[0], elapsed,
        )

        if val_acc > best_acc:
            best_acc = val_acc
            torch.save(model.state_dict(), args.checkpoint)
            logger.info("Checkpoint saved: %s (acc=%.4f)", args.checkpoint, best_acc)

    logger.info("Training complete. Best validation accuracy: %.4f", best_acc)

    # 4. Export to ONNX
    model.load_state_dict(torch.load(args.checkpoint, map_location=device, weights_only=True))
    model.to(device)
    model.eval()

    from services.biometric.inference.onnx_export import ONNXExporter

    onnx_path = ONNXExporter.export_arcface(
        model_path=args.checkpoint,
        output_path=args.output,
        input_shape=(1, 3, 112, 112),
        dynamic_batch=True,
        opset_version=17,
    )
    logger.info("ONNX export complete: %s", onnx_path)

    # 5. Validate ONNX
    validation = ONNXExporter.validate_onnx(onnx_path, atol=1e-3)
    logger.info(
        "ONNX validation: %s (max_diff=%.6f)",
        "PASSED" if validation["passed"] else "FAILED",
        validation["max_diff"],
    )

    if not validation["passed"]:
        logger.error("ONNX output diverges from PyTorch. Check export settings.")
        sys.exit(1)

    # 6. Optimize for NPU
    optimized_path = ONNXExporter.optimize_for_npu(onnx_path, npu="qaic")
    logger.info("Optimized ONNX for QAIC: %s", optimized_path)

    # 7. Deployment instructions
    sep = "=" * 60
    instructions = f"""
{sep}
NPU DEPLOYMENT INSTRUCTIONS
{sep}

Model exported to: {onnx_path}
Optimized model:   {optimized_path}
Embedding dim:     {args.embedding_dim}
Input shape:       (1, 3, 112, 112)

To deploy on Qualcomm AI 100 Edge:

 1. Transfer the optimized ONNX model to the edge device:
    scp {optimized_path} user@edge-terminal:/models/arcface.onnx

 2. On the edge terminal, run the NPU runtime:
    python -c "
    from services.biometric.inference.npu_runtime import NPURuntime
    runtime = NPURuntime.create('/models/arcface.onnx', preferred_backend='qaic')
    print('Backends:', runtime.available_backends)
    "

 3. Verify inference with a test image:
    python -c "
    import cv2, numpy as np
    img = cv2.imread('/path/to/face.jpg')
    img = cv2.resize(img, (112, 112)).astype(np.float32) / 255.0
    tensor = np.expand_dims(np.transpose(img, (2, 0, 1)), axis=0)
    emb = runtime.infer(tensor)
    print('Embedding shape:', emb.shape, 'Norm:', np.linalg.norm(emb))
    "

 4. For maximum throughput, batch multiple faces:
    python -c "
    batch = np.concatenate([tensor] * 16, axis=0)
    embs = runtime.infer(batch)
    print('Batch embedding shape:', embs.shape)
    "

 5. To benchmark latency:
    python scripts/benchmark_npu.py --model {optimized_path}

Requirements on edge:
  - onnxruntime>=1.18.0 with QNNExecutionProvider
  - Qualcomm AI 100 SDK (qaic-api)
  - Python 3.10+
{sep}
"""
    print(instructions)


if __name__ == "__main__":
    main()
