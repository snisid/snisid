import os
import sys
import pytest
from pathlib import Path
from unittest.mock import patch, MagicMock


@pytest.fixture(autouse=True)
def _mock_deps_and_reload():
    """Mock torch, ArcFaceModel, onnx, onnxruntime before importing module."""
    mock_torch = MagicMock()
    mock_torch.cuda.is_available.return_value = False
    mock_torch.device.return_value = "cpu"
    mock_torch.randn.return_value = MagicMock()
    mock_torch.no_grad.return_value.__enter__.return_value = None
    mock_torch.onnx.export = MagicMock()

    mock_onnx = MagicMock()
    mock_onnx.TensorProto.FLOAT = 1
    mock_onnx.load.return_value = MagicMock()
    mock_onnx.checker.check_model = MagicMock()

    mock_onnxruntime = MagicMock()
    mock_onnxruntime.get_device.return_value = "CPU"
    mock_onnxruntime.InferenceSession.return_value.run.return_value = [
        MagicMock()
    ]

    sys.modules["torch"] = mock_torch
    sys.modules["onnx"] = mock_onnx
    sys.modules["onnxruntime"] = mock_onnxruntime

    # Clear cached module under test
    for mod in list(sys.modules.keys()):
        if "onnx_export" in mod and mod != __name__:
            del sys.modules[mod]
    if "services.planes" in mod for mod in sys.modules:
        pass
    # Also clear ArcFaceModel's parent module if cached
    for mod in list(sys.modules.keys()):
        if "ai_face" in mod or "planes" in mod:
            del sys.modules[mod]

    yield

    for mod in list(sys.modules.keys()):
        if "onnx_export" in mod and mod != __name__:
            del sys.modules[mod]
    for pkg in ["torch", "onnx", "onnxruntime"]:
        if pkg in sys.modules:
            del sys.modules[pkg]


class TestONNXExporter:
    def test_input_shape_validation_standard(self):
        from services.biometric.inference.onnx_export import ONNXExporter

        mock_model = MagicMock()
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=mock_model,
        ):
            path = ONNXExporter.export_arcface(
                "dummy.pt", "/tmp/test.onnx"
            )
            assert isinstance(path, str)

    def test_dynamic_batch_enabled(self):
        from services.biometric.inference.onnx_export import ONNXExporter

        mock_model = MagicMock()
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=mock_model,
        ):
            ONNXExporter.export_arcface(
                "dummy.pt",
                "/tmp/test.onnx",
                dynamic_batch=True,
            )
            export_call = sys.modules["torch"].onnx.export
            export_call.assert_called_once()
            _, kwargs = export_call.call_args
            assert kwargs["dynamic_axes"] is not None
            assert "batch_size" in kwargs["dynamic_axes"]["input"].values()

    def test_dynamic_batch_disabled(self):
        from services.biometric.inference.onnx_export import ONNXExporter

        mock_model = MagicMock()
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=mock_model,
        ):
            ONNXExporter.export_arcface(
                "dummy.pt",
                "/tmp/test.onnx",
                dynamic_batch=False,
            )
            export_call = sys.modules["torch"].onnx.export
            export_call.assert_called_once()
            _, kwargs = export_call.call_args
            assert kwargs["dynamic_axes"] is None

    def test_opset_version_passed_correctly(self):
        from services.biometric.inference.onnx_export import ONNXExporter

        mock_model = MagicMock()
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=mock_model,
        ):
            ONNXExporter.export_arcface(
                "dummy.pt",
                "/tmp/test.onnx",
                opset_version=17,
            )
            export_call = sys.modules["torch"].onnx.export
            _, kwargs = export_call.call_args
            assert kwargs["opset_version"] == 17

    def test_opset_version_alternative(self):
        from services.biometric.inference.onnx_export import ONNXExporter

        mock_model = MagicMock()
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=mock_model,
        ):
            ONNXExporter.export_arcface(
                "dummy.pt",
                "/tmp/test.onnx",
                opset_version=19,
            )
            export_call = sys.modules["torch"].onnx.export
            _, kwargs = export_call.call_args
            assert kwargs["opset_version"] == 19

    def test_output_path_creation(self, tmp_path):
        from services.biometric.inference.onnx_export import ONNXExporter

        output = str(tmp_path / "models" / "arcface.onnx")
        mock_model = MagicMock()
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=mock_model,
        ):
            result = ONNXExporter.export_arcface("dummy.pt", output)
            assert Path(output).parent.exists()
            assert result == str(Path(output).resolve())

    def test_error_non_existent_model_path(self):
        from services.biometric.inference.onnx_export import ONNXExporter

        mock_torch = sys.modules["torch"]
        mock_torch.load.side_effect = FileNotFoundError("No such file")
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=MagicMock(),
        ):
            with pytest.raises(FileNotFoundError):
                ONNXExporter.export_arcface(
                    "/nonexistent/model.pt", "/tmp/test.onnx"
                )

    def test_state_dict_strips_margin_keys(self):
        from services.biometric.inference.onnx_export import ONNXExporter

        mock_state = {
            "margin.weight": MagicMock(),
            "backbone.conv1.weight": MagicMock(),
            "backbone.conv1.bias": MagicMock(),
        }
        mock_torch = sys.modules["torch"]
        mock_torch.load.return_value = mock_state

        mock_model = MagicMock()
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=mock_model,
        ):
            ONNXExporter.export_arcface("dummy.pt", "/tmp/test.onnx")
            loaded_state = mock_model.load_state_dict.call_args[0][0]
            assert "margin.weight" not in loaded_state
            assert "backbone.conv1.weight" in loaded_state

    def test_export_uses_cpu_when_cuda_unavailable(self):
        from services.biometric.inference.onnx_export import ONNXExporter

        mock_torch = sys.modules["torch"]
        mock_torch.cuda.is_available.return_value = False

        mock_model = MagicMock()
        with patch(
            "services.biometric.inference.onnx_export.ArcFaceModel",
            return_value=mock_model,
        ):
            ONNXExporter.export_arcface("dummy.pt", "/tmp/test.onnx")
            assert mock_torch.device.call_args[0][0] == "cpu"
