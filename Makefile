.PHONY: test-unit test-integration docker-up docker-dev-up vet fmt-check

test-unit:
	go test ./... -count=1 -coverprofile=coverage-unit.out -covermode=atomic -coverpkg=./...
	@go tool cover -func=coverage-unit.out | tail -1

test-integration:
	go test -tags=integration -count=1 -timeout=5m -coverprofile=coverage-integration.out -covermode=atomic -coverpkg=./... .
	@go tool cover -func=coverage-integration.out | tail -1

docker-up:
	docker compose up --build

docker-dev-up:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

vet:
	go vet ./...

fmt-check:
	@if gofmt -l $$(go list -f '{{.Dir}}' ./...) | grep -q .; then gofmt -l $$(go list -f '{{.Dir}}' ./...); exit 1; fi
