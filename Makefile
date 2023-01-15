generate_pb:
	ECHO Y | powershell rm pb\*
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

test:
	go test -v ./...
