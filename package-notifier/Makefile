.PHONY: proto swagger

# Generate protobuf and gRPC code
proto:
	protoc -I . \
		-I third_party \
		-I $(GOPATH)/src \
		-I $(GOPATH)/pkg/mod \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=. --openapiv2_opt=logtostderr=true \
		proto/notification.proto

# Generate Swagger documentation
swagger:
	swag init -g main.go -o docs

# Install required tools
install-tools:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Run the server
run:
	go run main.go

# Build the server
build:
	go build -o bin/notifier main.go 