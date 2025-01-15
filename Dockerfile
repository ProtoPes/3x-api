FROM golang:1.23.4 as go
COPY src /src
WORKDIR /src
RUN go build -v -o /usr/bin

FROM alpine:latest

RUN mkdir -p /opt/amnezia/awg
COPY scripts /opt/amnezia/awg/scripts
RUN chmod +x /opt/amnezia/awg/scripts
COPY --from=go /usr/bin/awg-gen-config /usr/bin/awg-gen-config

