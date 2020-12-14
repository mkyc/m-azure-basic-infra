ARG VARIANT="buster"

ARG GO_VERSION=1.15
FROM golang:${GO_VERSION}-${VARIANT} AS go

FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}

ENV TERRAFORM_VERSION 0.13.2

# Install Terraform, yq, make, AzureCLI
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && mv terraform /usr/local/bin/ \
    && terraform --version \

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
