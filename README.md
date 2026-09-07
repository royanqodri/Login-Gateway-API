<div align="center">

# 🔐 Login Gateway API

**A production-grade authentication & session gateway built with Go, Gin, PostgreSQL, and Redis.**

Multi-provider login (Email, Google, Facebook, Apple) · JWT auth · Real-time WebSocket notifications · Multi-tenant routing · Swagger docs

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go)
![Gin](https://img.shields.io/badge/Framework-Gin-008ECF?style=flat-square)
![GORM](https://img.shields.io/badge/ORM-GORM-informational?style=flat-square)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-336791?style=flat-square&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Cache-Redis-DC382D?style=flat-square&logo=redis&logoColor=white)
![Swagger](https://img.shields.io/badge/Docs-Swagger-85EA2D?style=flat-square&logo=swagger&logoColor=black)
![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)

</div>

---

## 📖 Overview

**Login Gateway API** is a standalone authentication and session-management service designed to sit in front of a larger system (e.g. a SaaS / multi-tenant platform) and handle everything related to identity: logging users in through multiple channels, issuing and validating JWTs, tracking sessions in Redis, and pushing real-time notifications over WebSockets.

It's built the way I build every backend service for clients — **layered, testable, and boring in the right places**: clear separation between `controller → service → repository → database`, config driven entirely by environment variables, structured logging, and auto-generated API docs so any frontend/mobile team can integrate without asking me a single question.

---

## ✨ Key Features

- 🔑 **Multi-provider authentication** — native email/password login plus OAuth sign-in with **Google**, **Facebook**, and **Apple**.
- 🪪 **JWT-based auth middleware** protecting private routes, with configurable secret & expiry.
- ⚡ **Redis-backed sessions** for fast session lookups and horizontal scalability.
- 🏢 **Multi-tenant routing** — private endpoints are scoped per customer via a `/:customer_no` path segment, so one gateway can serve many tenants.
- 🔔 **Real-time notifications** over WebSocket (subscribe/publish/history), backed by Redis pub/sub for multi-instance broadcasting.
- 📜 **Auto-generated Swagger/OpenAPI docs**, served behind their own auth middleware.
- 🧱 **Clean layered architecture** (`controller → service → repository`) with dedicated request/response/entity models — easy to extend or hand off to another team.
- 🪵 **Structured logging & alerting** — log rotation via `lumberjack` and real-time incident pings to Discord.
- 📦 **FTP & Google Cloud API integration** utilities for file/storage workflows.
- 🌐 **CORS-ready, graceful shutdown, configurable timeouts** — the operational details clients usually forget to ask for, done by default.

---

## 🧰 Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Web framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM / Database | [GORM](https://gorm.io) + PostgreSQL |
| Cache / Sessions | Redis ([go-redis](https://github.com/redis/go-redis)) |
| Auth | JWT ([golang-jwt](https://github.com/golang-jwt/jwt)), OAuth2 (Google, Facebook, Apple) |
| Real-time | WebSocket ([gorilla/websocket](https://github.com/gorilla/websocket)) |
| API docs | Swagger via [swaggo](https://github.com/swaggo/swag) |
| Config | [Viper](https://github.com/spf13/viper) (env-driven) |
| Logging | [Logrus](https://github.com/sirupsen/logrus) + lumberjack (rotation) + Discord webhooks |
| File transfer | [jlaffaye/ftp](https://github.com/jlaffaye/ftp) |

---

## 🏗️ Architecture

```
Client
  │
  ▼
Gin Router  ──►  Middleware (CORS, Auth, Logger, Request-count)
  │
  ▼
Controller  ──►  Service  ──►  Repository  ──►  PostgreSQL (GORM)
  │                  │
  │                  └──►  Redis (sessions, pub/sub)
  ▼
WebSocket Hub  ──►  Redis pub/sub  ──►  Connected clients
```

Every layer has a single responsibility: controllers only parse/validate HTTP, services hold business logic, repositories are the only place that talks to the database. This makes the service straightforward to unit-test and safe to extend without breaking existing flows.

---

## 📡 API Endpoints

### Public
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/login` | Email/password login |
| `POST` | `/login/google` | Sign in with Google |
| `POST` | `/login/facebook` | Sign in with Facebook |
| `POST` | `/login/apple` | Sign in with Apple |
| `POST` | `/user` | Create a user |
| `GET` | `/user` | List users |
| `GET` | `/test` | Health check |
| `GET` | `/swagger/*any` | API documentation (auth-protected) |

### Private — scoped per tenant, requires JWT
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/:customer_no/session-data` | Get the current session's data |

### Real-time notifications (WebSocket)
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/notification/sub` | Subscribe to a notification channel |
| `GET` | `/notification/pub` | Publish a notification |
| `GET` | `/notification/history` | Fetch notification history |
| `POST` | `/notification/post` | Push a notification |

Full request/response schemas are available live via **Swagger UI** once the service is running (see below).

---

## 🚀 Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL
- Redis


## 📂 Project Structure

```
├── cmd/            # entrypoint + generated Swagger docs
├── config/         # env-driven configuration (Viper)
├── controller/     # HTTP handlers (Gin)
├── service/        # business logic
├── repository/     # data access layer (GORM)
├── model/          # entity / request / response structs
├── middleware/      # auth, logging, CORS, request counting
├── database/       # DB connection + migrations
├── websocket/      # real-time notification hub
└── util/           # JWT, hashing, FTP, OAuth helpers, logging, etc.
```

---

## 👋 About Me

I'm **Royan**, a backend developer specializing in **Go** (Gin/Echo, GORM, PostgreSQL) — I build the kind of services shown in this repo: authentication gateways, reporting/scheduling systems, and data pipelines for production platforms. This project reflects how I structure real client work: clean architecture, environment-based config, documented APIs, and operational details (logging, graceful shutdown, timeouts) handled from day one, not bolted on later.

**Looking for a Go backend developer for your project?** I'd be glad to help with:
- REST/gRPC API design & development
- Authentication systems (JWT, OAuth, SSO, multi-tenant)
- Database design & query optimization (PostgreSQL, GORM)
- Real-time features (WebSocket, pub/sub, notifications)
- Refactoring & scaling existing Go services

📩 Feel free to reach out via Upwork or open an issue on this repo — happy to discuss your project.

---

## 📄 License

This project is available under the MIT License — see the [LICENSE](LICENSE) file for details.