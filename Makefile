.PHONY: all build frontend backend clean run install uninstall

APP_NAME := doudian
VERSION := 1.0.0
GO := go
NPM := npm

FRONTEND_DIR := frontend

all: build

## build: 构建前端和后端（生产版本）
build: frontend backend

## frontend: 构建前端
frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && $(NPM) install && $(NPM) run build
	@rm -rf static
	@cp -r $(FRONTEND_DIR)/dist static
	@echo "Frontend built successfully"

## backend: 构建后端
backend:
	@echo "Building backend..."
	CGO_ENABLED=0 $(GO) build -ldflags "-s -w" -o $(APP_NAME) .
	@echo "Backend built successfully"

## build-linux: 交叉编译 Linux amd64 版本
build-linux:
	@echo "Building Linux amd64 version..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w" -o $(APP_NAME)-linux-amd64 .
	@echo "Linux amd64 binary: $(APP_NAME)-linux-amd64"

## build-arm64: 交叉编译 Linux arm64 版本
build-arm64:
	@echo "Building Linux arm64 version..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -ldflags "-s -w" -o $(APP_NAME)-linux-arm64 .
	@echo "Linux arm64 binary: $(APP_NAME)-linux-arm64"

## run: 运行开发服务器
run:
	@echo "Starting server..."
	PORT=2095 DB_PATH=./instance/doudian.db JWT_SECRET=change-me-in-production ./$(APP_NAME)

## dev-frontend: 启动前端开发服务器
dev-frontend:
	cd $(FRONTEND_DIR) && $(NPM) run dev

## clean: 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -f $(APP_NAME) $(APP_NAME)-linux-*
	@rm -rf static
	@rm -rf $(FRONTEND_DIR)/dist
	@rm -rf $(FRONTEND_DIR)/node_modules
	@echo "Cleaned"

## install: 安装到系统（Linux）
install: build
	@echo "Installing $(APP_NAME)..."
	@mkdir -p /opt/$(APP_NAME)
	@mkdir -p /var/log/$(APP_NAME)
	@cp $(APP_NAME) /opt/$(APP_NAME)/
	@cp -r static /opt/$(APP_NAME)/
	@cp install.sh /opt/$(APP_NAME)/
	@echo "Installed to /opt/$(APP_NAME)"

## uninstall: 卸载（Linux）
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	@systemctl stop $(APP_NAME) 2>/dev/null || true
	@systemctl disable $(APP_NAME) 2>/dev/null || true
	@rm -f /etc/systemd/system/$(APP_NAME).service
	@rm -rf /opt/$(APP_NAME)
	@rm -rf /var/log/$(APP_NAME)
	@echo "Uninstalled"

## help: 显示帮助信息
help:
	@echo "DouDian Supplier Management System Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^## / {desc=$$0; next} /^[a-zA-Z_-]+:/ {gsub(/:.*/, "", $$1); printf "  %-20s %s\n", $$1, substr(desc, 4)}' Makefile
