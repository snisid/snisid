"""ArcFace to ONNX export with NPU optimization for Qualcomm AI 100 Edge."""

import logging
import os
from pathlib import Path
from typing import Optional

import numpy as np
import torch

logger = logging.getLogger(__name__)


class ONNXExporter:
    """Export PyTorch ArcFace model to ONNX for NPU deployment."""

    @staticmethod
    def export_arcface(
        model_path: str,
        output_path: str,
        input_shape: tuple[int, ...] = (1, 3, 112, 112),
        dynamic_batch: bool = True,
        opset_version: int = 17,
    ) -> str:
        """Export ArcFace ResNet50 to ONNX with optional dynamic batch size.

        Args:
            model_path: Path to PyTorch checkpoint (.pt or .pth).
            output_path: Destination path for the .onnx file.
            input_shape: Static input shape (batch, channels, height, width).
            dynamic_batch: If True, the first dimension (batch) is dynamic.
            opset_version: ONNX opset version (17+ recommended for NPUs).

        Returns:
            Absolute path to the exported ONNX model.
        """
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

        # Build the same ArcFace architecture used in training
        from services.planes.ai.ai_face.model import ArcFaceModel

        model = ArcFaceModel(embedding_size=512)
        state = torch.load(model_path, map_location=device, weights_only=True)

        # Strip ArcMarginProduct weights if present in checkpoint
        if any(k.startswith("margin.") for k in state.keys()):
            state = {k: v for k, v in state.items() if not k.startswith("margin.")}

        model.load_state_dict(state, strict=False)
        model.to(device)
        model.eval()

        dummy_input = torch.randn(*input_shape, device=device)

        # Define dynamic axes for flexible batch size
        dynamic_axes: Optional[dict[str, dict[int, str]]] = None
        if dynamic_batch:
            dynamic_axes = {
                "input": {0: "batch_size"},
                "embedding": {0: "batch_size"},
            }

        output_dir = Path(output_path).parent
        output_dir.mkdir(parents=True, exist_ok=True)

        torch.onnx.export(
            model,
            dummy_input,
            output_path,
            input_names=["input"],
            output_names=["embedding"],
            dynamic_axes=dynamic_axes,
            opset_version=opset_version,
            do_constant_folding=True,
            export_params=True,
        )

        logger.info("ArcFace exported to ONNX: %s", output_path)
        return os.path.abspath(output_path)

    @staticmethod
    def validate_onnx(model_path: str, atol: float = 1e-4) -> dict:
        """Validate ONNX model output matches PyTorch output within tolerance.

        Args:
            model_path: Path to .onnx file.
            atol: Absolute tolerance for element-wise comparison.

        Returns:
            Dict with keys: pytorch_output, onnx_output, max_diff, passed.
        """
        import onnx
        import onnxruntime as ort

        # Load ONNX model and check schema
        onnx_model = onnx.load(model_path)
        onnx.checker.check_model(onnx_model)

        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

        # Build PyTorch reference model
        from services.planes.ai.ai_face.model import ArcFaceModel

        ref_model = ArcFaceModel(embedding_size=512)
        ref_model.to(device)
        ref_model.eval()

        rng = np.random.default_rng(42)
        dummy_np = rng.standard_normal((1, 3, 112, 112)).astype(np.float32)

        # PyTorch inference
        with torch.no_grad():
            pt_out = ref_model(torch.from_numpy(dummy_np).to(device))
        pt_out_np = pt_out.cpu().numpy()

        # ONNX Runtime inference
        providers = ["CPUExecutionProvider"]
        if ort.get_device() == "GPU":
            providers = ["CUDAExecutionProvider", "CPUExecutionProvider"]

        sess = ort.InferenceSession(model_path, providers=providers)
        onnx_out = sess.run(["embedding"], {"input": dummy_np})[0]

        max_diff = float(np.max(np.abs(pt_out_np - onnx_out)))
        passed = max_diff < atol

        result = {
            "pytorch_output": pt_out_np,
            "onnx_output": onnx_out,
            "max_diff": max_diff,
            "passed": passed,
        }

        status = "PASSED" if passed else "FAILED"
        logger.info("ONNX validation %s (max_diff=%.6f)", status, max_diff)

        return result

    @staticmethod
    def optimize_for_npu(
        model_path: str,
        npu: str = "qaic",
        output_path: Optional[str] = None,
    ) -> str:
        """Optimize ONNX graph for Qualcomm AI 100 (QAIC) or other NPUs.

        Performs:
          - Constant folding
          - FP16 conversion (where supported)
          - Node fusion for NPU operators
          - Redundant node elimination

        Args:
            model_path: Path to input .onnx file.
            npu: Target NPU backend ("qaic", "hailo", "intel_npu").
            output_path: Destination path; defaults to model_path with
                         _optimized suffix.

        Returns:
            Path to the optimized ONNX model.
        """
        import onnx
        from onnx import optimizer as onnx_optimizer

        model = onnx.load(model_path)

        # Available ONNX optimization passes
        passes = [
            "eliminate_deadend",
            "eliminate_identity",
            "eliminate_nop_dropout",
            "eliminate_nop_monotone_argmax",
            "eliminate_nop_pad",
            "eliminate_nop_transpose",
            "eliminate_unused_initializer",
            "extract_constant_to_initializer",
            "fuse_bn_into_conv",
            "fuse_consecutive_concats",
            "fuse_consecutive_log_softmax",
            "fuse_consecutive_reduce_unsqueeze",
            "fuse_consecutive_squeezes",
            "fuse_consecutive_transposes",
            "fuse_matmul_add_bias_into_gemm",
            "fuse_pad_into_conv",
            "fuse_relu_into_conv",
            "nop",
            "saturate_float",
        ]

        model = onnx_optimizer.optimize(model, passes)

        # Convert constants to FP16 for QAIC (reduces bandwidth)
        if npu == "qaic":
            model = ONNXExporter._convert_initializers_to_fp16(model)

        if output_path is None:
            stem = Path(model_path).stem
            output_path = str(Path(model_path).parent / f"{stem}_optimized.onnx")

        onnx.save(model, output_path)
        logger.info("ONNX model optimized for %s: %s", npu.upper(), output_path)

        return output_path

    @staticmethod
    def _convert_initializers_to_fp16(model: "onnx.ModelProto") -> "onnx.ModelProto":
        """Convert float32 initializers to float16 to reduce memory on NPU."""
        import onnx
        from onnx import numpy_helper

        for initializer in model.graph.initializer:
            if initializer.data_type == onnx.TensorProto.FLOAT:
                arr = numpy_helper.to_array(initializer)
                arr_f16 = arr.astype(np.float16)
                new_init = numpy_helper.from_array(arr_f16, initializer.name)
                initializer.CopyFrom(new_init)

        return model
