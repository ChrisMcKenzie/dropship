echo "Generating Protobuf Code..."
protoc --go_out=plugins=grpc:. dropship/rpc.proto
