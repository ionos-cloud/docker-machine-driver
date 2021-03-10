# Introduction

## Overview

Rancher Driver is the official driver for Rancher Machine to use with Ionos Cloud. 

[Rancher Machine](https://github.com/rancher/machine) it is a fork of [Docker Machine](https://github.com/docker/machine) and it lets you create Docker hosts on your computer, on cloud providers and inside your own data center. 
For more information about Rancher Machine, you can visit the [GitHub Repository](https://github.com/rancher/machine) and the [Official Documentation](https://rancher.com/).
                                                                 
## Getting started

### Prerequisites
 
#### Installing Rancher Machine

This Ionos Cloud plugin will only work with Rancher Machine. Before we continue, you will need to install [Rancher Machine](https://github.com/rancher/machine/releases/).

#### Installing Go

The Ionos Cloud Rancher Driver is written in the Go programming language. Your system will need to have Go installed. Please refer to the [Go Install Documentation](https://golang.org/doc/install) if you do not have Go installed and configured for your system.

Remember to set `$GOPATH` and update `$PATH`. The following are just examples using the `export` command, you will need to adjust the paths for your particular installation.

```
export GOPATH=/usr/local/go
export PATH=$PATH:/usr/local/go/bin
```

### Installing

#### Local Version 

With the prerequisites taken care of, will need to run the following commands to install the Ionos Cloud Rancher Machine Driver:

```
git clone https://github.com/ionos-cloud/rancher-driver.git
```

After cloning the repository, you can build and install the driver itself:

```
cd $DIRECTORY_PATH/rancher-driver
make install
```

When successful, we will end up with a newly created `docker-machine-driver-ionoscloud` binary in `rancher-driver/bin/` and in `$GOPATH/bin/`. 

Depending how your `$PATH` is being set, you may need to copy the binary to `$PATH` in order to use the Rancher Driver. 

Note that the development version is a work-in-progress of a future stable release and can include bugs. Officially released versions will generally be more stable. Check the latest releases in the [Release Page](https://github.com/ionos-cloud/rancher-driver/releases).

### Usage

Before you start creating a Rancher Machine with the Ionos Cloud Rancher Driver, you need to authenticate in your Ionos Cloud account. Check the steps in the [Authentication](./usage/authentication.md) section.

For information about how to create a Rancher Machine with Ionos Cloud Rancher Driver, check the [Create a Machine](./usage/create-machine.md) section.

For information about how to create a Rancher Machine with Ionos Cloud Rancher Driver with [Swarm Mode](https://docs.docker.com/engine/swarm/), check the [Create a Swarm](./usage/create-swarm.md) section.

In order to see the available options and flags, check the [Options](./usage/options.md) section.

For more information about Rancher Machine commands on how to manage a machine, including examples, check the [Rancher Machine Commands](./usage/rancher-commands.md) section. 

## Feature Reference 

The IONOS Cloud Rancher Driver aims to offer access to all resources in the IONOS Cloud API and also offers some additional features that make the integration easier: 
- authentication for API calls
- handling of asynchronous requests 

## FAQ
- How can I open a bug/feature request?

Bugs & feature requests can be open on the repository issues: https://github.com/ionos-cloud/rancher-driver/issues/new/choose

- Can I contribute to the Rancher Driver?

Sure! Our repository is public, feel free to fork it and file a PR for one of the issues opened in the issues list. We will review it and work together to get it released.
