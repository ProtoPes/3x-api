#!/bin/sh

# This scripts copied from Amnezia client to Docker container to /opt/amnezia and launched every time container starts

echo "Container startup"

if [ ! -f wg0.conf ]; then
    awg-gen-config -c
    scripts/server.sh
    scripts/client.sh
fi

tail -f /dev/null
