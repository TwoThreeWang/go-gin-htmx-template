# 🚀 Gin + HTMX Project Template

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin-0081CF?style=flat-square)](https://gin-gonic.com/)
[![HTMX](https://img.shields.io/badge/Frontend-HTMX-3366CC?style=flat-square)](https://htmx.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

这是一个基于 **Go + Gin + HTMX** 的现代化 Web 项目模板。它旨在为开发者提供一个开箱即用、模块化且易于扩展的起点，特别适合追求极致性能与简洁前端交互的应用。

---

## ✨ 核心特性

- **⚡ 高性能后端**: 基于 Gin 框架，提供极速的路由响应与中间件支持。
- **🎨 声明式前端**: 集成 HTMX，无需编写复杂的 JavaScript 即可实现局部刷新、异步提交等现代交互。
- **🏗️ 模块化架构**: 严格的包结构划分（Handler, Service, Repository），确保代码的可维护性与可测试性。
- **🔐 安全认证**: 
  - 基于 Http-Only Cookie 的 JWT 身份验证。
  - **滑动续期 (Sliding Expiration)**：自动刷新活跃用户的 Token 有效期。
- **📊 统一响应封装**: 标准化的 API 返回格式，简化前后端对接。
- **🧩 模板引擎优化**: 使用 `multitemplate` 解决 Go 原生模板继承痛点，支持全局数据注入 (`RenderData`)。
- **🗄️ 数据库集成**: 预配置 GORM + PostgreSQL，支持自动迁移。

---

## 📂 目录结构

```text
gin-htmx-template/
├── cmd/server/         # 程序入口 (main.go)
├── internal/           # 核心业务逻辑 (私有)
│   ├── config/         # 配置管理 (Environment variables)
│   ├── handler/        # 控制层 (HTTP Handlers & RenderData)
│   ├── middleware/     # 中间件 (Auth, Logging, Recovery)
│   ├── model/          # 数据模型 (GORM Models)
│   ├── repository/     # 持久层 (Database DAL)
│   ├── router/         # 路由定义 & 模板加载 (LoadTemplates)
│   ├── service/        # 业务服务层 (Optional)
│   └── utils/          # 工具类 (Response, Logger)
├── web/                # 前端资源
│   ├── static/         # 静态文件 (CSS, JS, Images)
│   └── templates/      # 视图模板 (Layouts, Pages, Partials)
├── .env.example        # 环境变量模板
└── go.mod              # 依赖管理
```

---

## 🛠️ 快速开始

### 1. 环境准备
- **Go**: v1.21 或更高版本
- **PostgreSQL**: v15+

### 2. 克隆与安装
```bash
git clone https://github.com/TwoThreeWang/go-gin-htmx-template.git
cd go-gin-htmx-template
go mod tidy
```

### 3. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，配置您的数据库连接信息
```

### 4. 运行项目
```bash
go run ./cmd/server/main.go
```
访问 [http://localhost:5007](http://localhost:5007) 即可看到示例页面。

---

## 📖 示例页面展示

- **首页 (Home)**: 基础布局与公共数据注入演示。
- **用户列表 (HTMX)**: 演示 `hx-get` 局部刷新表格。
- **表单提交 (HTMX)**: 演示 `hx-post` 异步提交与成功反馈。
- **API 接口**: 演示 `utils.Success` 统一响应格式。

---

## 🧩 模板系统

项目使用 `multitemplate` 解决 Go 原生模板继承问题，支持全局数据注入 (`RenderData`)。

### 模板目录结构
```text
web/templates/
├── layouts/      # 布局模板 (base.html)
├── pages/        # 页面模板 (home.html, about.html, users.html, ...)
└── partials/     # 局部模板 (users-list.html, contact-form.html, ...)
```

### 模板命名规范

| 类型 | 文件路径 | 注册名称 | 使用示例 |
|------|----------|----------|----------|
| 页面模板 | `pages/home.html` | `home` | `c.HTML(200, "home", data)` |
| 局部模板 | `partials/users-list.html` | `partials/users-list` | `c.HTML(200, "partials/users-list", data)` |

### 添加新页面

1. 在 `web/templates/pages/` 下创建 HTML 文件（如 `new-page.html`）
2. 模板会自动注册，名称为文件名（不带扩展名）
3. 在 handler 中调用：`c.HTML(200, "new-page", data)`

### 添加新的局部模板

1. 在 `web/templates/partials/` 下创建 HTML 文件
2. 模板自动注册为 `partials/文件名`（不带扩展名）
3. 用于 HTMX 局部刷新：`c.HTML(200, "partials/xxx", data)`

---

## 🛡️ 开源协议

本项目采用 [MIT License](LICENSE) 开源协议。

---

## 🙏 鸣谢

- 感谢 [Moovie](https://github.com/TwoThreeWang/Moovie) 项目提供的架构灵感。
- 感谢 [Gin](https://github.com/gin-gonic/gin) 与 [HTMX](https://htmx.org/) 社区的卓越贡献。
