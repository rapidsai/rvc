# rvc

`rvc` (_RAPIDS version converter_) is a tool built with Golang that automatically
convert a `RAPIDS` version (CalVer) to a `ucxx` version (SemVer) and vice versa.

Note that in RAPIDS 25.10 the UCX-Py project was discontinued and archived.
However, for practical reasons we decided to maintain links, request names and
internal names, such as functions, still referring UCX-Py, such as in the URL
(see [AWS Lambda instance](#aws-lambda-instance)). This is beneficial from a
maintenance standpoint by avoiding the need to review all uses and rename them
all, which is not a strict requirement given these names are not user-facing.
Nevertheless, the current purpose for this repository is to support the
[UCXX](https://github.com/rapidsai/ucxx) project that superseded UCX-Py.

## Motivation

In June 2021, RAPIDS moved from a SemVer versioning to a CalVer versioning.
As `ucxx` is expected to be upstreamed to `ucx`, it is not possible to adopt
a CalVer versioning for it, as the versions would have been greater than the
current `ucx` version. `rvc` is designed to ease the conversion between both
versioning.

## Deliverables

`rvc` is published in two different ways:
  - As an AWS Lambda instance
  - As a CLI binary

### AWS Lambda instance

`rvc` is deployed as an AWS Lambda instance using the _Serverless_ framework.
The deployment configuration can be seen in the
[serverless.yaml](https://github.com/rapidsai/rvc/blob/main/serverless.yaml) file.
A deployment will happen automatically anytime a change is merged to the main branch.
See the [deploy.yaml](https://github.com/rapidsai/rvc/blob/main/.github/workflows/deploy.yaml)
GitHub Action for more details.

`rvc` is deployed at this endpoint: https://version.gpuci.io

Two different routes are exposed:
  - https://version.gpuci.io/ucx-py/{version}: Converts a `ucxx` version to a `RAPIDS` version
  - https://version.gpuci.io/rapids/{version}: Converts a `RAPIDS` version to a `ucxx` version

Examples:
```sh
$ RAPIDS_VER=$(curl -sL https://version.gpuci.io/ucx-py/0.46)
$ echo "${RAPIDS_VER}"
25.10
```
```sh
$ UCXX_VER=$(curl -sL https://version.gpuci.io/rapids/25.10)
$ echo "${UCX_VER}"
0.46
```

### CLI binary

`rvc` is also available as a CLI binary.

#### Installation

Download the latest binary from GitHub and add it to one of your `PATH` directories:

```sh
wget https://github.com/rapidsai/rvc/releases/latest/download/rvc -O rvc
chmod +x rvc
sudo mv ./rvc /usr/local/bin
```

#### Usage

```
Usage of rvc:
  -rapids string
        Rapids version
  -ucx-py string
        ucx-py version
```

`-rapids` and `ucx-py` options are mutually exclusive.

Examples:
```sh
$ rvc -rapids 25.10
0.46
```
```sh
$ rvc -ucx-py 0.56
25.10
```

## Contributing

### Add a new version mapping

If you need to add a new version mapping, you only need to update this
[map](https://github.com/rapidsai/rvc/blob/main/pkg/rvc/rvc.go#L15).

### Improving `rvc`

Requirements:
  - golang >= 1.17
  - serverless

A Makefile with some useful rules is provided:
  - `make build`: Run unit tests and build CLI and Serverless binaries
  - `make test`: Run unit tests
  - `make build_cli`: Build CLI binary
  - `make build_serverless`: Build Serverless binary
  - `make fmt`: Format code
  - `make coverage`: Compute and display tests coverage
