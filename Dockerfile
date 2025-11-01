ARG PHP_VERSION

FROM ghcr.io/onlishop/onlishop-cli-base:${PHP_VERSION}

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/onlishop-cli /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/onlishop-cli"]
CMD ["--help"]
