#!/bin/sh

set -e

# check dependencies
if [ -z "${DATABASE_HOST}" ]
then
  echo "DATABASE_HOST must be set"
  exit 1
fi

if [ -z "${DATABASE_PORT}" ]
then
  echo "DATABASE_PORT must be set"
  exit 1
fi

wait_for_database() {
  echo "Waiting for database"
  while ! (telnet $DATABASE_HOST $DATABASE_PORT > /dev/null 2>&1); do
    echo -n .
    sleep 1
  done
  echo
}


if [ "$1" = "start" ]
then
  echo "Initializing"
  wait_for_database
  echo "Starting"
  contacts
fi

exec "$@"
