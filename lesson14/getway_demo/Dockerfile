FROM golang:1.25.0 AS build

ARG TARGET_BINARY

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN service="${TARGET_BINARY##*/}" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app "./cmd/${service}"

FROM scratch

WORKDIR /app
COPY --from=build /out/app /app/app

EXPOSE 8080

CMD ["/app/app"]
