# 这是 Go 语言的 Dockerfile 示例
FROM golang:1.25.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
