FROM golang:1.16 AS builder

WORKDIR /go/src/k8s-secret-injector/

COPY go.mod /go/src/k8s-secret-injector/
COPY go.sum go.sum
RUN go mod download

COPY . /go/src/k8s-secret-injector/
RUN go get -v -t -d ./... \
    && go build -o k8s-secret-injector

FROM alpine:latest
COPY --from=builder /go/src/k8s-secret-injector/k8s-secret-injector /usr/local/bin/
RUN apk add --no-cache libc6-compat
