# PROMPT 287: AUTOMATED CAPACITY PLANNING & RIGHTSIZING

This architecture defines the continuous resource optimization and predictive scaling strategy for the SNISID platform, ensuring that national compute and storage investments are utilized with maximum efficiency.

---

## 1. Planning Architecture (Continuous Optimization)

SNISID utilizes a multi-tier optimization stack to balance performance and cost.

- **Collector (Prometheus/Thanos)**: Provides historical usage data (CPU, RAM, Storage, Network) for every workload.
- **Analyzer (Vertical Pod Autoscaler - VPA)**: Analyzes historical metrics to recommend or automatically apply the ideal `requests` and `limits` for pods.
- **Forecaster (Prophet/Custom AI)**: Predictive engine that identifies growth trends and forecasts hardware exhaustion dates.
- **FinOps Engine (Kubecost)**: Provides real-time visibility into the cost of every namespace, service, and agency.

---

## 2. Optimization Workflows (Rightsizing)

1.  **Metric Aggregation**: The system aggregates P95 and P99 resource usage over a rolling 30-day window.
2.  **Recommendation**: VPA generates "Rightsizing Recommendations" (e.g., "This service is only using 20% of its requested RAM; reduce request by 500MB").
3.  **Automated Application**: For non-critical workloads (Sandbox/Dev), the system automatically applies recommendations via GitOps.
4.  **Human Review**: For mission-critical workloads, recommendations are presented in the **Developer Portal** (Prompt 282) for manual approval and one-click application.

---

## 3. Capacity Planning (Forecasting)

- **Regional Growth Analysis**: The system monitors regional cluster utilization and predicts when a region will reach its "Safe Capacity" (80%).
- **Automated Hardware Requests**: When a capacity threshold is crossed, the system automatically generates a "National Hardware Procurement" request with documented evidence of growth.
- **Seasonal Modeling**: AI models national events (e.g., elections, censuses, or security crises) to preemptively scale the infrastructure before the surge occurs.

---

## 4. Analysis & Reporting

- **Waste Dashboard**: Highlighting "Zombie Workloads" and "Over-Provisioned Services" to encourage agency teams to optimize.
- **Efficiency Grading**: Every agency is assigned a "Resource Efficiency Grade" based on how well their actual usage aligns with their requested quotas.
- **National Cost Attribution**: Automated monthly reports showing exactly how much of the national cloud budget was consumed by each agency and project.

---

## 5. Governance Strategy

- **Incentivized Optimization**: Agencies that maintain a high efficiency grade receive priority for new hardware allocations.
- **Audit Ledger**: All rightsizing changes and capacity forecasts are recorded in the forensic ledger for multi-year resource planning.
- **Quota Rebalancing**: The system automatically suggests reallocating unused quotas from stagnant projects to rapidly growing ones, maximizing the value of existing sovereign hardware.

---

**PROMPT 287 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 288 — INFRASTRUCTURE DRIFT AUTO-CORRECTION.**
