# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 3: Final Production Image
FROM alpine:latest
WORKDIR /app
RUN apk add --no-ca-certificates tzdata
COPY --from=backend-builder /app/backend/main .
COPY --from=frontend-builder /app/frontend/dist ./dist

EXPOSE 8080
CMD ["./main"]
