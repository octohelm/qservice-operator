# syntax = docker/dockerfile:experimental

FROM golang:1.15-buster as builder

ARG GOPROXY=https://proxy.golang.org,direct
ENV GOBIN=/go/bin

WORKDIR /go/src
COPY ./ ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    make build

FROM debian:buster-slim

COPY --from=builder /go/bin/qservice-operator /go/bin/qservice-operator

ENTRYPOINT ["/go/bin/qservice-operator"]
