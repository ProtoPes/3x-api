#!/bin/sh

umask 177
PREFIX_DIR="/opt/amnezia/awg"
cd ${PREFIX_DIR}
source ${PREFIX_DIR}/scripts/values

CLIENT_NAME=$(awg-gen-config -n)
if [ ! -r ${PREFIX_DIR}/unused_ips.txt ]; then
    awg-gen-config -g
fi

## TODO: Remove hardcode
WIREGUARD_CLIENT_IP="10.8.1.$(awg-gen-config -i)/32"

WIREGUARD_CLIENT_PRIVATE_KEY=$(wg genkey)
echo $WIREGUARD_CLIENT_PRIVATE_KEY > ${PREFIX_DIR}/keys/${CLIENT_NAME}_private_key.key

WIREGUARD_CLIENT_PUBLIC_KEY=$(echo $WIREGUARD_CLIENT_PRIVATE_KEY | wg pubkey)
echo $WIREGUARD_CLIENT_PUBLIC_KEY > ${PREFIX_DIR}/keys/${CLIENT_NAME}_public_key.key

WIREGUARD_PSK=$(wg genpsk)
echo $WIREGUARD_PSK > ${PREFIX_DIR}/keys/${CLIENT_NAME}_psk.key
WIREGUARD_SERVER_PUBLIC_KEY=$(cat wireguard_server_public_key.key)


cat > ${PREFIX_DIR}/configs/${CLIENT_NAME}.conf <<EOF
[Interface]
Address = $WIREGUARD_CLIENT_IP
DNS = $PRIMARY_DNS, $SECONDARY_DNS
PrivateKey = $WIREGUARD_CLIENT_PRIVATE_KEY
Jc = $JUNK_PACKET_COUNT
Jmin = $JUNK_PACKET_MIN_SIZE
Jmax = $JUNK_PACKET_MAX_SIZE
S1 = $INIT_PACKET_JUNK_SIZE
S2 = $RESPONSE_PACKET_JUNK_SIZE
H1 = $INIT_PACKET_MAGIC_HEADER
H2 = $RESPONSE_PACKET_MAGIC_HEADER
H3 = $UNDERLOAD_PACKET_MAGIC_HEADER
H4 = $TRANSPORT_PACKET_MAGIC_HEADER

[Peer]
PublicKey = $WIREGUARD_SERVER_PUBLIC_KEY
PresharedKey = $WIREGUARD_PSK
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = $SERVER_IP_ADDRESS:$AWG_SERVER_PORT
PersistentKeepalive = 25
EOF

cat >> ${PREFIX_DIR}/wg0.conf <<EOF

[Peer]
PublicKey = $WIREGUARD_CLIENT_PUBLIC_KEY
PresharedKey = $WIREGUARD_PSK
AllowedIPs = $WIREGUARD_CLIENT_IP
EOF
