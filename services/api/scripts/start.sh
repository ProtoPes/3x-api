#!/bin/sh

# This scripts copied from Amnezia client to Docker container to /opt/amnezia and launched every time container starts

echo "Container startup"

# If config file do not exists, create one
if [ ! -f wg0.conf ]; then
    3x-api -c
    scripts/server.sh
    scripts/client.sh
fi

tail -f /dev/null
