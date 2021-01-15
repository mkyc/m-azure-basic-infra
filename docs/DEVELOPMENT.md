# Development

This document describes requirements to be able to develop abd build module. 

## Operating systems

We develop this module in Linux/MacOs environments mostly. 

If you need to develop from Windows you can use the included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the examples the same way but then from then `examples/basic_flow_devcontainer` folder.

To build module container Linux and MacOS should work fine. 

On Windows [WSL(2)](https://docs.microsoft.com/en-us/windows/wsl/install-win10) can be used. In WSL2 you can easily share the docker host between Windows and the WSL2 environment by selecting the WSL2 backend in Docker Desktop settings. In WSL1 you can achieve the same with [this](https://nickjanetakis.com/blog/setting-up-docker-for-windows-and-wsl-to-work-flawlessly) tutorial.


# Build image

In the main directory, run:

```shell
make build
```

:warning: This command uses the default VERSION variable (default: `dev`).

or directly using Docker:

```shell
docker build --tag epiphanyplatform/azbi:dev .
```

:warning: Re-run the above commands will overwrite your existing docker image (if exists). To bypass that, specify a different `--tag` parameter. 

# Release module

```shell
make release
```

or if you want to set different version number:

```shell
make release VERSION=number_of_your_choice
```

# Run tests

Tests are described in a separate [document](TESTS.md).

# Develop locally with e-structures repository

Assuming you have [e-structures](https://github.com/epiphany-platform/e-structures) repository downloaded locally to directory `../../epiphany-platform/e-structures/` relatively to this repository directory, you should introduce following changes to develop locally with changes introduced also locally in e-structures repository: 

* Add copying and removing instructions to Makefile build target: 

    Change: 
    
    ```
    build: guard-VERSION guard-IMAGE guard-USER
        docker build \
            --build-arg ARG_M_VERSION=$(VERSION) \
            --build-arg ARG_HOST_UID=$(HOST_UID) \
            --build-arg ARG_HOST_GID=$(HOST_GID) \
            -t $(IMAGE_NAME) \
            .
    ``` 
    
    into: 
    
    ```
    build: guard-IMAGE_NAME
        rm -rf ./tmp
        mkdir -p ./tmp
        cp -R ../../epiphany-platform/e-structures/ ./tmp
        docker build \
            --build-arg ARG_M_VERSION=$(VERSION) \
            --build-arg ARG_HOST_UID=$(HOST_UID) \
            --build-arg ARG_HOST_GID=$(HOST_GID) \
            -t $(IMAGE_NAME) \
            .
        rm -rf ./tmp
    ```

* Add new line to Dockerfile

    Change: 
    
    ```
    ...
    COPY . $GOPATH/src/$GO_MODULE_NAME
    ...
    ```
    
    into: 
    
    ```
    ...
    COPY . $GOPATH/src/$GO_MODULE_NAME
    COPY tmp $GOPATH/src/github.com/epiphany-platform/e-structures
    ...
    ```

* Add the new section to go.mod file

    The new section: 
    
    ```
    replace (
        github.com/epiphany-platform/e-structures => ../../epiphany-platform/e-structures
    )
    ```
