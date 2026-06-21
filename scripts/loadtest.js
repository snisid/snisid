// SNISID — k6 Load Test
// Usage: k6 run scripts/loadtest.js
//        k6 run -e PORT_OFFSET=90000 scripts/loadtest.js  (for e2e override ports)

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const PORT_OFFSET = parseInt(__ENV.PORT_OFFSET) || 0;

const HT_SERVICES = [
  { name: 'id-core',              port: 8201 },
  { name: 'civil-ht',             port: 8202 },
  { name: 'bio-ht',               port: 8203 },
  { name: 'card-ht',              port: 8204 },
  { name: 'pki-ht',               port: 8206 },
  { name: 'iam-ht',               port: 8207 },
  { name: 'interop-ht',           port: 8208 },
  { name: 'infra-ht',             port: 8209 },
  { name: 'cyber-ht',             port: 8210 },
  { name: 'offline-ht',           port: 8211 },
  { name: 'field-ht',             port: 8212 },
  { name: 'data-ht',              port: 8213 },
  { name: 'api-ht',               port: 8214 },
  { name: 'foves-ht',             port: 8215 },
  { name: 'lapi-ht',              port: 8216 },
  { name: 'fpr-ht',               port: 8205 },
  { name: 'sigint-ht',            port: 8301 },
  { name: 'humint-ht',            port: 8302 },
  { name: 'air-defense-ht',       port: 8303 },
  { name: 'mil-c2-ht',            port: 8304 },
  { name: 'bio-surveillance-ht',  port: 8305 },
  { name: 'executive-protection-ht', port: 8306 },
  { name: 'transport-security-ht',   port: 8307 },
  { name: 'radiation-safety-svc',     port: 8308 },
  { name: 'all-source-fusion-ht',     port: 8309 },
  { name: 'counterintel-ht',          port: 8310 },
  { name: 'critical-infra-protection-ht', port: 8311 },
  { name: 'fisa-court-svc',       port: 8312 },
  { name: 'classification-mgmt-ht',   port: 8313 },
];

const POST_SERVICES = [
  { name: 'id-core',  port: 8201,  endpoint: '/api/v1/citizens',       body: { firstName: 'Load', lastName: 'Test', dob: '1990-01-01', nationalId: 'LOAD-0001' } },
  { name: 'fpr-ht',   port: 8205,  endpoint: '/api/v1/warrants',       body: { subjectId: 'LOAD-0001', warrantType: 'Test', issuingAuthority: 'k6', reason: 'Load test' } },
  { name: 'sigint-ht', port: 8301, endpoint: '/api/v1/interception-targets', body: { fisaWarrantId: 'LOAD-FISA', targetIdentifier: 'load-test', commType: 'Test', authorizee: 'k6' } },
  { name: 'civil-ht', port: 8202,  endpoint: '/api/v1/birth-records',   body: { citizenId: 'LOAD-0001', registrant: 'Load Test', dateOfBirth: '1990-01-01', placeOfBirth: 'Test' } },
  { name: 'bio-surveillance-ht', port: 8305, endpoint: '/api/v1/outbreaks', body: { disease: 'Test', region: 'Load', confirmedCases: 1, suspectedCases: 1, severity: 'TEST', reportedBy: 'k6' } },
];

const errorRate = new Rate('errors');
const healthzTrend = new Trend('healthz_duration');

export let options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '30s', target: 25 },
    { duration: '30s', target: 50 },
    { duration: '30s', target: 25 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    errors: ['rate<0.01'],
  },
};

export default function () {
  // Test all /healthz endpoints
  for (const svc of HT_SERVICES) {
    const url = `http://localhost:${svc.port + PORT_OFFSET}/healthz`;
    const start = Date.now();
    const res = http.get(url);
    const duration = Date.now() - start;

    healthzTrend.add(duration);
    const ok = check(res, {
      [`${svc.name} healthz status 200`]: (r) => r.status === 200,
    });
    if (!ok) errorRate.add(1);

    sleep(0.1);
  }

  // Test POST endpoints on key services
  for (const svc of POST_SERVICES) {
    const url = `http://localhost:${svc.port + PORT_OFFSET}${svc.endpoint}`;
    const payload = JSON.stringify(svc.body);
    const params = { headers: { 'Content-Type': 'application/json' } };
    const res = http.post(url, payload, params);

    const ok = check(res, {
      [`${svc.name} POST ${svc.endpoint} 2xx`]: (r) => r.status >= 200 && r.status < 300,
    });
    if (!ok) errorRate.add(1);

    sleep(0.2);
  }
}
