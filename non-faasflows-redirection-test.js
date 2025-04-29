import http from 'k6/http';
import { check, sleep } from 'k6';


export const options = {
    stages: [
        { duration: '20s', target: 10 },  // First stage: 10 users for 20 seconds
        { duration: '40s', target: 50 },  // Second stage: 50 users for 40 seconds
        { duration: '90s', target: 200 }, // Third stage: 200 users for 90 seconds
    ],
};

export default function () {
    const validUserIds = [10, 20, 30, 40];
    const userId = validUserIds[Math.floor(Math.random() * validUserIds.length)];
    const origin = `${(Math.random() * 90).toFixed(6)},${(Math.random() * 90).toFixed(6)}`;

    const url = 'http://127.0.0.1:8080/function/homepage-nonflow';
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const payload = JSON.stringify({
        user_id: userId,
        origin: origin
    });

    const response = http.post(url, payload, params);

    check(response, {
        'status is 200': (r) => r.status === 200,
        'response body has content': (r) => r.body.length > 0,
    });

    sleep(1);
}