FROM golang:1.10-alpine as builder

RUN apk update && apk add git protobuf protobuf-dev &&\
    go get -v github.com/golang/dep/cmd/dep &&\
    go get -v google.golang.org/grpc &&\
    go get -v github.com/golang/protobuf/protoc-gen-go
COPY ./Gopkg.* /go/src/github.com/nlepage/grpc-hellow/go/hellow/
WORKDIR /go/src/github.com/nlepage/grpc-hellow/go/hellow
RUN dep ensure -v -vendor-only
COPY . /go/src/github.com/nlepage/grpc-hellow/go/hellow
COPY --from=nlepage/grpc-hellow:service /hellow.proto /go/src/github.com/nlepage/grpc-hellow/go/hellow/service/
RUN protoc -I. \
           --go_out=plugins=grpc:. \
           service/hellow.proto &&\
    CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"'

FROM scratch

COPY --from=builder /go/bin/hellow /

ENTRYPOINT ["/hellow"]
