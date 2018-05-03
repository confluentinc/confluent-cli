PROTO = shared/connect

compile-proto:
	protoc -I $(PROTO) $(PROTO)/*.proto --gogo_out=plugins=grpc:$(PROTO)

install-plugins:
	go install ./plugin/...

clean:
	rm $(PROTO)/*.pb.go
