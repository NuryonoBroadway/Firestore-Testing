proto:
	protoc --proto_path=proto --go_out=collectionx --go_opt=paths=source_relative \
	--go-grpc_out=collectionx --go-grpc_opt=paths=source_relative \
	proto/*.proto


.PHONY: proto
