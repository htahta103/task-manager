.PHONY: build test lint typecheck dev deploy-staging deploy-prod

build:
	cd backend && go build ./...

test:
	cd backend && go test ./...

lint:
	@test -z "$$(gofmt -l backend | tr -d '\n')" || (echo "gofmt needed. Run: gofmt -w backend"; exit 1)

typecheck:
	cd backend && go vet ./...

dev:
	docker compose --env-file .env up --build

deploy-staging:
	@echo "Set up Fly + Pages first, then deploy:"
	@echo "  fly deploy -a $$FLY_APP_STAGING"

deploy-prod:
	@echo "Set up Fly + Pages first, then deploy:"
	@echo "  fly deploy -a $$FLY_APP_PROD"
