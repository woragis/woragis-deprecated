"""
Performance tests for docs service
Run with: pytest tests/performance_test.py -v
"""

import pytest
import time
from concurrent.futures import ThreadPoolExecutor
from fastapi.testclient import TestClient


@pytest.fixture
def client():
    # Import here to avoid issues if service isn't set up
    try:
        from app.main import app
        return TestClient(app)
    except ImportError:
        pytest.skip("Docs service not available")


def test_health_endpoint_performance(client):
    """Test health endpoint performance"""
    iterations = 100
    latencies = []

    for _ in range(iterations):
        start = time.time()
        response = client.get("/healthz")
        latency = (time.time() - start) * 1000
        latencies.append(latency)
        assert response.status_code == 200

    avg_latency = sum(latencies) / len(latencies)
    sorted_latencies = sorted(latencies)
    p95_index = int(len(latencies) * 0.95)
    p95_latency = sorted_latencies[p95_index]

    print(f"\nHealth Endpoint Performance:")
    print(f"  Average Latency: {avg_latency:.2f}ms")
    print(f"  P95 Latency: {p95_latency:.2f}ms")

    assert avg_latency < 50
    assert p95_latency < 100


def test_concurrent_requests(client):
    """Test concurrent request handling"""
    concurrent_requests = 20
    requests_per_thread = 5

    def make_request():
        response = client.get("/healthz")
        return response.status_code == 200

    start = time.time()
    with ThreadPoolExecutor(max_workers=concurrent_requests) as executor:
        futures = [
            executor.submit(make_request)
            for _ in range(concurrent_requests * requests_per_thread)
        ]
        results = [f.result() for f in futures]

    duration = time.time() - start
    success_count = sum(results)

    print(f"\nConcurrent Requests: {success_count}/{len(results)}")
    print(f"Duration: {duration:.2f}s")

    assert success_count >= len(results) * 0.9
