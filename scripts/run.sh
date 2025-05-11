#!/bin/sh
set -e

if [ ! -f ./bin/server/subitocart ]; then
  echo "Binary not found. Run scripts/build.sh first."
  exit 1
fi

echo "Starting server on port 9090..."
./bin/server/subitocart