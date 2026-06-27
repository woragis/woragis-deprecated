import pytest
from fastapi.testclient import TestClient
from unittest.mock import patch
from app.main import app

client = TestClient(app)


@pytest.fixture
def run_payload():
    return {"max_applications": 2, "platforms": ["greenhouse"]}


@patch("app.llm.fit_check.should_apply", return_value={"apply": True, "score": 0.9, "reason": "Good fit"})
@patch("app.llm.text_gen.generate_application_text", return_value={"cover_letter": "Test", "answers": ["A1"]})
@patch("app.platforms.greenhouse.GreenhousePlatform.collect_jobs", return_value=[{"id": "job1", "url": "http://test/job1"}, {"id": "job2", "url": "http://test/job2"}])
@patch("app.platforms.greenhouse.GreenhousePlatform.apply_to_job", return_value=True)
def test_happy_path_run(mock_apply, mock_collect, mock_text, mock_fit, run_payload):
    # Start run
    resp = client.post("/runs", json=run_payload)
    assert resp.status_code == 200
    run_id = resp.json()["run_id"]
    # Simulate status polling (in real system, would be async)
    status = client.get(f"/runs/{run_id}").json()
    assert status["state"] in ("PENDING", "RUNNING", "COMPLETED")
    # (In a real async system, would wait for completion and check logs/events)
