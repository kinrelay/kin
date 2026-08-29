.PHONY: test test-ci

test: test-ci
	cd apps/api && GOTOOLCHAIN=local go test ./...

test-ci:
	./scripts/ci/test-detect-backend-changes.sh
