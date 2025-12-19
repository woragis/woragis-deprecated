"""
Performance tests for AI service
Run with: pytest app/tests/performance_test.py -v
"""

import pytest
import time
import asyncio
from concurrent.futures import ThreadPoolExecutor
from fastapi.testclient import TestClient
from app.main import app


@pytest.fixture
def client():
    return TestClient(app)


def test_health_endpoint_performance(client):
    """Test health endpoint performance"""
    iterations = 100
    latencies = []

    for _ in range(iterations):
        start = time.time()
        response = client.get("/healthz")
        latency = (time.time() - start) * 1000  # Convert to ms
        latencies.append(latency)
        assert response.status_code == 200

    avg_latency = sum(latencies) / len(latencies)
    sorted_latencies = sorted(latencies)
    p95_index = int(len(latencies) * 0.95)
    p95_latency = sorted_latencies[p95_index]

    print(f"\nHealth Endpoint Performance:")
    print(f"  Iterations: {iterations}")
    print(f"  Average Latency: {avg_latency:.2f}ms")
    print(f"  P95 Latency: {p95_latency:.2f}ms")

    assert avg_latency < 50, "Average latency should be under 50ms"
    assert p95_latency < 100, "P95 latency should be under 100ms"


def test_concurrent_requests(client):
    """Test API under concurrent load"""
    concurrent_requests = 20
    requests_per_thread = 5
    total_requests = concurrent_requests * requests_per_thread

    def make_request():
        response = client.get("/healthz")
        return response.status_code == 200

    start = time.time()
    with ThreadPoolExecutor(max_workers=concurrent_requests) as executor:
        futures = [
            executor.submit(make_request)
            for _ in range(total_requests)
        ]
        results = [f.result() for f in futures]

    duration = time.time() - start
    success_count = sum(results)

    throughput = total_requests / duration

    print(f"\nConcurrent Requests Test:")
    print(f"  Concurrent Threads: {concurrent_requests}")
    print(f"  Requests per Thread: {requests_per_thread}")
    print(f"  Total Requests: {total_requests}")
    print(f"  Successful: {success_count}")
    print(f"  Duration: {duration:.2f}s")
    print(f"  Throughput: {throughput:.2f} req/s")

    assert success_count >= total_requests * 0.9, "At least 90% should succeed"
    assert duration < 10, "Should complete in under 10 seconds"


@pytest.mark.skip(reason="Requires AI API keys - run manually")
def test_ai_endpoint_load(client):
    """Test AI endpoint under load (requires API keys)"""
    concurrent_requests = 10
    requests_per_thread = 2

    def make_ai_request():
        payload = {
            "messages": [{"role": "user", "content": "Hello"}],
            "model": "test-model"
        }
        response = client.post("/api/v1/chat", json=payload)
        return response.status_code in [200, 201]

    start = time.time()
    with ThreadPoolExecutor(max_workers=concurrent_requests) as executor:
        futures = [
            executor.submit(make_ai_request)
            for _ in range(concurrent_requests * requests_per_thread)
        ]
        results = [f.result() for f in futures]

    duration = time.time() - start
    success_count = sum(results)

    print(f"\nAI Endpoint Load Test:")
    print(f"  Total Requests: {concurrent_requests * requests_per_thread}")
    print(f"  Successful: {success_count}")
    print(f"  Duration: {duration:.2f}s")

    assert success_count > 0, "Some requests should succeed"


if __name__ == "__main__":
    pytest.main([__file__, "-v", "-s"])
