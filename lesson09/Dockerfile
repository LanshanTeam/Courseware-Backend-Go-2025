FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hello

FROM scratch
WORKDIR /app
COPY --from=builder /app/hello .
CMD ["./hello"]
