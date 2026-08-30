#!/usr/bin/env python3
"""Benchmark NPU vs CPU vs CUDA for face matching inference.

Measures:
  - Latency for 1, 10, 100, 1000 inferences
  - 1:N search time for 100, 1K, 10K, 100K galleries
  - Throughput (identities/second)

Usage:
    python scripts/benchmark_npu.py --model ./models/arcface.onnx
"""

import argparse
import logging
import time
from pathlib import Path

import numpy as np

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
)
logger = logging.getLogger("benchmark_npu")


def benchmark_inference_latency(
    runtime, label: str, batch_sizes: list[int], num_warmup: int = 10
) -> dict[int, dict]:
    """Measure average latency for different batch sizes.

    Returns:
        {batch_size: {"avg_ms": ..., "std_ms": ..., "throughput": ...}}
    """
    results = {}
    rng = np.random.default_rng(42)

    for batch_size in batch_sizes:
        dummy = rng.standard_normal((batch_size, 3, 112, 112)).astype(np.float32)
        durations: list[float] = []

        # Warmup
        for _ in range(num_warmup):
            runtime.infer(dummy)

        # Timed runs
        num_runs = max(1, 200 // batch_size)
        for _ in range(num_runs):
            t0 = time.perf_counter()
            runtime.infer(dummy)
            durations.append((time.perf_counter() - t0) * 1000)  # ms

        avg = float(np.mean(durations))
        std = float(np.std(durations))
        throughput = (batch_size * len(durations)) / (sum(durations) / 1000) if durations else 0

        results[batch_size] = {
            "avg_ms": round(avg, 3),
            "std_ms": round(std, 3),
            "throughput": round(throughput),
        }
        logger.info(
            "  %s batch=%3d | avg=%.2f ms | std=%.2f ms | %.0f img/s",
            label.ljust(6), batch_size, avg, std, throughput,
        )

    return results


def benchmark_1n_search(
    gallery_sizes: list[int], dimension: int = 512, top_k: int = 5, num_queries: int = 50
) -> dict[int, dict]:
    """Benchmark 1:N FAISS search time for various gallery sizes.

    Returns:
        {gallery_size: {"avg_ms": ..., "throughput_ids_per_sec": ...}}
    """
    try:
        import faiss
    except ImportError:
        logger.error("faiss is required for 1:N benchmark. Install: pip install faiss-cpu")
        return {}

    results = {}
    rng = np.random.default_rng(42)

    for gallery_size in gallery_sizes:
        gallery = rng.standard_normal((gallery_size, dimension)).astype(np.float32)
        gallery = gallery / np.linalg.norm(gallery, axis=1, keepdims=True)

        index = faiss.IndexFlatIP(dimension)
        index.add(gallery)

        queries = rng.standard_normal((num_queries, dimension)).astype(np.float32)
        queries = queries / np.linalg.norm(queries, axis=1, keepdims=True)

        durations: list[float] = []
        for i in range(num_queries):
            t0 = time.perf_counter()
            index.search(queries[i : i + 1], top_k)
            durations.append((time.perf_counter() - t0) * 1000)

        avg = float(np.mean(durations))
        ids_per_sec = 1.0 / (avg / 1000) if avg > 0 else 0

        results[gallery_size] = {
            "avg_ms": round(avg, 3),
            "throughput_qps": round(ids_per_sec),
        }
        logger.info(
            "  1:N gallery=%7d | avg=%.3f ms | %.0f queries/s",
            gallery_size, avg, ids_per_sec,
        )

    return results


def main():
    parser = argparse.ArgumentParser(description="Benchmark NPU vs CPU vs CUDA")
    parser.add_argument("--model", type=str, required=True, help="Path to ONNX model")
    parser.add_argument("--backends", type=str, nargs="+",
                        default=["cpu"], choices=["cpu", "cuda", "qaic", "all"],
                        help="Backends to benchmark")
    parser.add_argument("--batch-sizes", type=int, nargs="+",
                        default=[1, 10, 100, 1000],
                        help="Batch sizes for latency test")
    parser.add_argument("--gallery-sizes", type=int, nargs="+",
                        default=[100, 1000, 10000, 100000],
                        help="Gallery sizes for 1:N search test")
    parser.add_argument("--dimension", type=int, default=512,
                        help="Embedding dimension")
    args = parser.parse_args()

    if not Path(args.model).exists():
        logger.error("Model not found: %s", args.model)
        return 1

    backends = ["cpu", "cuda", "qaic"] if "all" in args.backends else args.backends

    separator = "=" * 72
    title = f"NPU BENCHMARK: {Path(args.model).name}"
    print(f"\n{separator}")
    print(f"  {title}")
    print(f"  Dimension: {args.dimension} | Batch sizes: {args.batch_sizes}")
    print(f"  Gallery sizes: {args.gallery_sizes}")
    print(f"{separator}\n")

    # Step 1: Inference latency per backend
    print(f"{'BACKEND':<10} {'BATCH':>6} {'AVG (ms)':>10} {'STD (ms)':>10} {'IMG/S':>10}")
    print("-" * 50)

    from services.biometric.inference.npu_runtime import NPURuntime

    all_inference_results: dict[str, dict] = {}

    for backend in backends:
        try:
            logger.info("Initializing backend: %s ...", backend)
            runtime = NPURuntime.create(args.model, preferred_backend=backend)
        except (ImportError, RuntimeError) as exc:
            logger.warning("Backend '%s' unavailable: %s", backend, exc)
            continue

        label = backend.upper()
        inf_results = benchmark_inference_latency(
            runtime, label, args.batch_sizes
        )
        all_inference_results[backend] = inf_results

    # Step 2: 1:N search benchmark
    print(f"\n{'GALLERY':<12} {'AVG (ms)':>10} {'QUERIES/S':>12}")
    print("-" * 36)

    search_results = benchmark_1n_search(
        args.gallery_sizes, args.dimension, top_k=5, num_queries=50
    )

    # Step 3: Summary table
    print(f"\n{separator}")
    print("  SUMMARY: Throughput (images/second) by backend and batch size")
    print(f"{separator}")
    print(f"{'BATCH':>8}", end="")
    for backend in all_inference_results:
        print(f"  {backend.upper():>12}", end="")
    print()

    batch_sizes = args.batch_sizes
    for bs in batch_sizes:
        print(f"{bs:>8}", end="")
        for backend in all_inference_results:
            tp = all_inference_results[backend].get(bs, {}).get("throughput", 0)
            print(f"  {tp:>12,}", end="")
        print()

    # Step 4: Recommendation
    print(f"\n{separator}")
    print("  DEPLOYMENT RECOMMENDATION")
    print(f"{separator}")

    # Find fastest backend
    fastest_backend = None
    max_tp = 0
    for backend, inf_res in all_inference_results.items():
        tp_bs1 = inf_res.get(1, {}).get("throughput", 0)
        if tp_bs1 > max_tp:
            max_tp = tp_bs1
            fastest_backend = backend

    if fastest_backend:
        print(f"  Fastest single-inference backend: {fastest_backend.upper()}")
        print(f"  Peak throughput: {max_tp:,} images/second")
    if search_results:
        max_gal = max(args.gallery_sizes)
        gal_time = search_results.get(max_gal, {}).get("avg_ms", 0)
        print(f"  1:N search over {max_gal:,} identities: {gal_time:.3f} ms")
        print(f"  Estimated throughput: {1_000_000 // max(1, int(gal_time + 0.01)):,} queries/second")

    print(f"\n{separator}\n")

    return 0


if __name__ == "__main__":
    exit(main())
