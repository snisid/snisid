import numpy as np

class BehavioralProfiler:
    """
    Pattern-of-Life modeling.
    Enriches intelligence by detecting subtle deviations from baseline behavior.
    """
    def __init__(self):
        self.profiles = {} # citizen_id -> profile

    def update_profile(self, citizen_id, transaction_val, location):
        if citizen_id not in self.profiles:
            self.profiles[citizen_id] = {"avg": transaction_val, "locations": [location]}
        else:
            self.profiles[citizen_id]["avg"] = (self.profiles[citizen_id]["avg"] + transaction_val) / 2
            if location not in self.profiles[citizen_id]["locations"]:
                self.profiles[citizen_id]["locations"].append(location)

    def is_anomalous(self, citizen_id, current_val, current_loc):
        profile = self.profiles.get(citizen_id)
        if not profile:
            return False
        
        # Simple heuristic: 3x average or new location
        if current_val > profile["avg"] * 3:
            return True
        if current_loc not in profile["locations"]:
            return True
            
        return False

if __name__ == "__main__":
    profiler = BehavioralProfiler()
    profiler.update_profile("CIT_001", 200, "Port-au-Prince")
    print(f"Anomaly: {profiler.is_anomalous('CIT_001', 10000, 'Paris')}")
