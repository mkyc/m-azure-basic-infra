FROM alpine:3.12.0

ENV M_WORKDIR "/workdir"
ENV M_RESOURCES "/resources"
ENV M_SHARED "/shared"
ENV M_CONFIG "azbi-config.yml"

WORKDIR /workdir
ENTRYPOINT ["make"]

RUN apk add --update --no-cache make=4.3-r0

ARG ARG_M_VERSION="unknown"
ENV M_VERSION=$ARG_M_VERSION

COPY resources /resources
COPY workdir /workdir
