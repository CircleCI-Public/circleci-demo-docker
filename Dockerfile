FROM alpine:3.13

ADD ./workdir/contacts /usr/bin/contacts
ADD ./db/migrations /migrations

ENTRYPOINT contacts
