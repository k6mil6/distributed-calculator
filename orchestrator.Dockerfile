FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY backend/internal/storage/migrations /usr/local/src/distributed-calculator/backend/internal/storage/migrations

COPY ./ ./
RUN go build -o ./bin/app backend/cmd/orchestrator/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app
COPY --from=builder /usr/local/src/distributed-calculator/backend/internal/storage/migrations /migrations
COPY config.hcl /config.hcl

CMD ["/app"]