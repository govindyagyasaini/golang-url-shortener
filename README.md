# 🚀 URL Shortener Service (Go + Fiber + Redis + Docker)

A **high-performance URL shortener service** built using **Go**, **Fiber**, **Redis**, and **Docker**, featuring **rate limiting**, **custom short URLs**, **TTL-based expiration**, and **containerized deployment**.

This project demonstrates real-world backend concepts such as caching, IP-based rate limiting, Redis usage, and Docker Compose orchestration.

---

## 📌 Features

* 🔗 Shorten long URLs into compact, shareable links
* 🚀 Fast redirection using Redis (in-memory datastore)
* ⏱ URL expiration with configurable TTL
* 🛡 IP-based rate limiting using Redis
* ✨ Optional custom short URLs
* 🐳 Fully containerized using Docker & Docker Compose
* 🧱 Clean project structure (routes, helpers, database)
* ⚡ Built with Fiber v3 (high-performance Go web framework)

---

## 🏗 Tech Stack

| Component        | Technology     |
| ---------------- | -------------- |
| Language         | Go (Golang)    |
| Web Framework    | Fiber v3       |
| Database         | Redis          |
| Containerization | Docker         |
| Orchestration    | Docker Compose |
| Validation       | govalidator    |
| UUID Generation  | google/uuid    |

---

## 📂 Project Structure

```
golang-url-shortener/
│
├── api/
│   ├── main.go
│   ├── routes/
│   │   ├── shorten.go
│   │   └── resolve.go
│   ├── helpers/
│   │   └── helpers.go
│   ├── database/
│   │   └── redis.go
│   ├── Dockerfile
│   └── .env
│
├── db/
│   └── Dockerfile
│
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## ⚙️ How It Works (High Level)

### Redis Databases

* **DB 0** → Stores short URL mappings

  ```
  shortID → originalURL
  ```
* **DB 1** → Stores rate-limiting counters

  ```
  IP → remaining requests
  ```

### Request Flow

1. Client sends a URL to `/api/v1`
2. Rate limit is checked using Redis
3. URL is validated and normalized
4. Short ID is generated (or custom ID used)
5. URL is stored in Redis with TTL
6. Short URL is returned
7. When accessed, short URL redirects to original URL

---

## 🐳 Running the Project (Docker)

### Prerequisites

* Docker
* Docker Compose

### Start the application

From the project root:

```bash
docker compose up --build
```

Services started:

* API → `http://localhost:3000`
* Redis → internal Docker network

---

## 🧪 API Usage

### 🔹 Shorten URL

**Endpoint**

```
POST /api/v1
```

**Request Body**

```json
{
  "url": "https://www.google.com",
  "short": "",
  "expiry": 24
}
```

**Response**

```json
{
  "url": "https://www.google.com",
  "short": "localhost:3000/a1B2c3",
  "expiry": 24,
  "rate_limit": 9,
  "rate_limit_reset": 30
}
```

---

### 🔹 Redirect to Original URL

**Endpoint**

```
GET /:shortID
```

**Example**

```
http://localhost:3000/a1B2c3
```

➡️ Redirects to:

```
https://www.google.com
```

---

## 🛡 Rate Limiting

* Implemented using Redis
* Based on client IP address
* Default quota: **10 requests / 30 minutes**
* Automatically resets after TTL expires

---

## 🔐 Environment Variables

Defined in `.env` file:

```env
DB_ADDR=db:6379
DB_PASS=
APP_PORT=:3000
DOMAIN=localhost:3000
API_QUOTA=10
```

---

## 🧠 Key Learning Outcomes

* Practical Redis usage in Go
* IP-based rate limiting
* TTL-based caching and expiration
* Docker multi-container architecture
* Clean backend project structure
* Fiber framework internals

---

## 🚀 Future Enhancements

* Authentication (JWT)
* Persistent storage fallback (PostgreSQL)
* Admin dashboard
* Analytics (click count per URL)
* API documentation using Swagger
* Unit & integration tests

---

## 👨‍💻 Author

**Govind Yagyasaini**
Backend Engineer | Go | Distributed Systems | Docker

🔗 LinkedIn: www.linkedin.com/in/govindyagyasaini
🐙 GitHub: https://github.com/govindyagyasaini/golang-url-shortener
