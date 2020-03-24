go:
	protoc -I proto/ proto/uiot.proto --go_out=plugins=grpc:proto/.
