# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR app/

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN apk update && apk add netcat-openbsd

COPY . .

RUN go build -o ./build/server ./cmd/server
RUN chmod +x wait-for.sh

EXPOSE $SERVER_PORT


CMD ["./build/server"]

