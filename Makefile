proto:
	protoc --proto_path=proto --go_out=collectionx/collectionx_service --go_opt=paths=source_relative \
	--go-grpc_out=collectionx/collectionx_service --go-grpc_opt=paths=source_relative \
	proto/*.proto


.PHONY: proto
