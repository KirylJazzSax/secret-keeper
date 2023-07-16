generate:
	rm internal/common/gen/$(name)/* || true
	rm docs/swagger/*.swagger.json || true
	protoc --proto_path=api/protobuf --go_out=internal/common/gen/$(name) --go_opt=paths=source_relative \
	--go-grpc_out=internal/common/gen/$(name) --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/common/gen/$(name) --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=secret-keeper \
	 api/protobuf/$(name).proto

test:
	go test -v ./...

up:
	docker compose up $(services) --force-recreate
