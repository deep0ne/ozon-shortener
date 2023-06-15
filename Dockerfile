# syntax=docker/dockerfile:1

FROM golang:1.19-alpine3.17 AS builder

# Set destination for COPY
WORKDIR /app

COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server/main.go

FROM alpine:3.17
RUN apk update && apk add netcat-openbsd
WORKDIR /app
COPY --from=builder /app/app .
COPY wait-for.sh .
RUN chmod +x wait-for.sh

EXPOSE 8082

# Run
CMD ["./app"]