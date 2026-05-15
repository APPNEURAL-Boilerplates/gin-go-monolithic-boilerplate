# syntax=docker/dockerfile:1

FROM golang:1.26.3-alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/api ./cmd/api

FROM alpine:3.22
RUN adduser -D -H -u 10001 appuser
WORKDIR /app

COPY --from=build /bin/api /app/api

ENV APP_ENV=production \
    GIN_MODE=release \
    PORT=8080 \
    LOG_LEVEL=info

USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/api"]
