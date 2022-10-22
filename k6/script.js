import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  http.get('http://localhost/test/prime?min=1&max=10000');
  sleep(1)
}
