FROM golang:1.20-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /sistem-asetku-backend ./

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /sistem-asetku-backend /sistem-asetku-backend
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/sistem-asetku-backend"]
