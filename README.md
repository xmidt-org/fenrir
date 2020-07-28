# fenrir
(pronounced “fen-rear”)

[![Build Status](https://travis-ci.com/xmidt-org/fenrir.svg?branch=main)](https://travis-ci.com/xmidt-org/fenrir)
[![codecov.io](http://codecov.io/github/xmidt-org/fenrir/coverage.svg?branch=main)](http://codecov.io/github/xmidt-org/fenrir?branch=main)
[![Code Climate](https://codeclimate.com/github/xmidt-org/fenrir/badges/gpa.svg)](https://codeclimate.com/github/xmidt-org/fenrir)
[![Issue Count](https://codeclimate.com/github/xmidt-org/fenrir/badges/issue_count.svg)](https://codeclimate.com/github/xmidt-org/fenrir)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/fenrir)](https://goreportcard.com/report/github.com/xmidt-org/fenrir)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/fenrir/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/xmidt-org/fenrir.svg)](CHANGELOG.md)

The service that prunes expired entries from the database.

## Summary

Fenrir prunes expired events from the database. For more information on how 
Fenrir fits into codex, check out [the codex README](https://github.com/xmidt-org/codex-deploy).

This project isn't currently being worked on, as yugabyte doesn't need this service.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Details](#details)
- [Build](#build)
- [Deploy](#deploy)
- [Contributing](#contributing)

## Code of Conduct

This project and everyone participating in it are governed by the [XMiDT Code Of Conduct](https://xmidt.io/code_of_conduct/). 
By participating, you agree to this Code.

## Details

Fenrir makes a call to the database at a configurable interval.  The database 
call tells the database to delete records whose deathdate has passed, but there 
is a configurable limit to the number of records to delete, to avoid timeout 
issues.

## Build

### Source

In order to build from the source, you need a working Go environment with 
version 1.11 or greater. Find more information on the [Go website](https://golang.org/doc/install).

You can directly use `go get` to put the Fenrir binary into your `GOPATH`:
```bash
GO111MODULE=on go get github.com/xmidt-org/fenrir
```

You can also clone the repository yourself and build using make:

```bash
mkdir -p $GOPATH/src/github.com/xmidt-org
cd $GOPATH/src/github.com/xmidt-org
git clone git@github.com:Comcast/fenrir.git
cd fenrir
make build
```

### Makefile

The Makefile has the following options you may find helpful:
* `make build`: builds the Fenrir binary
* `make rpm`: builds an rpm containing Fenrir
* `make docker`: builds a docker image for Fenrir, making sure to get all 
   dependencies
* `make local-docker`: builds a docker image for Fenrir with the assumption
   that the dependencies can be found already
* `make it`: runs `make docker`, then deploys Fenrir and a cockroachdb 
   database into docker.
* `make test`: runs unit tests with coverage for Fenrir
* `make clean`: deletes previously-built binaries and object files

### Docker

The docker image can be built either with the Makefile or by running a docker 
command.  Either option requires first getting the source code.

See [Makefile](#Makefile) on specifics of how to build the image that way.

For running a command, either you can run `docker build` after getting all 
dependencies, or make the command fetch the dependencies.  If you don't want to 
get the dependencies, run the following command:
```bash
docker build -t fenrir:local -f deploy/Dockerfile .
```
If you want to get the dependencies then build, run the following commands:
```bash
GO111MODULE=on go mod vendor
docker build -t fenrir:local -f deploy/Dockerfile.local .
```

For either command, if you want the tag to be a version instead of `local`, 
then replace `local` in the `docker build` command.

### Kubernetes

WIP. TODO: add info

## Deploy

For deploying on Docker or in Kubernetes, refer to the [deploy README](https://github.com/xmidt-org/codex-deploy/tree/main/deploy/README.md).

For running locally, ensure you have the binary [built](#Source).  If it's in 
your `GOPATH`, run:
```
fenrir
```
If the binary is in your current folder, run:
```
./fenrir
```

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md).
