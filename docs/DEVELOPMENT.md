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
# Develop terraform scripts

There is a simple way to develop terraform scripts independently of module.

1) have image built already
1) go to `resources` directory in terminal
1) run module init command ie.: `docker run --rm -v $(pwd)/shared:/shared -t epiphanyplatform/azbi:dev init`
1) run module plan command with debug switch ie.: `docker run --rm -v $(pwd)/shared:/shared -e SUBSCRIPTION_ID=xxx -e CLIENT_ID=yyy -e CLIENT_SECRET=zzz -e TENANT_ID=vvv epiphanyplatform/azbi:dev plan --debug`
1) that would provide you detailed output and one of first lines provides used .tfvars.json file. It will look similar to following: 
   ```
   ...
   2021-01-15T13:42:27Z DBG go/src/github.com/epiphany-platform/m-azure-basic-infrastructure/cmd/common.go:61 > templateTfVars
   2021-01-15T13:42:27Z INF go/src/github.com/epiphany-platform/m-azure-basic-infrastructure/cmd/common.go:68 > {"name":"epiphany","location":"northeurope","address_space":["10.0.0.0/16"],"subnets":[{"name":"main","address_prefixes":["10.0.1.0/24"]}],"vm_groups":[{"name":"vm-group0","vm_count":1,"vm_size":"Standard_DS2_v2","use_public_ip":true,"subnet_names":["main"],"vm_image":{"publisher":"Canonical","offer":"UbuntuServer","sku":"18.04-LTS","version":"18.04.202006101"}}],"rsa_pub_path":"/shared/vms_rsa.pub"}
   ...
   ```
1) you should copy provided JSON to file `./terraform/terraform.tfvars.json` relatively to `resources` directory
1) run `mkdir shared` also in that `resources` directory
1) run `ssh-keygen -t rsa -b 4096 -f ./shared/vms_rsa -N ''` to generate required key pair
1) now you can run terraform directly from `resources` directory

To be able to run version of terraform that scripts require you can use following docker snippet: 

```
docker run --rm -it -e ARM_SUBSCRIPTION_ID=xxx -e ARM_CLIENT_ID=yyy -e ARM_CLIENT_SECRET=zzz -e ARM_TENANT_ID=vvv -v $(pwd)/terraform:/workspace -v $(pwd)/shared:/shared -w /workspace hashicorp/terraform:0.13.2 apply
```
(notice `ARM_` prefix on passed envs)
