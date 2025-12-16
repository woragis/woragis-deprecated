.PHONY: test test-unit test-integration test-cov install clean

# Install dependencies
install:
	npm install

# Run all tests
test:
	npm run test:jest

# Run unit tests only
test-unit:
	npm run test:jest -- tests/__tests__/

# Run tests with coverage
test-cov:
	npm run test:jest:coverage

# Clean test artifacts
clean:
	rm -rf coverage
	rm -rf node_modules/.cache
	rm -rf .jest-cache
