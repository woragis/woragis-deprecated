.PHONY: test test-unit test-integration test-cov install clean

# Install dependencies
install:
	pip install -r requirements.txt

# Run all tests
test:
	pytest

# Run unit tests only
test-unit:
	pytest tests/unit/ -v

# Run integration tests only
test-integration:
	pytest tests/integration/ -v -m integration

# Run tests with coverage
test-cov:
	pytest --cov=app --cov-report=html --cov-report=term

# Clean test artifacts
clean:
	rm -rf .pytest_cache
	rm -rf htmlcov
	rm -rf .coverage
	rm -rf *.pyc
	rm -rf __pycache__
	find . -type d -name __pycache__ -exec rm -r {} +
	find . -type f -name "*.pyc" -delete
