# --- STEP 1 ---
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o taskservice cmd/api/main.go

# --- STEP 2 ---
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/taskservice /app/taskservice

RUN chmod +x /app/taskservice

EXPOSE 8080

CMD ["/app/taskservice"]