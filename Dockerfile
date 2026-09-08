# 构建阶段：前端
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --no-audit --no-fund
COPY frontend/ ./
RUN npm run build

# 构建阶段：后端
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/frontend/dist ./static
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o doudian .

# 运行阶段
FROM alpine:3.19
LABEL maintainer="doudian-admin"
LABEL description="DouDian - 抖店一键代发供应商管理系统"

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /opt/doudian

COPY --from=backend-builder /app/doudian .
COPY --from=backend-builder /app/static ./static

VOLUME ["/opt/doudian/data"]

ENV HOST=0.0.0.0
ENV PORT=2095
ENV DB_PATH=/opt/doudian/data/doudian.db
ENV JWT_SECRET=change-me-in-production
ENV SECRET_PATH=
ENV LOG_LEVEL=info

EXPOSE 2095

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:${PORT}/ || exit 1

CMD ["./doudian"]
