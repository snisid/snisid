"""NPU runtime abstraction with auto-fallback to CPU/ONNX Runtime."""

import logging
import platform
import sys
from abc import ABC, abstractmethod
from typing import Optional

import numpy as np

logger = logging.getLogger(__name__)


class NPURuntime(ABC):
    """Abstract NPU runtime with fallback to CPU/ONNX Runtime.

    Factory method ``create`` auto-selects the best available backend.
    """

    def __init__(self, model_path: str, backend: str = "auto"):
        self.model_path = model_path
        self.backend = backend
        self._session = None
        self._input_name: Optional[str] = None
        self._output_name: Optional[str] = None

    @abstractmethod
    def infer(self, input_tensor: np.ndarray) -> np.ndarray:
        """Run inference and return output as numpy array."""

    @property
    @abstractmethod
    def available_backends(self) -> list[str]:
        """Return available execution providers for this platform."""

    @classmethod
    def create(
        cls,
        model_path: str,
        preferred_backend: str = "auto",
    ) -> "NPURuntime":
        """Factory: auto-detect best backend and return appropriate runtime.

        Priority: QAIC (Qualcomm) > CUDA > CPU
        """
        if preferred_backend == "auto":
            # Probe available providers in priority order
            for backend_name, runtime_cls in [
                ("qaic", QaicNPURuntime),
                ("cuda", OnnxRuntimeCUDA),
                ("cpu", OnnxRuntimeCPU),
            ]:
                try:
                    instance: NPURuntime = runtime_cls(model_path)
                    _ = instance.available_backends
                    logger.info("Selected backend: %s", backend_name)
                    return instance
                except Exception as exc:
                    logger.debug("Backend %s unavailable: %s", backend_name, exc)
                    continue

            raise RuntimeError("No inference backend available on this platform.")

        backend_map = {
            "qaic": QaicNPURuntime,
            "cuda": OnnxRuntimeCUDA,
            "cpu": OnnxRuntimeCPU,
        }
        if preferred_backend not in backend_map:
            raise ValueError(
                f"Unknown backend '{preferred_backend}'. "
                f"Choose from: {list(backend_map)}"
            )

        instance = backend_map[preferred_backend](model_path)
        _ = instance.available_backends
        return instance


class QaicNPURuntime(NPURuntime):
    """Qualcomm AI 100 specific runtime using QAIC SDK.

    Requires ``qaic-api`` package and the QAIC runtime driver installed.
    Falls back to ONNX Runtime CPU if QAIC is unavailable.
    """

    def __init__(self, model_path: str):
        super().__init__(model_path, backend="qaic")
        self._load_session()

    def _load_session(self):
        try:
            import onnxruntime as ort
        except ImportError:
            raise RuntimeError("onnxruntime is required for QAIC runtime.")

        qaic_provider = "QNNExecutionProvider"
        available = ort.get_available_providers()

        if qaic_provider not in available:
            logger.warning(
                "QAIC provider not found in onnxruntime. "
                "Available: %s. Falling back to CPU.",
                available,
            )
            self._session = ort.InferenceSession(
                self.model_path, providers=["CPUExecutionProvider"]
            )
        else:
            self._session = ort.InferenceSession(
                self.model_path,
                providers=[
                    (qaic_provider, {"backend_path": "QnnBackend.dll"}),
                    "CPUExecutionProvider",
                ],
            )

        self._input_name = self._session.get_inputs()[0].name
        self._output_name = self._session.get_outputs()[0].name
        logger.info("QAIC NPU session loaded for %s", self.model_path)

    def infer(self, input_tensor: np.ndarray) -> np.ndarray:
        if self._session is None:
            raise RuntimeError("QAIC session not initialized.")
        return self._session.run(
            [self._output_name], {self._input_name: input_tensor}
        )[0]

    @property
    def available_backends(self) -> list[str]:
        try:
            import onnxruntime as ort

            return ort.get_available_providers()
        except ImportError:
            return []


class OnnxRuntimeCPU(NPURuntime):
    """CPU fallback using ONNX Runtime."""

    def __init__(self, model_path: str):
        super().__init__(model_path, backend="cpu")
        self._load_session()

    def _load_session(self):
        try:
            import onnxruntime as ort
        except ImportError:
            raise RuntimeError("onnxruntime is required for CPU runtime.")

        sess_options = ort.SessionOptions()
        sess_options.graph_optimization_level = ort.GraphOptimizationLevel.ORT_ENABLE_ALL
        sess_options.intra_op_num_threads = max(1, (os_cpu_count() or 4))

        self._session = ort.InferenceSession(
            self.model_path,
            sess_options=sess_options,
            providers=["CPUExecutionProvider"],
        )
        self._input_name = self._session.get_inputs()[0].name
        self._output_name = self._session.get_outputs()[0].name
        logger.info("ONNX Runtime CPU session loaded for %s", self.model_path)

    def infer(self, input_tensor: np.ndarray) -> np.ndarray:
        if self._session is None:
            raise RuntimeError("CPU session not initialized.")
        return self._session.run(
            [self._output_name], {self._input_name: input_tensor}
        )[0]

    @property
    def available_backends(self) -> list[str]:
        try:
            import onnxruntime as ort

            return ort.get_available_providers()
        except ImportError:
            return []


class OnnxRuntimeCUDA(NPURuntime):
    """CUDA fallback using ONNX Runtime with CUDA provider."""

    def __init__(self, model_path: str):
        super().__init__(model_path, backend="cuda")
        self._load_session()

    def _load_session(self):
        try:
            import onnxruntime as ort
        except ImportError:
            raise RuntimeError("onnxruntime is required for CUDA runtime.")

        cuda_provider = "CUDAExecutionProvider"
        available = ort.get_available_providers()

        if cuda_provider not in available:
            logger.warning(
                "CUDA provider not found. Available: %s. Falling back to CPU.",
                available,
            )
            providers = ["CPUExecutionProvider"]
        else:
            providers = [
                (cuda_provider, {"device_id": 0, "arena_extend_strategy": "kNextPowerOfTwo"}),
                "CPUExecutionProvider",
            ]

        self._session = ort.InferenceSession(
            self.model_path, providers=providers
        )
        self._input_name = self._session.get_inputs()[0].name
        self._output_name = self._session.get_outputs()[0].name
        logger.info("ONNX Runtime CUDA session loaded for %s", self.model_path)

    def infer(self, input_tensor: np.ndarray) -> np.ndarray:
        if self._session is None:
            raise RuntimeError("CUDA session not initialized.")
        return self._session.run(
            [self._output_name], {self._input_name: input_tensor}
        )[0]

    @property
    def available_backends(self) -> list[str]:
        try:
            import onnxruntime as ort

            return ort.get_available_providers()
        except ImportError:
            return []


def os_cpu_count() -> int:
    """Cross-platform logical CPU count."""
    try:
        if platform.system() == "Windows":
            import os as _os

            return int(_os.environ.get("NUMBER_OF_PROCESSORS", 4))
        return len(open("/proc/stat").readlines()) if sys.platform != "win32" else 4
    except Exception:
        return 4
