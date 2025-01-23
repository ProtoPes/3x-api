FROM golang:1.23.5-alpine3.21 AS go
COPY src /src
WORKDIR /src
RUN go build -v -o /usr/bin

FROM alpine:3.21.2
RUN apk --update upgrade --no-cache \
    && apk add --no-cache dumb-init wireguard-tools-wg
WORKDIR /opt/amnezia/awg
COPY scripts /opt/amnezia/awg/scripts
RUN chmod +x /opt/amnezia/awg/scripts/*
COPY --from=go /usr/bin/awg-gen-config /usr/bin/awg-gen-config
CMD ["dumb-init", "scripts/start.sh"]

