# Docker Machine Driver

## Overview

Docker Machine Driver is the official driver for Docker Machine to use with IONOS Cloud. It adds support for creating Docker Machines hosted on the IONOS Cloud.

[Docker Machine](https://github.com/docker/machine) lets you create Docker hosts on your computer and inside your own data center. It creates servers, installs Docker on them, then configures the Docker client to talk to them. For more information about Docker Machine, check the official [GitHub Repository](https://github.com/docker/machine).

## Getting started

### Prerequisites

#### Installing Docker Machine

This Ionos Cloud plugin works with Docker Machine and with Rancher Machine as well. Before we continue, you will need to install [Docker Machine](https://docs.docker.com/machine/install-machine/) or [Rancher Machine](https://github.com/rancher/machine/releases/).

#### Installing Go

The Ionos Cloud Docker Machine Driver is written in the Go programming language. Your system will need to have Go installed. Please refer to the [Go Install Documentation](https://golang.org/doc/install) if you do not have Go installed and configured for your system.

Remember to set `$GOPATH` and update `$PATH`. The following are just examples using the `export` command, you will need to adjust the paths for your particular installation.

```text
export GOPATH=/usr/local/go
export PATH=$PATH:/usr/local/go/bin
```

### Installing

#### Released Binaries

Check the [Release Page](https://github.com/ionos-cloud/docker-machine-driver/releases) and find the corresponding archive for your operating system and architecture. You can download the archive from your browser or you can follow the next steps:

```text
# Check if /usr/local/bin is part of your PATH
echo $PATH

# Download and extract the binary (<version> is the full semantic version): 
curl -sL https://github.com/ionos-cloud/docker-machine-driver/releases/download/v<version>/docker-machine-driver-<version>-linux-amd64.tar.gz | tar -xzv

# Move the binary somewhere in your $PATH:
sudo mv ~/docker-machine-driver-ionoscloud /usr/local/bin

# See options for the driver to use with the Docker Machine
docker-machine create --help --driver ionoscloud
```

For Windows users, you can download the latest release available on [Release Page](https://github.com/ionos-cloud/docker-machine-driver/releases), unzip it and copy the binary in your `PATH`. You can follow this \[official guide\]\([https://msdn.microsoft.com/en-us/library/office/ee537574\(v=office.14\).aspx](https://msdn.microsoft.com/en-us/library/office/ee537574%28v=office.14%29.aspx)\) that explains how to add tools to your `PATH`.

#### Local Version

With the prerequisites taken care of, will need to run the following commands to install the Ionos Cloud Docker Machine Driver:

```text
git clone https://github.com/ionos-cloud/docker-machine-driver.git
```

After cloning the repository, you can build and install the driver itself:

```text
cd $DIRECTORY_PATH/docker-machine-driver
make install
```

When successful, we will end up with a newly created `docker-machine-driver-ionoscloud` binary in `docker-machine-driver/bin/` and in `$GOPATH/bin/`.

Depending on how your `$PATH` is being set, you may need to copy the binary to `$PATH` in order to use the Docker Machine Driver.

```text
sudo cp $DIRECTORY_PATH/docker-machine-driver/bin/docker-machine-driver-ionoscloud /usr/local/bin/docker-machine-driver-ionoscloud
```

Note that the development version is a work-in-progress of a future stable release and can include bugs. Officially released versions will generally be more stable. Check the latest releases in the [Release Page](https://github.com/ionos-cloud/docker-machine-driver/releases).

### Usage

Before you start using the Ionos Cloud Docker Machine Driver, you need to authenticate in your Ionos Cloud account. Check the steps in the [Authentication](usage/authentication.md) section.

In order to see the available options and flags, check the [Options](usage/options.md) section.

For more information about Docker/Rancher Machine commands on how to manage a machine, including examples, check the [Commands](usage/commands.md) section.

### Docker Support

For information on how to create a Docker Machine with Ionos Cloud Docker Machine Driver, check the [Docker Machine](docker/docker-machine.md) section.

For information on how to create a Docker Machine with Ionos Cloud Docker Machine Driver with [Swarm Mode](https://docs.docker.com/engine/swarm/), check the [Docker Swarm](docker/docker-swarm.md) section.

For more details about possible issues, check the [Troubleshooting](docker/troubleshooting.md) section.

### Rancher Support

For information on how to create a Rancher Machine with Ionos Cloud Docker Machine Driver, check the [Rancher Machine](https://github.com/ionos-cloud/docker-machine-driver/tree/634b60ff47c8a3294d5955ea0eb19bd3c18ac454/docs/rancher/rancher-machine.md) section.

For information on how to create a Rancher Cluster via Rancher UI, using Ionos Cloud Docker Machine Driver, check the [Rancher Cluster]() section.

## Feature Reference

The IONOS Cloud Docker Machine Driver aims to offer access to all resources in the IONOS Cloud API and also offers some additional features that make the integration easier:

* authentication for API calls
* handling of asynchronous requests 

## FAQ

* How can I open a bug/feature request?

Bugs & feature requests can be open on the repository issues: [https://github.com/ionos-cloud/docker-machine-driver/issues/new/choose](https://github.com/ionos-cloud/docker-machine-driver/issues/new/choose)

* Can I contribute to the Docker Machine Driver?

Sure! Our repository is public, feel free to fork it and file a PR for one of the issues opened in the issues list. We will review it and work together to get it released.

