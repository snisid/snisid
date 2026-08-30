class VeraStrategist:
    """
    Vera: The Strategic Reasoning Agent.
    Role: Decide policy strategy based on global world state.
    """
    def __init__(self, risk_threshold=0.7):
        self.risk_threshold = risk_threshold

    def evaluate_strategy(self, world_state):
        print(f"VERA: Analyzing world state: {world_state}")
        
        if world_state["fraud_rate"] > self.risk_threshold:
            return "increase_enforcement"
        
        if world_state["economic_load"] > 0.9:
            return "relax_policy"
            
        return "maintain_stability"

if __name__ == "__main__":
    vera = VeraStrategist()
    strategy = vera.evaluate_strategy({"fraud_rate": 0.82, "economic_load": 0.4})
    print(f"VERA STRATEGY: {strategy}")
