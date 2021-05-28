# golang builder
FROM golang:1.16.2 as builder
ARG ARG_GO_MODULE_NAME="github.com/epiphany-platform/m-azure-basic-infrastructure"
ENV GO_MODULE_NAME=$ARG_GO_MODULE_NAME
ARG ARG_M_VERSION="dev"
ENV M_VERSION=$ARG_M_VERSION
RUN mkdir -p $GOPATH/src/$GO_MODULE_NAME
COPY . $GOPATH/src/$GO_MODULE_NAME
WORKDIR $GOPATH/src/$GO_MODULE_NAME
RUN cd /tmp && go get github.com/ahmetb/govvv
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w $(govvv -flags -pkg $GO_MODULE_NAME/cmd -version $M_VERSION)" -o /m-azure-basic-infrastructure $GO_MODULE_NAME

# terraform init
FROM hashicorp/terraform:0.13.7 as initializer

COPY resources /resources
RUN cd /resources/terraform && terraform init

# main
FROM hashicorp/terraform:0.13.7

ENV RESOURCES "/resources"
ENV SHARED "/shared"

WORKDIR /workdir
ENTRYPOINT ["/workdir/m-azure-basic-infrastructure"]

COPY --from=initializer /resources/ /resources/
COPY --from=builder /m-azure-basic-infrastructure /workdir

ARG ARG_HOST_UID=1000
ARG ARG_HOST_GID=1000
RUN chown -R $ARG_HOST_UID:$ARG_HOST_GID \
    /workdir \
    /resources

USER $ARG_HOST_UID:$ARG_HOST_GID
