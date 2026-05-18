.PHONY: test-unit test-integration docker-up

test-unit:
	go test ./... -count=1 -coverprofile=coverage-unit.out -covermode=atomic -coverpkg=./...
	@go tool cover -func=coverage-unit.out | tail -1

test-integration:
	go test -tags=integration -count=1 -timeout=5m -coverprofile=coverage-integration.out -covermode=atomic -coverpkg=./... .
	@go tool cover -func=coverage-integration.out | tail -1

docker-up:
	docker compose up --build
