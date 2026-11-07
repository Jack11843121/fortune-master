# 多阶段构建：前端 (Node) → 后端 (Go) → 精简运行镜像
FROM node:20-alpine AS fe
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

FROM golang:1.27-alpine AS be
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=fe /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /fortune ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
WORKDIR /app
COPY --from=be /fortune .
COPY config.yaml ./config.yaml
EXPOSE 8080
CMD ["./fortune", "-config", "config.yaml"]