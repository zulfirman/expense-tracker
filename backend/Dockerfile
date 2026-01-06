########################################
# Builder
########################################
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 1) Cache dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# 2) Copy source
COPY . .

# 3) Build static binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -trimpath -o webserver ./main.go


########################################
# Final
########################################
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/webserver .

# Run as nonroot user provided by distroless
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/webserver"]

