# ---------- STAGE 1: BUILD ----------
FROM golang:alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main

# ---------- STAGE 2: RUN ----------
FROM alpine

RUN adduser -S -D -H -h /app appuser
USER appuser

WORKDIR /app
COPY --from=builder /build/main /app/main

EXPOSE 3000
CMD ["./main"]
