import http from 'k6/http';
import { check } from 'k6';

const baseURL = "http://gocrypto:8080"
const endpoints = [
    {
        route: "/api/exchange-rates",
        method: "GET",
        load: 10000,
        expected: {
            status: 200
        }
    }
]

export function RequestEP(route, method, expected) {
    
    let res;
    if (method == "POST") res = http.post(route);
    else res = http.get(route);

    check(res, { 'status was 200': (r) => r.status == expected.status });

}

export function LoadEP(ep) {
    for (let i = 0; i < ep.load; i++) {
        RequestEP(baseURL + ep.route, ep.method, ep.expected)
    }
}

export default function () {
    endpoints.forEach(LoadEP)
}