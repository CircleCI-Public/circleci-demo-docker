FROM alpine:3.24

ADD ./workdir/contacts /usr/bin/contacts
ADD ./db/migrations /migrations

ENTRYPOINT contacts