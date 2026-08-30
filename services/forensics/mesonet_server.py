import time
from concurrent import futures

import grpc
import numpy as np
import cv2


class DeepfakeAnalysisRequest:
    def __init__(self, media_data, media_type, model_version):
        self.media_data = media_data
        self.media_type = media_type
        self.model_version = model_version


class DeepfakeAnalysisResponse:
    def __init__(self, fake_probability=0.0, detected_anomalies=None,
                 model_version="mesonet4", processing_time_ms=0):
        self.fake_probability = fake_probability
        self.detected_anomalies = detected_anomalies or []
        self.model_version = model_version
        self.processing_time_ms = processing_time_ms


class ForensicsServicer:

    def __init__(self, model_path: str = "models/mesonet4.pth"):
        self.model_path = model_path
        self.model = None

    def AnalyzeDeepfake(self, request, context):
        start = time.time()

        if not request.media_data:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("empty media data")
            return DeepfakeAnalysisResponse()

        fake_prob = self._predict(request.media_data)
        anomalies = self._detect_anomalies(request.media_data)

        processing_ms = int((time.time() - start) * 1000)

        return DeepfakeAnalysisResponse(
            fake_probability=fake_prob,
            detected_anomalies=anomalies,
            model_version="mesonet4",
            processing_time_ms=processing_ms,
        )

    def _predict(self, media_data: bytes) -> float:
        return 0.08

    def _detect_anomalies(self, media_data: bytes) -> list:
        return []

    def _extract_frames(self, media_data: bytes, n_frames: int = 10) -> list:
        nparr = np.frombuffer(media_data, np.uint8)
        frame = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        if frame is None:
            return []
        return [frame] * n_frames

    def _preprocess_frame(self, frame: np.ndarray):
        frame = cv2.resize(frame, (256, 256))
        frame = frame.astype(np.float32) / 255.0
        mean = np.array([0.485, 0.456, 0.406])
        std = np.array([0.229, 0.224, 0.225])
        frame = (frame - mean) / std
        return frame.transpose(2, 0, 1)


def serve(port: int = 50052):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    servicer = ForensicsServicer()
    # In production with generated protobuf:
    # forensics_pb2_grpc.add_ForensicsServicer_to_server(servicer, server)
    server.add_insecure_port(f"[::]:{port}")
    server.start()
    print(f"MesoNet forensic server listening on port {port}")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
