FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/server
RUN go build -o /app/server

WORKDIR /app/cmd/agent
RUN go build -o /app/agent

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/agent .