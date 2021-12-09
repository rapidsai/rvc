# rvc

`rvc` (_RAPIDS version converter_) is a tool built with Golang that automatically
convert a `RAPIDS` version (CalVer) to a `ucx-py` version (SemVer) and vice versa.

`rvc` is published in two different ways:
  - As an AWS Lambda instance
  - As a CLI binary

## AWS Lambda instance

`rvc` is deployed as an AWS Lambda instance using the _Serverless_ framework.
The deployment configuration can be seen in the
[serverless.yaml](https://github.com/rapidsai/rvc/blob/main/serverless.yaml) file.
A deployment will happen automatically anytime a change is merged to the main branch.
See the [deploy.yaml](https://github.com/rapidsai/rvc/blob/main/.github/workflows/deploy.yaml)
GitHub Action for more details.

`rvc` is deployed at this endpoint: https://version.gpuci.io

Two different routes are exposed:
  - https://version.gpuci.io/ucx-py/{version}: Converts a `RAPIDS` version to a `ucx-py` version
  - https://version.gpuci.io/rapids/{version}: Converts a `ucx-py` version to a `RAPIDS` version

## CLI binary

`rvc` is also available as a CLI binary.

### Installation

Download the latest binary from GitHub and add it to one of your `PATH` directories:

```
wget https://github.com/rapidsai/rvc/releases/latest/download/rvc -O rvc
chmod +x rvc
sudo mv ./rvc /usr/local/bin
```

### Usage

```
Usage of rvc:
  -rapids string
        Rapids version
  -ucx-py string
        ucx-py version
```

`-rapids` and `ucx-py` options are mutually exclusive.

Examples:
```
$ rvc -rapids 21.12
0.23
```
```
$ rvc -ucx-py 0.22
21.10
```
