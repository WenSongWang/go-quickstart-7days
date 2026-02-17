# Day6 示例：多阶段构建 Go 镜像
# 在项目根目录执行: docker build -f day6/Dockerfile -t go-7days-app .

FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY go.sum* ./
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./day7/cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /server .
EXPOSE 8080
CMD ["./server"]
