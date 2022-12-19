FROM go:latest AS builder
WORKDIR /app
COPY ./ /app

RUN go mod tidy
RUN go build /app/cmd/main.go


