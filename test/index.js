import http from 'k6/http';
import { check } from 'k6';
import { SharedArray } from 'k6/data';
import papaparse from './papaparse.min.js';

const data = new SharedArray('Test texts', function () {
    return papaparse.parse(open('./data.csv'), { header: true }).data;
});

const serviceUrl = 'http://localhost:8081/similarity';

export const options = {
    stages: [
        { duration: '30s', target: 100 },
        { duration: '1m', target: 100 },
        { duration: '30s', target: 200 },
        { duration: '1m', target: 200 },
        { duration: '30s', target: 400 },
        { duration: '1m', target: 400 },
        { duration: '30s', target: 800 },
        { duration: '1m', target: 800 },
        { duration: '30s', target: 1600 },
        { duration: '1m', target: 1600 },
        { duration: '30s', target: 3200 },
        { duration: '1m', target: 3200 },
        { duration: '30s', target: 5000 },
        { duration: '1m', target: 5000 },
        { duration: '2m', target: 0 },
    ],
    thresholds: {
        "http_req_failed": ["rate<0.01"],
        "http_req_duration": ["p(95)<1500"],
    },
};

export default function () {
    const text1 = data[Math.floor(Math.random() * data.length)].Text;
    const text2 = data[Math.floor(Math.random() * data.length)].Text;

    const payload = JSON.stringify({
        text1: text1,
        text2: text2,
    });

    const params = {
        headers: { 'Content-Type': 'application/json' },
    };

    const response = http.post(serviceUrl, payload, {
        ...params,
        tags: { service: 'Go' },
    });
    check(response, {
        OK: (r) => r.status === 200,
        'Similarity Returned': (r) =>
            typeof JSON.parse(r.body).similarity === 'number',
    });
}
