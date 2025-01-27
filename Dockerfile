FROM alpine:3.21.2
RUN apk --update upgrade --no-cache \
    && apk add --no-cache dumb-init wireguard-tools-wg
WORKDIR /configs/awg
COPY bin/api-server /usr/bin/3x-api
ENTRYPOINT ["dumb-init", "/usr/bin/3x-api"]
CMD ["init"]
