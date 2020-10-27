FROM --platform=${BUILDPLATFORM} golang:1.15-buster as builder

ARG GOPROXY
WORKDIR /go/src
COPY ./ ./

ENV GOBIN=/go/bin
ARG TARGETARCH
RUN GOARCH=${TARGETARCH} make build

FROM ghcr.io/querycap/distroless/static-debian10:latest

ARG TARGETARCH
COPY --from=builder /go/bin/qservice-operator-linux-${TARGETARCH} /go/bin/qservice-operator

ENTRYPOINT ["/go/bin/qservice-operator"]