# Development image: Go toolchain plus air for live reload.
# The project is bind-mounted at /app, so code changes rebuild without a new image.
FROM golang:1.25

ENV GOMODCACHE=/go/pkg/mod \
    GOCACHE=/go/cache \
    GOFLAGS=-buildvcs=false \
    HOME=/tmp

# Match the host user so the telemetry resource detector can resolve the uid
ARG UID=1000
ARG GID=1000
RUN groupadd -g ${GID} dev 2>/dev/null || true \
    && useradd -u ${UID} -g ${GID} -M -d /tmp dev 2>/dev/null || true

RUN GOBIN=/usr/local/bin go install github.com/air-verse/air@v1.62.0

WORKDIR /app

# Warm the module cache so the first `air` run does not download everything.
COPY go.mod go.sum ./
RUN go mod download && mkdir -p /go/cache && chmod -R a+rwX /go/pkg/mod /go/cache

EXPOSE 8000

CMD ["air", "-c", ".air.toml"]
