.PHONY: vet test-unit test-integration merge-coverage docker-up docker-up-dev

vet:
	go vet ./...

# Full-module statement coverage (matches CI unit job).
test-unit: vet
	go test ./... -count=1 -coverprofile=coverage-unit.out -covermode=atomic -coverpkg=./...
	@go tool cover -func=coverage-unit.out | tail -1

# Requires MySQL (see AGENTS.md) and config.toml at repo root.
test-integration: vet
	go test -tags=integration -count=1 -timeout=5m -coverprofile=coverage-integration.out -covermode=atomic -coverpkg=./... .

merge-coverage: vet
	go test ./... -count=1 -coverprofile=coverage-unit.out -covermode=atomic -coverpkg=./...
	go test -tags=integration -count=1 -timeout=5m -coverprofile=coverage-integration.out -covermode=atomic -coverpkg=./... .
	python3 scripts/merge_cover_profiles.py coverage-merged.out coverage-unit.out coverage-integration.out
	@go tool cover -func=coverage-merged.out | tail -1

docker-up:
	docker compose up -d

docker-up-dev:
	docker compose --profile dev up -d
