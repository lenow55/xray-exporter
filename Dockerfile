# syntax=docker/dockerfile:1

FROM golang:1.24-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY *.go ./
COPY internal/ ./internal/

ARG TARGETOS=linux
ARG TARGETARCH
ARG TARGETVARIANT

ARG VERSION=dev
ARG COMMIT=none

RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	set -eux; \
	export CGO_ENABLED=0; \
	export GOOS="${TARGETOS}" GOARCH="${TARGETARCH}"; \
	if [ "${TARGETARCH}" = "arm" ] && [ "${TARGETVARIANT:-}" = "v7" ]; then export GOARM=7; fi; \
	BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"; \
	go build -trimpath \
		-ldflags="-s -w -X main.buildVersion=${VERSION} -X main.buildCommit=${COMMIT} -X main.buildDate=${BUILD_DATE}" \
		-o /out/xray-exporter .

FROM alpine:3.21
RUN apk --no-cache add ca-certificates

COPY --from=builder /out/xray-exporter /usr/bin/xray-exporter

EXPOSE 9550
ENTRYPOINT ["/usr/bin/xray-exporter"]
