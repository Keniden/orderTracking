import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const url = 'http://localhost:8000/api/lists/add_order';
  const payload = JSON.stringify({
    title: 'Test Order', description: 'test', price: 100,
    from_location: 'WH1', to_location: 'HUB1'
  });
  const params = { headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer test' } };
  const res = http.post(url, payload, params);
  check(res, { 'status is 200': (r) => r.status === 200 });
  sleep(1);
}

