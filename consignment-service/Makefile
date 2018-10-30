build:
	protoc -I. --go_out=plugins=micro:/Users/fjordan/go/src/github.com/fjordansilva/micros/consignment-service \
	  proto/consignment/consignment.proto
	docker build -t consignment-service .

run:
	docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns consignment-service