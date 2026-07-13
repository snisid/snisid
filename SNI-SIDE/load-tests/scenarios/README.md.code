# SNI-SIDE Load Testing Scenarios

## Scenario 1: Normal Operations (500 concurrent analysts)

```bash
locust -f locustfile.py --host=https://api.sniside.ht \
    --users=500 --spawn-rate=20 --run-time=1h \
    --headless --csv=sniside-normal
```

**Expected metrics:**
- P95 latency: <500ms
- P99 latency: <1s
- Error rate: <0.1%
- Throughput: ~500 req/s

## Scenario 2: Peak Crisis (2,000 concurrent users)

```bash
locust -f locustfile.py --host=https://api.sniside.ht \
    --users=2000 --spawn-rate=50 --run-time=30m \
    --headless --csv=sniside-crisis
```

**Expected metrics:**
- P95 latency: <1s
- P99 latency: <2s
- Error rate: <0.5%
- Throughput: ~2,000 req/s

## Scenario 3: ALPR Burst Ingest (10,000 reads/s)

```bash
locust -f locustfile.py --host=https://api.sniside.ht \
    --users=50 --spawn-rate=10 --run-time=10m \
    --headless --csv=sniside-alpr-burst \
    --tags alpr_bulk
```

**Expected metrics:**
- Throughput: 10,000 reads/s (50 users × 50 reads × ~4/s)
- P95 latency: <200ms
- Error rate: <0.01%

## Scenario 4: Intelligence Generation (100 concurrent)

```bash
locust -f locustfile.py --host=https://api.sniside.ht \
    --users=100 --spawn-rate=5 --run-time=15m \
    --headless --csv=sniside-graphrag \
    --tags graphrag
```

**Expected metrics:**
- P95 latency: <5s (LLM generation)
- Error rate: <1%
- Throughput: 20 reports/min

## Scenario 5: Sustained Load (8h)

```bash
locust -f locustfile.py --host=https://api.sniside.ht \
    --users=500 --spawn-rate=10 --run-time=8h \
    --headless --csv=sniside-sustained
```

**Expected metrics:**
- No memory leak
- No degradation over time
- All KPIs within Scenario 1 bounds
