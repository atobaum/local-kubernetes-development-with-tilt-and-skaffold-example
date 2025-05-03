FROM golang:1.24 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/main main.go

FROM debian:bullseye-slim
COPY --from=builder /app/build/main /app/main
ENTRYPOINT ["/app/main"]