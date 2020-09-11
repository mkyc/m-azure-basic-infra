FROM hashicorp/terraform:0.12.25 as initializer

COPY resources/terraform /files
RUN TF_DATA_DIR=/wrkdr/.terraform terraform init /files
RUN ls -laR /wrkdr/.terraform/


FROM alpine:3.12.0

ENV M_WORKDIR "/workdir"
ENV M_RESOURCES "/resources"
ENV M_SHARED "/shared"

WORKDIR /workdir
ENTRYPOINT ["make"]

RUN apk add --update --no-cache make=4.3-r0 terraform=0.12.25-r0 git &&\
    wget https://github.com/mikefarah/yq/releases/download/3.3.2/yq_linux_amd64 -O /usr/bin/yq &&\
    chmod +x /usr/bin/yq

ARG ARG_M_VERSION="unknown"
ENV M_VERSION=$ARG_M_VERSION

COPY resources /resources
COPY --from=initializer /wrkdr/.terraform/ /resources/terraform/.terraform/
COPY workdir /workdir

ARG ARG_HOST_UID=1000
ARG ARG_HOST_GID=1000
RUN chown -R $ARG_HOST_UID:$ARG_HOST_GID \
    /workdir \
    /resources

USER $ARG_HOST_UID:$ARG_HOST_GID