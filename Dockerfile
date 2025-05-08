FROM golang:1.19-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/zkkyc \
    ./cmd/zkKYC

FROM alpine:latest

RUN apk add --no-cache ca-certificates && \
    adduser -D appuser

USER appuser
WORKDIR /home/appuser

COPY --from=builder --chown=appuser:appuser /app/zkkyc .

EXPOSE 8080

ENTRYPOINT ["./zkkyc"]