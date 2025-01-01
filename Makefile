build:
	go build cmd/main.go

local-deploy:
	docker-compose -f deployments/docker-compose.yml up -d --force-recreate --remove-orphans

run:
	go run cmd/main.go

unit-test:
	go test -v ./internal/...

integration-test:
	go test -v ./test/...

test: unit-test integration-test
