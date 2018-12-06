FROM alpine:3.5

COPY ./workdir/contacts /usr/bin/contacts
COPY ./db/migrations /migrations

ENTRYPOINT contacts
