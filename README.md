# 🔗 URL Shortener Service

A production-ready URL shortener microservice built with **Go**, **Fiber v3**, and **Redis**, featuring IP-based rate limiting, custom short URLs, TTL-based expiration, and full Docker containerization.

## 📋 Overview

This project demonstrates a scalable URL shortening service similar to bit.ly or TinyURL, showcasing real-world backend engineering concepts including:

- **Distributed caching** with Redis
- **IP-based rate limiting** for API protection
- **Multi-database architecture** using Redis logical databases
- **RESTful API design** with proper error handling
- **Container orchestration** with Docker Compose
- **Clean architecture** with separation of concerns

Perfect for learning modern Go backend development, Redis usage patterns, and microservices deployment strategies.

---

## ✨ Features

### Core Functionality
- ✅ **URL Shortening**: Convert long URLs into compact, shareable links
- ✅ **Fast Redirection**: Sub-millisecond redirects using Redis in-memory storage
- ✅ **Custom Short URLs**: Optional user-defined short identifiers
- ✅ **TTL-Based Expiration**: Automatic URL cleanup with configurable expiry (default: 24 hours)
- ✅ **Click Tracking**: Counter increments on each redirect (analytics foundation)

### Security & Performance
- 🛡️ **IP-Based Rate Limiting**: Configurable request quota (default: 10 requests per 30 minutes)
- 🔒 **URL Validation**: Prevents invalid URLs and self-referencing links
- ⚡ **High Performance**: Fiber framework with minimal latency
- 🐳 **Production-Ready**: Multi-stage Docker builds for optimized container size

### Developer Experience
- 📝 **Clean Code Structure**: Modular design with routes, helpers, and database layers
- 🔧 **Environment Configuration**: Easy deployment with `.env` files
- 📊 **Request Logging**: Built-in middleware for debugging
- 🚀 **One-Command Deployment**: Docker Compose orchestration

---

## 🏗️ Architecture

### System Design

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────┐
│   Fiber API Server (Port 3000)  │
│  ┌──────────────────────────┐   │
│  │  Rate Limiter Middleware │   │
│  └──────────────────────────┘   │
│  ┌──────────────────────────┐   │
│  │   Routes & Handlers      │   │
│  │  • POST /api/v1          │   │
│  │  • GET /:url             │   │
│  └──────────────────────────┘   │
└────────┬─────────────────────────┘
         │
         ▼
┌─────────────────────────┐
│   Redis (Port 6379)     │
│  ┌──────────────────┐   │
│  │  DB 0: URL Store │   │
│  │  shortID → URL   │   │
│  └──────────────────┘   │
│  ┌──────────────────┐   │
│  │  DB 1: Rate Limit│   │
│  │  IP → Quota      │   │
│  └──────────────────┘   │
└─────────────────────────┘
```

### Redis Database Strategy

**DB 0 - URL Mappings**
```
Key: <shortID>
Value: <originalURL>
TTL: User-defined (default: 24 hours)
```

**DB 1 - Rate Limiting & Analytics**
```
Key: <client_IP>
Value: Remaining quota
TTL: 30 minutes

Key: "counter"
Value: Total redirects
TTL: None (persistent)
```

---

## 🛠️ Tech Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Language** | Go 1.21+ | High-performance backend |
| **Web Framework** | Fiber v3 | Fast HTTP routing & middleware |
| **Database** | Redis 7.x | In-memory data store |
| **Containerization** | Docker | Application packaging |
| **Orchestration** | Docker Compose | Multi-container management |
| **URL Validation** | govalidator | Input sanitization |
| **UUID Generation** | google/uuid | Short ID generation |
| **Environment Config** | godotenv | Configuration management |

---

## 📂 Project Structure

```
golang-url-shortener/
│
├── api/                          # Main application
│   ├── main.go                   # Application entry point
│   ├── routes/
│   │   ├── shorten.go           # POST /api/v1 - URL shortening logic
│   │   └── resolve.go           # GET /:url - Redirection logic
│   ├── helpers/
│   │   └── helpers.go           # URL validation & normalization
│   ├── database/
│   │   └── redis.go             # Redis client factory
│   ├── Dockerfile               # Multi-stage build for API
│   ├── .env                     # Environment configuration
│   ├── go.mod                   # Go module dependencies
│   └── go.sum                   # Dependency checksums
│
├── db/
│   └── Dockerfile               # Redis container configuration
│
├── docker-compose.yml           # Service orchestration
└── README.md                    # Project documentation
```

---

## 🚀 Getting Started

### Prerequisites

- **Docker** (20.10+)
- **Docker Compose** (2.0+)
- **Go** 1.21+ (for local development)

### Installation & Running

#### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/govindyagyasaini/golang-url-shortener.git
cd golang-url-shortener

# Start all services
docker compose up --build
```

**Services Started:**
- API Server: `http://localhost:3000`
- Redis: Internal network (port 6379)

#### Option 2: Local Development

```bash
# Install dependencies
cd api
go mod download

# Start Redis separately
docker run -d -p 6379:6379 redis:alpine

# Configure environment
cp .env.example .env

# Run the application
go run main.go
```

---

## 📡 API Documentation

### 1️⃣ Shorten URL

Create a shortened URL with optional custom identifier and expiry.

**Endpoint:**
```http
POST /api/v1
Content-Type: application/json
```

**Request Body:**
```json
{
  "url": "https://www.example.com/very/long/url/path",
  "short": "mylink",           // Optional: custom short ID
  "expiry": 48                  // Optional: hours (default: 24)
}
```

**Success Response (200 OK):**
```json
{
  "url": "https://www.example.com/very/long/url/path",
  "short": "localhost:3000/mylink",
  "expiry": 48,
  "rate_limit": 9,              // Remaining requests
  "rate_limit_reset": 1800      // Seconds until reset
}
```

**Error Responses:**

| Status | Error | Cause |
|--------|-------|-------|
| `400 Bad Request` | `"invalid URL"` | Malformed or invalid URL |
| `400 Bad Request` | `"invalid domain"` | Self-referencing URL |
| `400 Bad Request` | `"cannot parse JSON"` | Invalid JSON payload |
| `403 Forbidden` | `"short URL already exists"` | Custom short ID in use |
| `429 Too Many Requests` | `"rate limit exceeded"` | Quota exhausted |
| `500 Internal Server Error` | `"cannot save URL"` | Redis error |

**Example with cURL:**
```bash
curl -X POST http://localhost:3000/api/v1 \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.github.com/govindyagyasaini",
    "short": "github",
    "expiry": 72
  }'
```

---

### 2️⃣ Resolve Short URL

Redirect to the original URL using the short identifier.

**Endpoint:**
```http
GET /:shortID
```

**Example:**
```bash
curl -L http://localhost:3000/mylink
# Redirects to https://www.example.com/very/long/url/path
```

**Response:**
```http
HTTP/1.1 301 Moved Permanently
Location: https://www.example.com/very/long/url/path
```

**Error Responses:**

| Status | Error | Cause |
|--------|-------|-------|
| `404 Not Found` | `"short url not found"` | Invalid or expired short ID |
| `500 Internal Server Error` | `"database error"` | Redis connection error |

**Behavior:**
- Increments global redirect counter in Redis DB 1
- Returns `301 Moved Permanently` for SEO benefits
- Automatically expires after configured TTL

---

## 🛡️ Rate Limiting

### How It Works

- **Scope**: Per client IP address
- **Default Quota**: 10 requests per 30 minutes
- **Storage**: Redis DB 1 with TTL
- **Reset**: Automatic after 30-minute window

### Configuration

Edit `.env` file:
```env
API_QUOTA=10  # Requests per window
```

### Example Flow

```bash
# Request 1-10: Success
curl -X POST http://localhost:3000/api/v1 -d '{"url":"https://example.com"}'
# Response: "rate_limit": 9, 8, 7... 0

# Request 11: Rate limited
curl -X POST http://localhost:3000/api/v1 -d '{"url":"https://example.com"}'
# Response: 429 Too Many Requests
# {
#   "error": "rate limit exceeded",
#   "rate_limit_reset": 25.5  // minutes remaining
# }
```

---

## ⚙️ Configuration

### Environment Variables

Create `api/.env` file:

```env
# Redis Configuration
DB_ADDR=db:6379              # Redis host (use 'localhost:6379' for local dev)
DB_PASS=                     # Redis password (empty for no auth)

# Application Configuration
APP_PORT=:3000               # API server port
DOMAIN=localhost:3000        # Base domain for short URLs

# Rate Limiting
API_QUOTA=10                 # Max requests per IP per window
```

### Docker Compose Configuration

The `docker-compose.yml` orchestrates two services:

```yaml
services:
  api:                        # Go application
    build: ./api
    ports:
      - "3000:3000"
    depends_on:
      - db
    env_file:
      - ./api/.env

  db:                         # Redis database
    build: ./db
    ports:
      - "6379:6379"
```

---

## 🧠 Key Implementation Details

### URL Validation & Normalization

**helpers/helpers.go**
```go
// Ensures URL has HTTP/HTTPS scheme
func EnforceHTTP(url string) string

// Prevents self-referencing URLs (e.g., shortening your own domain)
func RemoveDomainError(url string) bool
```

### Short ID Generation

```go
// Auto-generated: First 6 chars of UUID v4
id := uuid.New().String()[:6]  // Example: "a1b2c3"

// Custom: User-provided identifier
id := body.CustomShort          // Example: "mylink"
```

### Redis Client Factory

**database/redis.go**
```go
// Creates Redis client for specific database
func CreateClient(dbNo int) *redis.Client
```

Usage:
- `CreateClient(0)` → URL storage
- `CreateClient(1)` → Rate limiting

---

## 🧪 Testing

### Manual Testing

**1. Shorten a URL:**
```bash
curl -X POST http://localhost:3000/api/v1 \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'
```

**2. Access the short URL:**
```bash
curl -L http://localhost:3000/<shortID>
```

**3. Test rate limiting:**
```bash
# Run 11 times quickly
for i in {1..11}; do
  curl -X POST http://localhost:3000/api/v1 \
    -H "Content-Type: application/json" \
    -d '{"url":"https://example.com"}'
done
```

### Redis Verification

```bash
# Connect to Redis container
docker exec -it <redis_container_id> redis-cli

# Check URL mappings (DB 0)
SELECT 0
KEYS *
GET <shortID>

# Check rate limits (DB 1)
SELECT 1
KEYS *
GET <your_ip>
```

---

## 🐛 Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| `connection refused` | Redis not running | `docker compose up db` |
| `invalid URL` error | Missing http:// | URLs auto-prefixed with http:// |
| Rate limit immediately | Clock skew | Check system time synchronization |
| Port 3000 in use | Another service running | Change `APP_PORT` in `.env` |
| `short URL already exists` | Collision | Use custom short ID or retry |

---

## 🚀 Future Enhancements

### Priority Features
- [ ] **Authentication**: JWT-based API key system
- [ ] **Analytics Dashboard**: Click tracking, geographic data, referrer stats
- [ ] **Persistent Storage**: PostgreSQL fallback for critical URLs
- [ ] **QR Code Generation**: Auto-generate QR codes for short URLs
- [ ] **Custom Domains**: Support user-defined domains (e.g., `go.mysite.com/abc`)

### Technical Improvements
- [ ] **Comprehensive Testing**: Unit & integration tests with testify
- [ ] **Monitoring**: Prometheus metrics + Grafana dashboards
- [ ] **API Documentation**: OpenAPI/Swagger specification
- [ ] **Graceful Shutdown**: Proper signal handling
- [ ] **Health Checks**: `/health` and `/ready` endpoints
- [ ] **Kubernetes Deployment**: Helm charts for production deployment

---

## 📚 Learning Outcomes

This project demonstrates:

✅ **Backend Engineering**
- RESTful API design with proper HTTP status codes
- Middleware pattern implementation
- Error handling and validation strategies

✅ **Database Design**
- Redis multi-database architecture
- TTL-based cache invalidation
- Counter patterns for analytics

✅ **DevOps Practices**
- Multi-stage Docker builds for optimization
- Docker Compose for local development
- Environment-based configuration

✅ **Production Considerations**
- Rate limiting for API protection
- Input validation and sanitization
- Logging and observability

---

## 👤 Author

**Govind Yagyasaini**

Backend Developer | Go • Cloud • Distributed Systems

- 🔗 LinkedIn: [linkedin.com/in/govindyagyasaini](https://www.linkedin.com/in/govindyagyasaini)
- 🐙 GitHub: [@govindyagyasaini](https://github.com/govindyagyasaini)

---

## 🙏 Acknowledgments

- [Fiber](https://github.com/gofiber/fiber) - Express-inspired web framework for Go
- [go-redis](https://github.com/go-redis/redis) - Type-safe Redis client for Go
- [govalidator](https://github.com/asaskevich/govalidator) - Package of validators and sanitizers
- [godotenv](https://github.com/joho/godotenv) - Go port of Ruby's dotenv library

---

## 📞 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
