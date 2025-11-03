FROM golang:1.24.0 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

FROM debian:stable-slim
WORKDIR /app
COPY --from=build /app/api /app/api
ENTRYPOINT ["/app/api"]

