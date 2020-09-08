FROM alpine:3.12.0

ENV M_WORKDIR "/workdir"
ENV M_RESOURCES "/resources"
ENV M_SHARED "/shared"
ENV M_CONFIG_NAME "azbi-config.yml"

WORKDIR /workdir
ENTRYPOINT ["make"]

RUN apk add --update --no-cache make=4.3-r0 terraform=0.12.25-r0 curl &&\
    wget $(curl -s https://api.github.com/repos/mikefarah/yq/releases/latest | grep browser_download_url | grep linux_amd64 | cut -d '"' -f 4) -O /usr/bin/yq &&\
    chmod +x /usr/bin/yq

ARG ARG_M_VERSION="unknown"
ENV M_VERSION=$ARG_M_VERSION

COPY resources /resources
COPY workdir /workdir
