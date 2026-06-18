import importlib
import sys
import pytest
import numpy as np
from unittest.mock import patch, MagicMock


@pytest.fixture(autouse=True)
def _mock_onnx_and_reload():
    mock_ort = MagicMock()
    mock_ort.SessionOptions.return_value = MagicMock()
    mock_ort.GraphOptimizationLevel = MagicMock()
    mock_ort.get_available_providers.return_value = [
        "CPUExecutionProvider",
        "QNNExecutionProvider",
    ]
    mock_session = MagicMock()
    mock_session.get_inputs.return_value = [MagicMock(name="input")]
    mock_session.get_outputs.return_value = [MagicMock(name="output")]
    mock_session.run.return_value = [np.random.rand(1, 512).astype(np.float32)]
    mock_ort.InferenceSession.return_value = mock_session
    mock_ort.__version__ = "1.15.0"

    # Ensure onnxruntime mock is in sys.modules before module loads
    sys.modules["onnxruntime"] = mock_ort

    # Force reimport of the module under test
    for mod in list(sys.modules.keys()):
        if "npu_runtime" in mod and mod != __name__:
            del sys.modules[mod]

    yield

    # Cleanup: remove the module from cache so next test gets fresh import
    for mod in list(sys.modules.keys()):
        if "npu_runtime" in mod and mod != __name__:
            del sys.modules[mod]
    if "onnxruntime" in sys.modules:
        del sys.modules["onnxruntime"]


class TestNPURuntimeFactory:
    def test_create_auto_selects_first_available(self):
        from services.biometric.inference.npu_runtime import (
            NPURuntime,
            QaicNPURuntime,
        )

        with patch.object(
            QaicNPURuntime, "available_backends", new_callable=MagicMock
        ) as mock_backends:
            mock_backends.return_value = [
                "QNNExecutionProvider",
                "CPUExecutionProvider",
            ]
            instance = NPURuntime.create("model.onnx", preferred_backend="auto")
            assert isinstance(instance, QaicNPURuntime)

    def test_create_cpu_fallback_when_qaic_unavailable(self):
        from services.biometric.inference.npu_runtime import (
            NPURuntime,
            QaicNPURuntime,
            OnnxRuntimeCUDA,
            OnnxRuntimeCPU,
        )

        with patch.object(
            QaicNPURuntime, "available_backends", new_callable=MagicMock
        ) as qaic_b:
            qaic_b.side_effect = RuntimeError("QAIC unavailable")
            with patch.object(
                OnnxRuntimeCUDA, "available_backends", new_callable=MagicMock
            ) as cuda_b:
                cuda_b.side_effect = RuntimeError("CUDA unavailable")
                with patch.object(
                    OnnxRuntimeCPU, "available_backends", new_callable=MagicMock
                ) as cpu_b:
                    cpu_b.return_value = ["CPUExecutionProvider"]
                    instance = NPURuntime.create(
                        "model.onnx", preferred_backend="auto"
                    )
                    assert isinstance(instance, OnnxRuntimeCPU)

    def test_create_no_backend_raises(self):
        from services.biometric.inference.npu_runtime import (
            NPURuntime,
            QaicNPURuntime,
            OnnxRuntimeCUDA,
            OnnxRuntimeCPU,
        )

        with patch.object(
            QaicNPURuntime, "available_backends", new_callable=MagicMock
        ) as qaic_b:
            qaic_b.side_effect = RuntimeError("QAIC unavailable")
            with patch.object(
                OnnxRuntimeCUDA, "available_backends", new_callable=MagicMock
            ) as cuda_b:
                cuda_b.side_effect = RuntimeError("CUDA unavailable")
                with patch.object(
                    OnnxRuntimeCPU, "available_backends", new_callable=MagicMock
                ) as cpu_b:
                    cpu_b.side_effect = RuntimeError("CPU unavailable")
                    with pytest.raises(
                        RuntimeError, match="No inference backend available"
                    ):
                        NPURuntime.create("model.onnx")

    def test_create_with_explicit_cpu_backend(self):
        from services.biometric.inference.npu_runtime import (
            NPURuntime,
            OnnxRuntimeCPU,
        )

        with patch.object(
            OnnxRuntimeCPU, "available_backends", new_callable=MagicMock
        ) as cpu_b:
            cpu_b.return_value = ["CPUExecutionProvider"]
            instance = NPURuntime.create(
                "model.onnx", preferred_backend="cpu"
            )
            assert isinstance(instance, OnnxRuntimeCPU)

    def test_create_unknown_backend_raises(self):
        from services.biometric.inference.npu_runtime import NPURuntime

        with pytest.raises(ValueError, match="Unknown backend"):
            NPURuntime.create("model.onnx", preferred_backend="unknown")


class TestOnnxRuntimeCPU:
    def test_infer_returns_numpy_array(self):
        from services.biometric.inference.npu_runtime import OnnxRuntimeCPU

        runtime = OnnxRuntimeCPU("model.onnx")
        result = runtime.infer(np.random.rand(1, 3, 112, 112).astype(np.float32))
        assert isinstance(result, np.ndarray)
        assert result.shape[1] == 512

    def test_infer_raises_when_session_none(self):
        from services.biometric.inference.npu_runtime import OnnxRuntimeCPU

        runtime = OnnxRuntimeCPU("model.onnx")
        runtime._session = None
        with pytest.raises(RuntimeError, match="CPU session not initialized"):
            runtime.infer(np.random.rand(1, 3, 112, 112).astype(np.float32))

    def test_available_backends_returns_list(self):
        from services.biometric.inference.npu_runtime import OnnxRuntimeCPU

        backends = OnnxRuntimeCPU("model.onnx").available_backends
        assert isinstance(backends, list)

    def test_infer_passes_input_to_session(self):
        from services.biometric.inference.npu_runtime import OnnxRuntimeCPU

        runtime = OnnxRuntimeCPU("model.onnx")
        input_tensor = np.random.rand(1, 3, 112, 112).astype(np.float32)
        runtime.infer(input_tensor)
        runtime._session.run.assert_called_once()

    def test_session_options_set(self):
        from services.biometric.inference.npu_runtime import (
            OnnxRuntimeCPU,
            ort,
        )

        runtime = OnnxRuntimeCPU("model.onnx")
        ort.InferenceSession.assert_called_once()
        _, kwargs = ort.InferenceSession.call_args
        assert "providers" in kwargs
        assert "CPUExecutionProvider" in kwargs["providers"]


class TestQaicNPURuntime:
    def test_qaic_creates_session(self):
        from services.biometric.inference.npu_runtime import (
            QaicNPURuntime,
            ort,
        )

        runtime = QaicNPURuntime("model.onnx")
        assert runtime._session is not None
        ort.InferenceSession.assert_called_once()

    def test_qaic_infer_returns_embedding(self):
        from services.biometric.inference.npu_runtime import QaicNPURuntime

        runtime = QaicNPURuntime("model.onnx")
        result = runtime.infer(np.random.rand(1, 3, 112, 112).astype(np.float32))
        assert result.shape[1] == 512


class TestOnnxRuntimeCUDA:
    def test_cuda_creates_session(self):
        from services.biometric.inference.npu_runtime import (
            OnnxRuntimeCUDA,
            ort,
        )

        runtime = OnnxRuntimeCUDA("model.onnx")
        assert runtime._session is not None
        ort.InferenceSession.assert_called_once()

    def test_cuda_infer_returns_embedding(self):
        from services.biometric.inference.npu_runtime import OnnxRuntimeCUDA

        runtime = OnnxRuntimeCUDA("model.onnx")
        result = runtime.infer(np.random.rand(1, 3, 112, 112).astype(np.float32))
        assert isinstance(result, np.ndarray)
