FROM golang:1.25-alpine AS builder
WORKDIR /app

ARG VERSION

COPY . .

RUN go mod tidy && go build -ldflags "-s -w -X main.Version=${VERSION}" -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]