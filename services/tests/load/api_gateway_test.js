// api_gateway_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 1000 }, // Ramp up to 1k users
    { duration: '3m', target: 5000 }, // Ramp up to 5k users
    { duration: '1m', target: 10000 }, // Peak load 10k users
    { duration: '2m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'], // 95% of requests must be under 200ms
    http_req_failed: ['rate<0.01'],    // Less than 1% failure rate
  },
};

export default function () {
  const url = 'http://api-gateway.nexus.svc/v1/identity/verify';
  const payload = JSON.stringify({
    citizen_id: 'citizen_test_123',
    biometric_token: 'valid_token_abc'
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': 'sovereign_production_key_xyz'
    },
  };

  const res = http.post(url, payload, params);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'transaction matches': (r) => r.json().risk_level !== undefined,
  });

  sleep(1);
}
