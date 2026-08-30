import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 }, // Ramp up to 50 users
    { duration: '1m', target: 50 },  // Stay at 50 users
    { duration: '30s', target: 0 },  // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(99)<500'], // 99% of requests must be below 500ms
  },
};

export default function () {
  const url = 'http://localhost/api/v1/identities';
  const payload = JSON.stringify({
    firstName: 'Test',
    lastName: 'User',
    dob: '2000-01-01',
    agency: 'AGENCY-TEST',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer YOUR_TEST_TOKEN',
    },
  };

  const res = http.post(url, payload, params);
  check(res, {
    'is status 201': (r) => r.status === 201,
  });

  sleep(1);
}
