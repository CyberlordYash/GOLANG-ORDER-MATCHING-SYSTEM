# ---------- build stage ----------
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /oms ./main.go

# ---------- final image ----------
FROM gcr.io/distroless/static
COPY --from=builder /oms /oms
EXPOSE 8080
ENTRYPOINT ["/oms"]
