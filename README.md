# Introduction

![CI](https://github.com/ionos-cloud/docker-machine-driver/workflows/CI/badge.svg)
[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=docker-machine-driver&metric=alert_status)](https://sonarcloud.io/dashboard?id=docker-machine-driver)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=docker-machine-driver&metric=bugs)](https://sonarcloud.io/dashboard?id=docker-machine-driver)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=docker-machine-driver&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=docker-machine-driver)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=docker-machine-driver&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=docker-machine-driver)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=docker-machine-driver&metric=security_rating)](https://sonarcloud.io/dashboard?id=docker-machine-driver)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=docker-machine-driver&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=docker-machine-driver)
[![Release](https://img.shields.io/github/v/release/ionos-cloud/docker-machine-driver.svg)](https://github.com/ionos-cloud/docker-machine-driver/releases/latest)
[![Release Date](https://img.shields.io/github/release-date/ionos-cloud/docker-machine-driver.svg)](https://github.com/ionos-cloud/docker-machine-driver/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/ionos-cloud/docker-machine-driver.svg)](https://github.com/ionos-cloud/docker-machine-driver)

![Alt text](.github/IONOS.CLOUD.BLU.svg?raw=true "Title")

> This library adds the support for creating Docker Machines hosted on IONOS Cloud.

## Overview

Docker Machine Driver is the official driver for Docker Machine to use with IONOS Cloud. It adds support for creating Docker Machines hosted on the IONOS Cloud. 

[Docker Machine](https://github.com/docker/machine) lets you create Docker hosts on your computer and inside your own data center. It creates servers, installs Docker on them, then configures the Docker client to talk to them. For more information about Docker Machine, check the official [GitHub Repository](https://github.com/docker/machine).

## Getting started

## Setup

### Option 1: Use with Rancher docker image

* Run the Rancher docker image on a publicly accessible server, reachable on ports `80` and `443`: `docker run -d --name=rancher-server --restart=unless-stopped -p 80:80 -p 443:443 --privileged rancher/rancher:v2.6.8`
* Add the Node Driver. Usage of the UI driver is **highly recommended**, but **not necessary**: `[server-ip]/dashboard/c/local/manager/pages/rke-drivers`
```markdown
  * Download URL: https://github.com/ionos-cloud/docker-machine-driver/releases/download/v<driver-version(v6.1.0-rc.1)>/docker-machine-driver-<driver-version(v6.1.0-rc.1)>-linux-amd64.tar.gz
  * Custom UI URL:  https://cdn.jsdelivr.net/gh/ionos-cloud/ui-driver-ionoscloud@<ui-driver-version(0.1.0)>/releases/<ui-driver-version(0.1.0)>/component.js
  * Whitelist Domains: https://cdn.jsdelivr.net
```
* The Docker Machine Driver for Ionoscloud is ready to use. Refer to the [Rancher Cluster](docs/rancher/rancher-cluster.md) section for version-specific instructions and further help with creating RKE1 templates and provisioning clusters.


### Option 2: Use with docker-machine or rancher-machine CLIs

Check the [Release Page](https://github.com/ionos-cloud/docker-machine-driver/releases) and find the corresponding archive for your operating system and architecture. You can download the archive from your browser or you can follow the next steps:

```text
# Check if /usr/local/bin is part of your PATH
echo $PATH

# Download and extract the binary (<version> is the full semantic version): 
curl -sL https://github.com/ionos-cloud/docker-machine-driver/releases/download/v<version>/docker-machine-driver-<version>-linux-amd64.tar.gz | tar -xzv

# Move the binary somewhere in your $PATH:
sudo mv ~/docker-machine-driver-ionoscloud /usr/local/bin

# See options for the driver to use with the Docker Machine CLI
docker-machine create --help --driver ionoscloud

# See options for the driver to use with the Rancher Machine CLI
rancher-machine create --help --driver ionoscloud
```

For Windows users, you can download the latest release available on [Release Page](https://github.com/ionos-cloud/docker-machine-driver/releases), unzip it and copy the binary in your `PATH`. You can follow this [official guide](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) that explains how to add tools to your `PATH`.

### Building From Source

#### Prerequisites
Please refer to the [Go Install Documentation](https://golang.org/doc/install) if you do not have Go installed and configured for your system.

Run the following commands to install the Ionos Cloud Docker Machine Driver:

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

Before you start using the Ionos Cloud Docker Machine Driver, you need to authenticate in your Ionos Cloud account. Check the steps in the [Authentication](docs/usage/authentication.md) section.

In order to see the available options and flags, check the [Options](docs/usage/options.md) section.

For more information about Docker/Rancher Machine commands on how to manage a machine, including examples, check the [Commands](docs/usage/commands.md) section.

### Docker Support

For information on how to create a Docker Machine with Ionos Cloud Docker Machine Driver, check the [Docker Machine](docs/docker/docker-machine.md) section.

For information on how to create a Docker Machine with Ionos Cloud Docker Machine Driver with [Swarm Mode](https://docs.docker.com/engine/swarm/), check the [Docker Swarm](docs/docker/docker-swarm.md) section.

For more details about possible issues, check the [Troubleshooting](docs/docker/troubleshooting.md) section.

### Rancher Support

For information on how to create a Rancher Machine with Ionos Cloud Docker Machine Driver, check the [Rancher Machine](docs/rancher/rancher-machine.md) section.

For information on how to create a Rancher Cluster via Rancher UI, using Ionos Cloud Docker Machine Driver, check the [Rancher Cluster](docs/rancher/rancher-cluster.md) section.

## Feature Reference

The IONOS Cloud Docker Machine Driver aims to offer access to all resources in the IONOS Cloud API and also offers some additional features that make the integration easier:

* authentication for API calls
* handling of asynchronous requests

## Contributing

Bugs & feature requests can be open on the repository issues: [https://github.com/ionos-cloud/docker-machine-driver/issues/new/choose](https://github.com/ionos-cloud/docker-machine-driver/issues/new/choose)

### Can I contribute to the Docker Machine Driver?

Sure! Our repository is public, feel free to fork it and file a PR for one of the issues opened in the issues list. We will review it and work together to get it released.

## License

[Apache 2.0](LICENSE)
