FROM golang:1.15.2 as builder
ARG ARG_GO_MODULE_NAME="github.com/epiphany-platform/m-azure-basic-infrastructure"
ENV GO_MODULE_NAME=$ARG_GO_MODULE_NAME
ARG ARG_M_VERSION="unknown"
ENV M_VERSION=$ARG_M_VERSION

RUN mkdir -p $GOPATH/src/$GO_MODULE_NAME
COPY . $GOPATH/src/$GO_MODULE_NAME
WORKDIR $GOPATH/src/$GO_MODULE_NAME
RUN go get -d -v ./...
RUN go get github.com/ahmetb/govvv
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="$(govvv -flags -pkg $GO_MODULE_NAME/cmd -version $M_VERSION)" -x -o /runner $GO_MODULE_NAME

FROM hashicorp/terraform:0.13.2 as initializer

COPY resources /resources
RUN cd /resources/terraform && terraform init

FROM hashicorp/terraform:0.13.2

ENV M_WORKDIR "/workdir"
ENV M_RESOURCES "/resources"
ENV M_SHARED "/shared"

WORKDIR /workdir
ENTRYPOINT ["make"]

RUN apk add --update --no-cache make=4.3-r0 &&\
    wget https://github.com/mikefarah/yq/releases/download/3.3.4/yq_linux_amd64 -O /usr/bin/yq &&\
    chmod +x /usr/bin/yq

ARG ARG_M_VERSION="unknown"
ENV M_VERSION=$ARG_M_VERSION

COPY --from=initializer /resources/ /resources/
COPY workdir /workdir
COPY --from=builder /runner /workdir

ARG ARG_HOST_UID=1000
ARG ARG_HOST_GID=1000
RUN chown -R $ARG_HOST_UID:$ARG_HOST_GID \
    /workdir \
    /resources

USER $ARG_HOST_UID:$ARG_HOST_GID
