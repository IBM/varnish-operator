#!/bin/sh

# Start k-watcher
while ! varnishadm -t 1 help >/dev/null 2>&1; do
  echo "waiting for varnish to start..."
done
echo "** Starting k-watcher"
/k-watcher
