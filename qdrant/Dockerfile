# syntax=docker/dockerfile:1

FROM golang:1.26.1-alpine AS builder

ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/swarm-deploy/init-jobs/qdrant

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X 'main.Version=${APP_VERSION}' -X 'main.BuildTime=${BUILD_TIME}'" -o /go/bin/init-job /go/src/github.com/swarm-deploy/init-jobs/qdrant/main.go

######################################################

FROM scratch

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/bin/init-job /go/bin/init-job

LABEL org.opencontainers.image.title="Init Job for create Qdrant collections"
LABEL org.opencontainers.image.description="Init Job for create Qdrant collections"
LABEL org.opencontainers.image.url="https://github.com/swarm-deploy/init-jobs"
LABEL org.opencontainers.image.source="https://github.com/swarm-deploy/init-jobs"
LABEL org.opencontainers.image.vendor="swarm-deploy"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL org.opencontainers.image.created="$BUILD_TIME"
LABEL org.opencontainers.image.licenses="Apache 2.0"
LABEL org.swarm-deploy.service.type="InitJob"

ENTRYPOINT ["/go/bin/init-job"]
