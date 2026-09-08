#!/bin/bash
# ==============================================================================
# DouDian - 抖店一键代发供应商管理系统 - 一键卸载脚本
#
# 用法:
#   bash uninstall.sh
#
# 选项:
#   NONINTERACTIVE=1    非交互模式（不询问，直接卸载）
#   KEEP_DATA=1         保留数据库和备份（仅删除程序文件）
# ==============================================================================

red='\033[0;31m'
green='\033[0;32m'
blue='\033[0;34m'
yellow='\033[0;33m'
plain='\033[0m'

APP_NAME="doudian"
APP_DIR="/opt/doudian"
SERVICE_FILE="/etc/systemd/system/doudian.service"
ENV_FILE="/etc/default/doudian"
LOG_DIR="/var/log/doudian"
DATA_DIR="/opt/doudian/data"
NGINX_CONF="/etc/nginx/conf.d/doudian.conf"
NGINX_SITE="/etc/nginx/sites-enabled/doudian"

check_root() {
    if [[ $EUID -ne 0 ]]; then
        echo -e "${red}请使用 root 用户运行此脚本${plain}"
        echo -e "${yellow}  sudo bash uninstall.sh${plain}"
        exit 1
    fi
}

confirm_uninstall() {
    if [[ "${NONINTERACTIVE:-0}" == "1" ]]; then
        return 0
    fi
    
    echo -e "${yellow}========================================${plain}"
    echo -e "${yellow}  卸载确认${plain}"
    echo -e "${yellow}========================================${plain}"
    echo ""
    echo -e "  将删除以下内容:"
    echo -e "    1. systemd 服务 (${SERVICE_FILE})"
    echo -e "    2. 应用目录 (${APP_DIR})"
    echo -e "    3. 配置文件 (${ENV_FILE})"
    echo -e "    4. 日志目录 (${LOG_DIR})"
    echo -e "    5. Nginx 配置 (如存在)"
    echo -e "    6. 防火墙端口规则"
    echo ""
    
    if [[ "${KEEP_DATA:-0}" == "1" ]]; then
        echo -e "  ${green}KEEP_DATA=1: 将保留数据库和备份文件${plain}"
    else
        echo -e "  ${red}将删除所有数据库和备份文件${plain}"
    fi
    
    echo ""
    read -p "确认卸载? (y/N): " confirm
    if [[ "${confirm}" != "y" && "${confirm}" != "Y" ]]; then
        echo -e "${yellow}已取消卸载${plain}"
        exit 0
    fi
}

stop_service() {
    echo -e "${blue}[1/6] 停止服务...${plain}"
    
    if systemctl is-active --quiet ${APP_NAME} 2>/dev/null; then
        systemctl stop ${APP_NAME}
        echo -e "  ${green}服务已停止${plain}"
    else
        echo -e "  ${yellow}服务未在运行${plain}"
    fi
}

disable_service() {
    echo -e "${blue}[2/6] 禁用并删除 systemd 服务...${plain}"
    
    if systemctl is-enabled --quiet ${APP_NAME} 2>/dev/null; then
        systemctl disable ${APP_NAME}
        echo -e "  ${green}服务已禁用${plain}"
    else
        echo -e "  ${yellow}服务未启用${plain}"
    fi
    
    if [[ -f "${SERVICE_FILE}" ]]; then
        rm -f "${SERVICE_FILE}"
        systemctl daemon-reload
        echo -e "  ${green}服务文件已删除${plain}"
    else
        echo -e "  ${yellow}服务文件不存在${plain}"
    fi
}

backup_before_remove() {
    if [[ "${KEEP_DATA:-0}" == "1" ]]; then
        echo -e "${blue}[3/6] 保留数据模式 - 备份数据到 /tmp/doudian-backup...${plain}"
        local backup_dir="/tmp/doudian-backup-$(date +%Y%m%d%H%M%S)"
        mkdir -p "${backup_dir}"
        
        if [[ -d "${DATA_DIR}" ]]; then
            cp -rf "${DATA_DIR}" "${backup_dir}/data"
            echo -e "  ${green}数据库已备份到 ${backup_dir}/data${plain}"
        fi
        
        if [[ -f "${ENV_FILE}" ]]; then
            cp -f "${ENV_FILE}" "${backup_dir}/doudian.env"
            echo -e "  ${green}配置已备份到 ${backup_dir}/doudian.env${plain}"
        fi
        
        echo -e "  ${yellow}数据备份路径: ${backup_dir}${plain}"
    else
        echo -e "${blue}[3/6] 跳过数据备份${plain}"
    fi
}

remove_app_files() {
    echo -e "${blue}[4/6] 删除应用文件...${plain}"
    
    if [[ "${KEEP_DATA:-0}" == "1" ]]; then
        if [[ -d "${APP_DIR}" ]]; then
            find "${APP_DIR}" -mindepth 1 ! -path "${DATA_DIR}*" ! -path "${DATA_DIR}" -exec rm -rf {} + 2>/dev/null
            echo -e "  ${green}程序文件已删除，数据目录已保留${plain}"
        fi
    else
        if [[ -d "${APP_DIR}" ]]; then
            rm -rf "${APP_DIR}"
            echo -e "  ${green}应用目录已删除: ${APP_DIR}${plain}"
        else
            echo -e "  ${yellow}应用目录不存在${plain}"
        fi
    fi
}

remove_config_and_logs() {
    echo -e "${blue}[5/6] 删除配置和日志...${plain}"
    
    if [[ -f "${ENV_FILE}" ]]; then
        rm -f "${ENV_FILE}"
        echo -e "  ${green}配置文件已删除: ${ENV_FILE}${plain}"
    else
        echo -e "  ${yellow}配置文件不存在${plain}"
    fi
    
    if [[ -d "${LOG_DIR}" ]]; then
        rm -rf "${LOG_DIR}"
        echo -e "  ${green}日志目录已删除: ${LOG_DIR}${plain}"
    else
        echo -e "  ${yellow}日志目录不存在${plain}"
    fi
}

remove_nginx_and_firewall() {
    echo -e "${blue}[6/6] 清理 Nginx 配置和防火墙...${plain}"
    
    local nginx_changed=0
    
    for conf in "${NGINX_CONF}" "${NGINX_SITE}"; do
        if [[ -f "${conf}" ]]; then
            rm -f "${conf}"
            echo -e "  ${green}Nginx 配置已删除: ${conf}${plain}"
            nginx_changed=1
        fi
    done
    
    if [[ ${nginx_changed} -eq 1 ]]; then
        if command -v nginx &>/dev/null; then
            nginx -t 2>/dev/null && systemctl reload nginx 2>/dev/null
            echo -e "  ${green}Nginx 已重载${plain}"
        fi
    else
        echo -e "  ${yellow}未找到 Nginx 配置${plain}"
    fi
    
    local port="${PORT:-2095}"
    
    if command -v ufw &>/dev/null; then
        ufw delete allow ${port}/tcp 2>/dev/null
        echo -e "  ${green}ufw 已移除端口 ${port} 规则${plain}"
    elif command -v firewall-cmd &>/dev/null; then
        firewall-cmd --permanent --remove-port=${port}/tcp 2>/dev/null
        firewall-cmd --reload 2>/dev/null
        echo -e "  ${green}firewalld 已移除端口 ${port} 规则${plain}"
    elif command -v iptables &>/dev/null; then
        iptables -D INPUT -p tcp --dport ${port} -j ACCEPT 2>/dev/null
        echo -e "  ${green}iptables 已移除端口 ${port} 规则${plain}"
    else
        echo -e "  ${yellow}未检测到防火墙工具${plain}"
    fi
}

print_uninstall_result() {
    echo ""
    echo -e "${green}========================================${plain}"
    echo -e "${green}  DouDian 卸载完成！${plain}"
    echo -e "${green}========================================${plain}"
    echo ""
    
    if [[ "${KEEP_DATA:-0}" == "1" ]]; then
        echo -e "  ${yellow}数据保留路径: /tmp/doudian-backup-*${plain}"
        echo -e "  ${yellow}如需重新安装: bash install.sh${plain}"
    else
        echo -e "  ${green}所有组件已完全移除${plain}"
    fi
    echo ""
}

main() {
    echo -e "${green}========================================${plain}"
    echo -e "${green}  DouDian 一键卸载脚本${plain}"
    echo -e "${green}========================================${plain}"
    echo ""
    
    check_root
    confirm_uninstall
    
    stop_service
    disable_service
    backup_before_remove
    remove_app_files
    remove_config_and_logs
    remove_nginx_and_firewall
    
    print_uninstall_result
}

main "$@"
