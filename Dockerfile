FROM golang:1.21-alpine as builder

WORKDIR /app
COPY main.go .

# Build the binary
RUN go build -o server main.go

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
