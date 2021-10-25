# syntax=docker/dockerfile:1

FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -v -o contacts

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/contacts /contacts

ADD ./db/migrations /migrations
ADD docker-entrypoint.sh /

EXPOSE 8080

# ENV DOCKERIZE_VERSION v0.6.1
# RUN apt-get update && apt-get install -y wget
# ENV DOCKERIZE_VERSION v0.6.1
# RUN wget "https://github.com/jwilder/dockerize/releases/download/${DOCKERIZE_VERSION}/dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz" \
#     && tar -C /usr/local/bin -xzvf "dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz" \
#     && rm "dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz"

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["start"]
