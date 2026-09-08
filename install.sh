#!/bin/bash

# DouDian - 抖店一键代发供应商管理系统
# 一键安装脚本 (3x-ui 风格部署)
# 支持: Debian / Ubuntu / CentOS / Rocky Linux

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PLAIN='\033[0m'

APP_NAME="doudian"
APP_DIR="/opt/doudian"
LOG_DIR="/var/log/doudian"
DATA_DIR="/opt/doudian/data"
CONFIG_FILE="/etc/default/doudian"
SERVICE_FILE="/etc/systemd/system/doudian.service"
DEFAULT_PORT=2095

detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        OS_VER=$VERSION_ID
    elif [ -f /etc/redhat-release ]; then
        OS="centos"
        OS_VER=$(rpm -q --qf "%{VERSION}" centos-release 2>/dev/null || echo "7")
    else
        OS=$(uname -s)
        OS_VER=$(uname -r)
    fi
    echo "检测到系统: $OS $OS_VER"
}

check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}错误：请使用 root 用户运行此脚本${PLAIN}"
        exit 1
    fi
}

check_installed() {
    if [ -d "$APP_DIR" ] && [ -f "$APP_DIR/doudian" ]; then
        echo -e "${YELLOW}检测到已安装的 $APP_NAME，将进行更新...${PLAIN}"
        return 1
    fi
    return 0
}

install_deps() {
    echo -e "${BLUE}安装系统依赖...${PLAIN}"
    case "$OS" in
        debian|ubuntu)
            apt-get update -qq
            apt-get install -y -qq curl wget tar sqlite3 openssl ca-certificates
            ;;
        centos|rocky|rhel|fedora)
            yum install -y -q curl wget tar sqlite openssl ca-certificates
            ;;
        *)
            echo -e "${YELLOW}未知系统类型，尝试继续...${PLAIN}"
            ;;
    esac
    echo -e "${GREEN}依赖安装完成${PLAIN}"
}

get_arch() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo -e "${RED}不支持的架构: $ARCH${PLAIN}"; exit 1 ;;
    esac
    echo "架构: $ARCH"
}

generate_secret() {
    openssl rand -hex 32
}

get_user_input() {
    echo ""
    echo -e "${BLUE}=== 配置参数 ===${PLAIN}"
    
    read -p "请输入服务端口 [默认 $DEFAULT_PORT]: " INPUT_PORT
    PORT=${INPUT_PORT:-$DEFAULT_PORT}
    
    if ! [[ "$PORT" =~ ^[0-9]+$ ]] || [ "$PORT" -lt 1 ] || [ "$PORT" -gt 65535 ]; then
        echo -e "${RED}无效端口号，使用默认端口 $DEFAULT_PORT${PLAIN}"
        PORT=$DEFAULT_PORT
    fi
    
    if ss -tlnp | grep -q ":$PORT "; then
        echo -e "${YELLOW}警告：端口 $PORT 已被占用${PLAIN}"
    fi
    
    read -p "请输入安全路径前缀（留空则不启用）: " SECRET_PATH
    
    read -p "请输入管理员密码 [默认 admin123]: " ADMIN_PASSWORD
    ADMIN_PASSWORD=${ADMIN_PASSWORD:-admin123}
    
    echo ""
    echo -e "${BLUE}=== 配置确认 ===${PLAIN}"
    echo "端口: $PORT"
    echo "安全路径: ${SECRET_PATH:-无}"
    echo "管理员密码: $ADMIN_PASSWORD"
    echo ""
    read -p "确认安装？[Y/n]: " CONFIRM
    CONFIRM=${CONFIRM:-Y}
    if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
        echo "安装已取消"
        exit 0
    fi
}

create_dirs() {
    mkdir -p "$APP_DIR"
    mkdir -p "$DATA_DIR"
    mkdir -p "$LOG_DIR"
}

write_config() {
    JWT_SECRET=$(generate_secret)
    
    cat > "$CONFIG_FILE" << EOF
# DouDian 供应商管理系统配置文件
# 生成时间: $(date)

HOST=0.0.0.0
PORT=$PORT
DB_PATH=$DATA_DIR/doudian.db
JWT_SECRET=$JWT_SECRET
SECRET_PATH=$SECRET_PATH
LOG_LEVEL=info
EOF
    chmod 600 "$CONFIG_FILE"
    echo -e "${GREEN}配置文件已写入: $CONFIG_FILE${PLAIN}"
}

create_service() {
    cat > "$SERVICE_FILE" << EOF
[Unit]
Description=DouDian Supplier Management System
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$APP_DIR
EnvironmentFile=$CONFIG_FILE
ExecStart=$APP_DIR/doudian
Restart=always
RestartSec=5
StandardOutput=append:$LOG_DIR/access.log
StandardError=append:$LOG_DIR/error.log

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    systemctl enable doudian
    echo -e "${GREEN}Systemd 服务已创建${PLAIN}"
}

copy_files() {
    if [ -f "./doudian" ]; then
        cp ./doudian "$APP_DIR/"
        chmod +x "$APP_DIR/doudian"
    else
        echo -e "${YELLOW}警告：未找到 doudian 二进制文件，需要手动编译${PLAIN}"
    fi
    
    if [ -d "./static" ]; then
        cp -r ./static "$APP_DIR/"
    fi
    
    if [ -f "./manage.sh" ]; then
        cp ./manage.sh "$APP_DIR/"
        chmod +x "$APP_DIR/manage.sh"
    fi
}

start_service() {
    echo -e "${BLUE}启动服务...${PLAIN}"
    systemctl start doudian
    sleep 2
    
    if systemctl is-active --quiet doudian; then
        echo -e "${GREEN}服务启动成功${PLAIN}"
    else
        echo -e "${RED}服务启动失败，请检查日志: $LOG_DIR/error.log${PLAIN}"
        systemctl status doudian --no-pager -n 20
        exit 1
    fi
}

get_public_ip() {
    PUBLIC_IP=$(curl -s --max-time 5 https://api.ipify.org 2>/dev/null || echo "YOUR_SERVER_IP")
}

print_complete() {
    get_public_ip
    
    BASE_URL="http://$PUBLIC_IP:$PORT"
    if [ -n "$SECRET_PATH" ]; then
        BASE_URL="$BASE_URL/$SECRET_PATH"
    fi
    
    echo ""
    echo -e "${GREEN}========================================${PLAIN}"
    echo -e "${GREEN}  DouDian 安装完成！${PLAIN}"
    echo -e "${GREEN}========================================${PLAIN}"
    echo ""
    echo -e "  访问地址:  ${BLUE}$BASE_URL/${PLAIN}"
    echo -e "  用户名:    admin"
    echo -e "  密码:      $ADMIN_PASSWORD"
    echo ""
    echo -e "  配置文件:  $CONFIG_FILE"
    echo -e "  日志目录:  $LOG_DIR/"
    echo -e "  数据目录:  $DATA_DIR/"
    echo ""
    echo -e "  管理命令:"
    echo -e "    systemctl start doudian     # 启动"
    echo -e "    systemctl stop doudian      # 停止"
    echo -e "    systemctl restart doudian   # 重启"
    echo -e "    systemctl status doudian    # 状态"
    echo -e "    journalctl -u doudian -f    # 查看日志"
    echo ""
    echo -e "  管理脚本: $APP_DIR/manage.sh"
    echo ""
    echo -e "${YELLOW}提示：请确保防火墙已开放 $PORT 端口${PLAIN}"
    echo ""
}

main() {
    echo ""
    echo -e "${GREEN}========================================${PLAIN}"
    echo -e "${GREEN}  DouDian - 抖店代发供应商管理系统${PLAIN}"
    echo -e "${GREEN}========================================${PLAIN}"
    echo ""
    
    check_root
    detect_os
    get_arch
    check_installed
    IS_UPDATE=$?
    
    install_deps
    get_user_input
    create_dirs
    write_config
    copy_files
    create_service
    start_service
    print_complete
}

main "$@"
