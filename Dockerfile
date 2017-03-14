FROM alpine:3.5

ADD ./workdir/contacts /usr/bin/contacts
ADD ./db/migrations /migrations

ENTRYPOINT contacts