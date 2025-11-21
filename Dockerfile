# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/web ./cmd/web

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/web /app/web
COPY ui/ /app/ui

EXPOSE 4000

CMD ["/app/web"]
