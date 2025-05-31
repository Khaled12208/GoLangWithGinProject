#!/bin/sh

set -e

host="$1"
shift
cmd="$@"

echo "Waiting for MySQL to be ready..."
until mysql -h "$host" -u root -pmysecretpassword --ssl-mode=DISABLED -e 'SELECT 1' >/dev/null 2>&1; do
  echo "MySQL is unavailable - sleeping"
  sleep 1
done

echo "MySQL is up - executing command"
exec $cmd 