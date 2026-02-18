# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/crypto-api ./cmd/api/main.go


FROM alpine:3.20
RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN adduser -D -H -u 10001 appuser
USER appuser
WORKDIR /app

COPY --from=build /out/crypto-api /app/crypto-api
COPY --from=build /src/frontend /app/frontend

EXPOSE 8080
ENTRYPOINT ["/app/crypto-api"]
