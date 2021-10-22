#!/bin/sh

set -e

# check dependencies
if [ -z "${DATABASE_URL}" ]
then
  echo "DATABASE_URL must be set"
  exit 1
fi

if [ ! $(command -v dockerize) ]
then
  echo "dockerize not found"
  exit 1
fi

if [ ! $(command -v python3) ]
then
  echo "python3 not found"
  exit 1
fi

wait_for_database() {
  local db_host=$(python3 -c "import os; from urllib.parse import urlparse; p = urlparse(os.environ['DATABASE_URL']); print(p.hostname)")
  local db_port=$(python3 -c "import os; from urllib.parse import urlparse; p = urlparse(os.environ['DATABASE_URL']); print(p.port)")

  echo "Waiting for database..."
  dockerize -wait "tcp://${db_host}:${db_port}" -timeout 30s
  echo "Detected that database is up!"
}

if [ "$1" = "start" ]
then
  echo "Initializing"
  wait_for_database
  echo "Starting"
  contacts
fi

exec "$@"
