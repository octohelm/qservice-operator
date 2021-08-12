FROM gcr.io/distroless/static-debian10:latest

ARG TARGETARCH
COPY ./bin/qservice-operator_linux_${TARGETARCH}/qservice-operator /go/bin/qservice-operator

ENTRYPOINT ["/go/bin/qservice-operator"]