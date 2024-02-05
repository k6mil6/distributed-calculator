FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/app backend/cmd/agent/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /

CMD ["/app"]