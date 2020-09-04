FROM alpine:3.12.0

ENV M_WORKDIR "/workdir"
ENV M_RESOURCES "/resources"
ENV M_SHARED "/shared"
ENV M_CONFIG "azbi-config.yml"

WORKDIR /workdir
ENTRYPOINT ["make"]

COPY resources /resources
COPY workdir /workdir
