FROM golang:alpine as builder
MAINTAINER Jack Murdock <jack_murdock@comcast.com>

WORKDIR /go/src/github.com/comcast/codex-fenrir

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

COPY . .
RUN GO111MODULE=on go mod vendor
RUN GO111MODULE=on go build -o fenrir_linux_amd64

FROM alpine

RUN apk --no-cache add ca-certificates
RUN mkdir -p /etc/fenrir
VOLUME /etc/fenrir

EXPOSE 8080

COPY --from=builder /go/src/github.com/comcast/codex-fenrir/fenrir_linux_amd64 /
COPY fenrir.yaml /
ENTRYPOINT ["/fenrir_linux_amd64"]
