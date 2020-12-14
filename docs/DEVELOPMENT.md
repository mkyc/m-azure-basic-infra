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
 