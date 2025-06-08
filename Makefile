build:
	go build cmd/main.go

local-deploy:
	docker-compose -f deployments/docker-compose.yml up -d --force-recreate --remove-orphans

run:
	go run cmd/main.go

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/recipe-manager/infrastructure/grpc/proto/calculator.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/recipe-manager/infrastructure/grpc/proto/ingredients_balancer.proto
	mkdir -p internal/recipe-manager/infrastructure/grpc/proto/generated
	mv internal/recipe-manager/infrastructure/grpc/proto/*.pb.go internal/recipe-manager/infrastructure/grpc/proto/generated/

unit-test:
	go test -v ./internal/...

integration-test:
	go test -v ./test/...

test: unit-test integration-test
