FROM go:latest AS builder
WORKDIR /app
COPY ./ /app

go mod tidy
go build /app/cmd/main.go


