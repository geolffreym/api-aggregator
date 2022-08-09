import http from "k6/http";
import { check, group, sleep } from "k6";
import { Rate, Counter } from "k6/metrics";

const endpoint = "http://localhost:3333"
// A custom metric to track failure rates
const failureRate = new Rate("check_failure_rate");
const contentMatchCounter = new Counter('content_match_error');
// sample endpoints
const blockByNumber = `${endpoint}/v1/block/by/number/0x5BAD55/true`,
    transactionByNumber = `${endpoint}/v1/tx/by/number/0x5BAD55/0x0`,
    blockByHash = `${endpoint}/v1/block/by/hash/0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35/true`,
    transactionByHash = `${endpoint}/v1/tx/by/hash/0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35/0x0`

// Options
export let options = {
    stages: [
        // Linearly ramp up from 1 to 100 VUs during first minute
        { target: 100, duration: "1m" },
        // Hold at 50 VUs for the next 3 minutes and 30 seconds
        { target: 50, duration: "3m30s" },
        // Linearly ramp down from 50 to 30 VUs over the last 1 minute
        { target: 30, duration: "1m" },
        // Linearly ramp down from 50 to 0 VUs over the last 30 seconds
        { target: 0, duration: "30s" }
    ],
    thresholds: {
        'content_match_error': ['count < 5'], // 5 or fewer total errors are tolerated
        'group_duration{group:::Blocks}': ['avg < 200'],
        'group_duration{group:::Transactions}': ['avg < 200'],
        // We want the 95th percentile of all HTTP request durations to be less than 500ms
        "http_req_duration": ["p(95)<500"],
        // Thresholds based on the custom metric we defined and use to track application failures
        "check_failure_rate": [
            // Global failure rate should be less than 1%
            "rate<0.01",
            // Abort the test early if it climbs over 5%
            { threshold: "rate<=0.05", abortOnFail: true },
        ],
    },
};

// Main function
export default function () {

    // TODO: more test could be added for expected body results
    // Expected valid result body for valid request
    const expectedHash = "0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35"
    const singleRequestToTransactionByNumber = http.get(transactionByNumber)
    const expectedResponse = singleRequestToTransactionByNumber.json("blockHash") == expectedHash
    contentMatchCounter.add(!expectedResponse)

    group("Blocks", function () {
        // Execute multiple requests in parallel like a browser, to fetch some static resources
        let resps = http.batch([
            ["GET", blockByHash],
            ["GET", blockByNumber],
        ]);

        // Combine check() call with failure tracking
        failureRate.add(!check(resps, {
            // Expected status for block endpoints
            "status is 200": (r) => r[0].status === 200 && r[1].status === 200,
        }));
    });


    group("Transactions", function () {
        // Execute multiple requests in parallel like a browser, to fetch some static resources
        let resps = http.batch([
            ["GET", transactionByNumber],
            ["GET", transactionByHash],
        ]);

        // Combine check() call with failure tracking
        failureRate.add(!check(resps, {
            // Expected status for transactions endpoints
            "status is 200": (r) => r[0].status === 200 && r[1].status === 200,
        }));
    });


    sleep(Math.random() * 3 + 2); // Random sleep between 2s and 5s
}