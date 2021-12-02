# Builder

FROM golang:1.16-alpine AS builder
ARG GITHUB_PATH=github.com/ozonmp/lgc-location-api

WORKDIR /home/${GITHUB_PATH}

RUN apk add --update make git protoc protobuf protobuf-dev curl
COPY ../.. .
RUN make build-retranslator

# Retranslator

FROM alpine:latest as server
ARG GITHUB_PATH=github.com/ozonmp/lgc-location-api
LABEL org.opencontainers.image.source=https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/retranslator .
COPY --from=builder /home/${GITHUB_PATH}/config-rt.yaml .

RUN chown root:root retranslator

CMD ["./retranslator"]
