#!/bin/sh

set -e

wait_for_postgres() {
  echo Waiting for Postgres
  while ! (telnet $CONTACTS_DB_HOST $CONTACTS_DB_PORT > /dev/null 2>&1); do
    echo -n .
    sleep 1
  done
  echo
}

if [ "$1" = "start" ]
then
  echo Initializing
  wait_for_postgres
  echo Starting
  contacts
fi

exec "$@"
