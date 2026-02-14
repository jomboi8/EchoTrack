# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o echostrike ./cmd/echostrike

# Final Stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/echostrike .

# Expose nothing by default (it's a client), but good practice to document
ENTRYPOINT ["./echostrike"]
