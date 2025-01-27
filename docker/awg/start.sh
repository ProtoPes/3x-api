#!/bin/sh

# This scripts copied from Amnezia client to Docker container to /opt/amnezia and launched every time container starts

echo "Container startup"
if [ -z "$AWG_SUBNET_IP" ] || [ -z "$WIREGUARD_SUBNET_CIDR" ]; then
    echo "No subnet environment variable provided! Abort!"
    exit 1
fi
#ifconfig eth0:0 $SERVER_IP_ADDRESS netmask 255.255.255.255 up

# kill daemons in case of restart
wg-quick down wg0

# start daemons if configured
if [ -f wg0.conf ]; then
    wg-quick up wg0

    # Allow traffic on the TUN interface.
    iptables -A INPUT -i wg0 -j ACCEPT
    iptables -A FORWARD -i wg0 -j ACCEPT
    iptables -A OUTPUT -o wg0 -j ACCEPT

    # Allow forwarding traffic only from the VPN.
    iptables -A FORWARD -i wg0 -o eth0 -s "$AWG_SUBNET_IP"/"$WIREGUARD_SUBNET_CIDR" -j ACCEPT
    iptables -A FORWARD -i wg0 -o eth1 -s "$AWG_SUBNET_IP"/"$WIREGUARD_SUBNET_CIDR" -j ACCEPT

    iptables -A FORWARD -m state --state ESTABLISHED,RELATED -j ACCEPT

    iptables -t nat -A POSTROUTING -s "$AWG_SUBNET_IP"/"$WIREGUARD_SUBNET_CIDR" -o eth0 -j MASQUERADE
    iptables -t nat -A POSTROUTING -s "$AWG_SUBNET_IP"/"$WIREGUARD_SUBNET_CIDR" -o eth1 -j MASQUERADE

    find ./wg0.conf | entr -n -p -s 'wg syncconf wg0 <(wg-quick strip wg0)'
    # tail -f /dev/null
else
    echo "No configuration found, exit"
    exit 1
fi
