# DouDian - 抖店一键代发供应商管理系统

基于 Go + React 构建的抖店代发供应商管理系统，支持多店铺管理、供应商管理、商品管理、订单管理和采购单管理。

## 技术栈

- **后端**: Go 1.22+ / Gin / GORM / SQLite (modernc.org/sqlite, 纯 Go 实现)
- **前端**: React 18 / Ant Design 5 / Vite 5
- **认证**: JWT + 安全路径双重防护
- **部署**: Systemd / Docker / 脚本管理

## 功能模块

- **仪表板**: 销售统计、订单概览、数据可视化
- **抖店管理**: 多店铺区分管理
- **订单管理**: 订单CRUD + CSV导出
- **采购单管理**: 采购单CRUD + CSV导出
- **商品管理**: 商品CRUD + CSV导入/导出 + 模板下载
- **供应商管理**: 供应商CRUD + CSV导入/导出 + 模板下载
- **工具**: 数据库备份、系统信息

## 项目结构

```
DouDian/
├── main.go                    # 程序入口
├── internal/
│   ├── config/                # 配置管理
│   ├── database/              # 数据库层 (GORM + SQLite)
│   │   └── model/             # 数据模型
│   ├── eventbus/              # 事件总线
│   ├── service/               # 业务逻辑层
│   ├── util/                  # 工具 (JWT)
│   └── web/
│       ├── controller/        # API 控制器
│       ├── middleware/        # 中间件 (Auth, CORS)
│       └── routes.go           # 路由配置
├── frontend/
│   ├── src/
│   │   ├── api/               # API 请求
│   │   ├── components/        # 通用组件
│   │   ├── layouts/           # 布局
│   │   ├── pages/             # 功能页面
│   │   └── utils/             # 工具函数
│   ├── index.html
│   └── package.json
├── nginx/
│   └── doudian.conf           # Nginx 反向代理配置
├── Makefile                   # 构建自动化
├── install.sh                 # 一键安装脚本
├── manage.sh                  # 管理面板
├── uninstall.sh               # 一键卸载脚本
├── Dockerfile                 # Docker 构建
├── docker-compose.yml         # Docker Compose
├── .env.example               # 环境变量示例
├── go.mod / go.sum            # Go 依赖管理
└── .gitignore
```

## 快速开始

### 本地开发

```bash
# 启动后端
go run .

# 启动前端开发服务器
cd frontend && npm install && npm run dev
```

### 构建生产版本

```bash
make build
```

### Docker 部署

```bash
docker-compose up -d
```

### VPS 一键部署

```bash
# 交叉编译 Linux 版本
make build-linux

# 上传到服务器后执行安装
bash install.sh
```

## 安全配置

- **安全路径**: 通过 `SECRET_PATH` 环境变量设置随机访问路径
- **JWT 认证**: 登录后发放 JWT Token，所有 API 需认证
- **默认账号**: admin / admin123（生产环境请修改）

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| HOST | 0.0.0.0 | 监听地址 |
| PORT | 2095 | 监听端口 |
| DB_PATH | ./instance/doudian.db | 数据库路径 |
| JWT_SECRET | change-me | JWT 密钥 |
| SECRET_PATH | (空) | 安全路径前缀 |
| LOGIN_USERNAME | admin | 登录用户名 |
| LOGIN_PASSWORD | admin123 | 登录密码 |
| LOG_LEVEL | info | 日志级别 |

## License

MIT
