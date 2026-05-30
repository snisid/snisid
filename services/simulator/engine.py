import time
import random

class CyberCivilizationSimulator:
    """
    Closed-Loop Digital Twin Engine
    Simulates: Cyber Economy, Internet Warfare, and Infrastructure State.
    """
    def __init__(self):
        self.world_state = {
            "cyber_gdp": 1000.0,
            "attack_pressure": 0.1,
            "infra_health": 1.0,
            "agent_efficiency": 0.85
        }

    def simulate_economy(self):
        # Compute economy grows with infrastructure health and agent efficiency
        growth = (self.world_state["infra_health"] * self.world_state["agent_efficiency"]) * 0.05
        self.world_state["cyber_gdp"] += growth
        return self.world_state["cyber_gdp"]

    def simulate_warfare(self):
        # Attack pressure increases randomly, but is mitigated by agent efficiency
        new_threats = random.uniform(0, 0.2)
        mitigation = self.world_state["agent_efficiency"] * 0.1
        self.world_state["attack_pressure"] = max(0, self.world_state["attack_pressure"] + new_threats - mitigation)
        return self.world_state["attack_pressure"]

    def simulate_infra(self):
        # Infrastructure health degrades with attack pressure
        degredation = self.world_state["attack_pressure"] * 0.05
        self.world_state["infra_health"] = max(0, self.world_state["infra_health"] - degredation)
        return self.world_state["infra_health"]

    def step(self):
        gdp = self.simulate_economy()
        pressure = self.simulate_warfare()
        health = self.simulate_infra()
        
        return {
            "timestamp": time.time(),
            "cyber_gdp": gdp,
            "attack_pressure": pressure,
            "infra_health": health,
            "stability_score": (health * 0.7) + (1 - pressure) * 0.3
        }

if __name__ == "__main__":
    sim = CyberCivilizationSimulator()
    for i in range(10):
        print(f"Step {i}: {sim.step()}")
        time.sleep(1)
