# syntax=docker/dockerfile:1

FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o contacts

FROM alpine:3.13

# Install dependencies such as telnet
RUN apk update && apk add busybox-extras

WORKDIR /

COPY --from=build /app/contacts /contacts

COPY ./db/migrations /migrations
COPY ./docker-entrypoint.sh /docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]
