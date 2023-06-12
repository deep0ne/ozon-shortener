# syntax=docker/dockerfile:1

FROM golang:1.18 AS builder

# Set destination for COPY
WORKDIR /app

COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server/main.go

FROM scratch
COPY --from=builder /app/app /app

EXPOSE 8080

# Run
CMD ["./app"]