.PHONY: test

test:
	cd apps/api && GOTOOLCHAIN=local go test ./...
