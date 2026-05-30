import numpy as np
from scipy.stats import ks_2samp
import logging

class DriftDetector:
    def __init__(self, threshold=0.05):
        self.threshold = threshold
        self.reference_data = None
        logging.basicConfig(level=logging.INFO)
        self.logger = logging.getLogger("DriftDetector")

    def set_reference(self, data):
        self.reference_data = np.array(data)
        self.logger.info("ML-DRIFT: Reference baseline established.")

    def detect(self, current_data):
        if self.reference_data is None:
            raise ValueError("Reference data must be set before detection.")
        
        current_data = np.array(current_data)
        
        # Perform Kolmogorov-Smirnov test for distribution drift
        statistic, p_value = ks_2samp(self.reference_data, current_data)
        
        is_drifted = p_value < self.threshold
        
        if is_drifted:
            self.logger.warning(f"🚨 ML-DRIFT: Significant drift detected! p-value: {p_value:.4f}")
        else:
            self.logger.info(f"✅ ML-DRIFT: Distribution stable. p-value: {p_value:.4f}")
            
        return {
            "is_drifted": is_drifted,
            "p_value": p_value,
            "statistic": statistic
        }

if __name__ == "__main__":
    detector = DriftDetector()
    
    # Simulate reference and drifted data
    ref = np.random.normal(0, 1, 1000)
    drifted = np.random.normal(0.5, 1.2, 1000)
    
    detector.set_reference(ref)
    detector.detect(drifted)
