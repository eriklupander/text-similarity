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
        { duration: '30s', target: 250 }, 
        { duration: '1m', target: 250 },
        { duration: '30s', target: 500 },
        { duration: '1m', target: 500 },
        { duration: '30s', target: 1000 },
        { duration: '1m', target: 1000 },
        { duration: '30s', target: 2500 },
        { duration: '1m', target: 2500 },
        { duration: '30s', target: 5000 },
        { duration: '1m', target: 5000 },
        { duration: '30s', target: 10000 },
        { duration: '1m', target: 10000 },
    ],
    thresholds: {
        "http_req_failed": [{
            "threshold": "rate<0.01",
            abortOnFail: true,
            delayAbortEval: '10s',
        }],
        "http_req_duration": [{
            "threshold": "p(95)<1500",
            abortOnFail: true,
            delayAbortEval: '10s',
        }],
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
