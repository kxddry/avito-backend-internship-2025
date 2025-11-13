FROM golang:1.25.4-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY "cmd/" "./cmd/"

RUN go build -o app ./cmd/app

FROM alpine:latest

COPY --from=builder /build/app /app/app

ENTRYPOINT ["/app/app"]