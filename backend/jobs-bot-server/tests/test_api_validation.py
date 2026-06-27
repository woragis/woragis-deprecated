import pytest
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)


def test_run_request_validation():
    # Invalid platforms
    resp = client.post(
        "/runs", json={"max_applications": 5, "platforms": ["invalid"]})
    assert resp.status_code == 400
    assert resp.json()["detail"]["code"].startswith("VALID001")

    # Invalid max_applications
    resp = client.post(
        "/runs", json={"max_applications": 0, "platforms": ["greenhouse"]})
    assert resp.status_code == 400
    assert resp.json()["detail"]["code"].startswith("VALID002")

    # Valid request
    resp = client.post(
        "/runs", json={"max_applications": 2, "platforms": ["greenhouse"]})
    assert resp.status_code == 200
