ARG BASE_IMAGE=alpine

FROM golang:1.19.4-alpine3.16 AS builder

WORKDIR /usr/local/src/bee

RUN apk add --no-cache --update alpine-sdk ca-certificates openssl

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT=""

ENV GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT}

ARG GOPROXY

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make release-binary

RUN make build
FROM alpine:3.17.0 AS stager

RUN mkdir -p /var/bee
RUN mkdir -p /etc/bee
COPY configs/config.docker.yaml /etc/bee/
COPY configs/boot.yaml /etc/bee/

# For Dependabot to detect base image versions
FROM alpine:3.17.0 AS alpine
FROM gcr.io/distroless/static:latest AS distroless

FROM $BASE_IMAGE

# Dex connectors, such as GitHub and Google logins require root certificates.
# Proper installations should manage those certificates, but it's a bad user
# experience when this doesn't work out of the box.
#
# See https://go.dev/src/crypto/x509/root_linux.go for Go root CA bundle locations.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=stager --chown=1001:1001 /var/bee /var/bee
COPY --from=stager --chown=1001:1001 /etc/bee /etc/bee

# Copy module files for CVE scanning / dependency analysis.
COPY --from=builder /usr/local/src/bee/go.mod /usr/local/src/bee/go.sum /usr/local/src/bee/
#COPY --from=builder /usr/local/src/bee/api/v2/go.mod /usr/local/src/bee/api/v2/go.sum /usr/local/src/bee/api/v2/

COPY --from=builder /go/bin/bee /usr/local/bin/bee
#COPY --from=builder /go/bin/docker-entrypoint /usr/local/bin/docker-entrypoint
#COPY --from=builder /usr/local/src/bee/web /srv/bee/web

#COPY --from=gomplate /usr/local/bin/gomplate /usr/local/bin/gomplate

USER 1001:1001

#ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
CMD ["bee", "serve", "/etc/bee/config.docker.yaml"]
