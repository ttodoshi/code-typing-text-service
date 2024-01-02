FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src/app

COPY go.mod go.sum ./
COPY internal internal
COPY pkg pkg
COPY cmd cmd

RUN go mod tidy

RUN go build -o ./bin/app ./cmd/main/main.go

FROM alpine:latest

WORKDIR /usr/local/src/app

COPY --from=0 /usr/local/src/app/bin/app .

CMD ["./app"]