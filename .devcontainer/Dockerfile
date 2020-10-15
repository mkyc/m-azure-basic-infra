ARG VARIANT="buster"

ARG GO_VERSION=1.15
FROM golang:${GO_VERSION}-${VARIANT} AS go

FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}

ENV TERRAFORM_VERSION 0.13.2
ENV YQ_VERSION=3.3.4

# Install Terraform, yq, make, AzureCLI
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && mv terraform /usr/local/bin/ \
    && terraform --version \

    && wget https://github.com/mikefarah/yq/releases/download//${YQ_VERSION}/yq_linux_amd64 -O /usr/bin/yq \
    && chmod +x /usr/bin/yq \

    && apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends make \

    && wget -O - https://aka.ms/InstallAzureCLIDeb | bash

# Setup and configure GO for writting Terratests
USER root
COPY --from=go /usr/local/go /usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH \
    CGO_ENABLED=0 \
    GO111MODULE=on
RUN go get -v \
    # Base Go tools needed for VS code Go extension
    golang.org/x/tools/gopls \
    golang.org/x/tools/cmd/guru \
    golang.org/x/tools/cmd/gorename \
    github.com/go-delve/delve/cmd/dlv \
    2>&1 && \
    rm -rf $GOPATH/pkg/* $GOPATH/src/* /root/.cache/go-build && \
    chown -R ${USER_UID}:${USER_GID} $GOPATH && \
    chmod -R 777 $GOPATH

USER ${USER}
