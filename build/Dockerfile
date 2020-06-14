FROM golang:1.14-alpine as builder

ARG GOPROXY=https://proxy.golang.org,direct

ENV OPERATOR_PKG=github.com/octohelm/qservice-operator

COPY ./ /go/src/${OPERATOR_PKG}
WORKDIR /go/src/${OPERATOR_PKG}/cmd/manager
RUN go build -ldflags "-X ${OPERATOR_PKG}/version.Version=$(cat .version)" -o qservice-operator && cp qservice-operator /go/bin/

FROM alpine

ENV OPERATOR=/usr/local/bin/qservice-operator \
    USER_UID=1001 \
    USER_NAME=qservice-operator

COPY --from=builder /go/bin/qservice-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
