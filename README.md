# Ionos Cloud Rancher Driver

![CI](https://github.com/ionos-cloud/rancher-driver/workflows/CI/badge.svg)

## Overview

This is the official Rancher Driver for use with the Ionos Cloud. Rancher Machine it is a fork of [Docker Machine](https://github.com/docker/machine) and it lets you create Docker hosts on your computer, on cloud providers and inside your own data center. 
For more information about Rancher Machine, you can visit [the Github Repository](https://github.com/rancher/machine) and [the Official Documentation](https://rancher.com/).
                                                                 
## Getting started

### Install Rancher Machine

This Ionos Cloud plugin will only work with Rancher Machine. Before we continue, you will need to install [Rancher Machine](https://github.com/rancher/machine/releases/).

### Install Go

The Ionos Cloud Rancher Driver is written in the Go programming language. Your system will need to have Go installed. Please refer to the [Go Install Documentation](https://golang.org/doc/install) if you do not have Go installed and configured for your system.

Remember to set `$GOPATH` and update `$PATH`. The following are just examples using the `export` command, you will need to adjust the paths for your particular installation.

```
export GOPATH=/usr/local/go
export PATH=$PATH:/usr/local/go/bin
```

### Install Rancher Driver

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

Note that the development version is a work-in-progress of a future stable release and can include bugs. Officially released versions will generally be more stable. Check the latest releases in [the Release Page](https://github.com/ionos-cloud/rancher-driver/releases).

## How to Use

### Authenticate with Ionos Cloud

Before you create a Ionos Cloud Rancher Machine you will need to set two environment variables containing your Ionos Cloud credentials. These would be the same username and password that you use to log into the Ionos Cloud DCD (Data Center Designer):

```
export IONOSCLOUD_USERNAME="ionoscloud_username"
export IONOSCLOUD_PASSWORD="ionoscloud_password"
```

It is possible to pass your credentials on the command-line using `--ionoscloud-username` and `--ionoscloud-password` if you prefer.

### Create a Machine

Now run `rancher-machine create` with the relevant parameters. This example will use mostly default values and will therefore be created in the `us/las` location.

```
rancher-machine create --driver ionoscloud test-machine
```

It should produce results similar to this:

```
Running pre-create checks...
Creating machine...
(test-machine) Creating SSH key...
(test-machine) DataCenter Created
(test-machine) LAN Created
(test-machine) Server Created
(test-machine) Volume Attached to Server
(test-machine) NIC Attached to Server
(test-machine) 158.222.102.154
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with ubuntu(systemd)...
Installing Docker from: https://get.docker.com
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: rancher-machine env test-machine
```

### Available Options

To get detailed information about the possible options to use in order to create a Rancher Machine with Ionos Cloud Driver, run the command:

```
rancher-machine create --help --driver ionoscloud
```

Options with their default values:

```
   --driver, -d "virtualbox"										Driver to create machine with. [$MACHINE_DRIVER]
   --engine-env [--engine-env option --engine-env option]						Specify environment variables to set in the engine
   --engine-insecure-registry [--engine-insecure-registry option --engine-insecure-registry option]	Specify insecure registries to allow with the created engine
   --engine-install-url "https://get.docker.com"							Custom URL to use for engine installation [$MACHINE_DOCKER_INSTALL_URL]
   --engine-label [--engine-label option --engine-label option]						Specify labels for the created engine
   --engine-opt [--engine-opt option --engine-opt option]						Specify arbitrary flags to include with the created engine in the form flag=value
   --engine-registry-mirror [--engine-registry-mirror option --engine-registry-mirror option]		Specify registry mirrors to use [$ENGINE_REGISTRY_MIRROR]
   --engine-storage-driver 										Specify a storage driver to use with the engine
   --ionoscloud-datacenter-id 										Ionos Cloud Virtual Data Center Id
   --ionoscloud-disk-size "50"										Ionos Cloud Volume Disk-Size (10, 50, 100, 200, 400) [$IONOSCLOUD_DISK_SIZE]
   --ionoscloud-disk-type "HDD"										Ionos Cloud Volume Disk-Type (HDD, SSD) [$IONOSCLOUD_DISK_TYPE]
   --ionoscloud-endpoint "https://api.ionos.com/cloudapi/v5"						Ionos Cloud API Endpoint [$IONOSCLOUD_ENDPOINT]
   --ionoscloud-image "ubuntu:latest"									Ionos Cloud Image Alias [$IONOSCLOUD_IMAGE]
   --ionoscloud-location "us/las"									Ionos Cloud Location [$IONOSCLOUD_LOCATION]
   --ionoscloud-password 										Ionos Cloud Password [$IONOSCLOUD_PASSWORD]
   --ionoscloud-server-availability-zone "AUTO"								Ionos Cloud Server Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)
   --ionoscloud-server-cores "4"									Ionos Cloud Server Cores (2, 3, 4, 5, 6, etc.) [$IONOSCLOUD_SERVER_CORES]
   --ionoscloud-server-cpu-family "AMD_OPTERON"								Ionos Cloud Server CPU families (AMD_OPTERON,INTEL_XEON) [$IONOSCLOUD_SERVER_CPU_FAMILY]
   --ionoscloud-server-ram "2048"									Ionos Cloud Server Ram (1024, 2048, 3072, 4096, etc.) [$IONOSCLOUD_SERVER_RAM]
   --ionoscloud-username 										Ionos Cloud Username [$IONOSCLOUD_USERNAME]
   --ionoscloud-volume-availability-zone "AUTO"								Ionos Cloud Volume Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)
   --swarm												Configure Machine to join a Swarm cluster
   --swarm-addr 											addr to advertise for Swarm (default: detect and use the machine IP)
   --swarm-discovery 											Discovery service to use with Swarm
   --swarm-experimental											Enable Swarm experimental features
   --swarm-host "tcp://0.0.0.0:3376"									ip/socket to listen on for Swarm master
   --swarm-image "swarm:latest"										Specify Docker image to use for Swarm [$MACHINE_SWARM_IMAGE]
   --swarm-join-opt [--swarm-join-opt option --swarm-join-opt option]					Define arbitrary flags for Swarm join
   --swarm-master											Configure Machine to be a Swarm master
   --swarm-opt [--swarm-opt option --swarm-opt option]							Define arbitrary flags for Swarm master
   --swarm-strategy "spread"										Define a default scheduling strategy for Swarm
   --tls-san [--tls-san option --tls-san option]							Support extra SANs for TLS certs
```

###  Set up environment

To configure your shell run:

```
eval $(rancher-machine env test-machine)
```

### Inspect a Machine

To see more information about the machine created, use:

```
rancher-machine inspect test-machine
```

### List Machines

To list the machines you have created, use the command:

```
rancher-machine ls
```

It will return information about your machines, similar to this:

```
NAME           ACTIVE   DRIVER         STATE     URL                         SWARM   DOCKER    ERRORS
test-machine   *        ionoscloud     Running   tcp://162.254.26.156:2376           v20.10.5
```

### Get Status

To get the status of a Rancher Machine created, run: 

```
rancher-machine status test-machine
```

### Start a Machine

To start a Rancher Machine, run: 

```
rancher-machine start test-machine
```

### Stop a Machine

To stop a Rancher Machine, run: 

```
rancher-machine stop test-machine
```

### Restart a Machine

To restart a Rancher Machine, run: 

```
rancher-machine restart test-machine
```

### Remove a Machine

To remove a Rancher Machine and all the resources associated with it, run: 

```
rancher-machine rm test-machine
```

It should produce results similar to this:

```
About to remove test-machine
WARNING: This action will delete both local reference and remote instance.
Are you sure? (y/n): y
(test-machine) Starting deleting resources...
(test-machine) NIC Deleted
(test-machine) Volume Deleted
(test-machine) Server Deleted
(test-machine) LAN Deleted
(test-machine) DataCenter Deleted
(test-machine) IPBlock Deleted
Successfully removed test-machine
```

The remove command can also be used with `--force` or `-f` flag. 

For more available commands and options to manage a Rancher Machine, use:

```
rancher-machine help
```

### Create a Swarm

You can use Rancher Machine to provision Swarm clusters. 

Before you create a swarm of Ionos Cloud machines, run this command:

```
docker swarm init
```

Then use the output `${DOCKER_SWARM_TOKEN}` to create the swarm and set a swarm master:

```
rancher-machine create --driver ionoscloud --swarm --swarm-master --swarm-discovery token://${DOCKER_SWARM_TOKEN} swarm-master-test
```

It should produce results similar to this:

```
Running pre-create checks...
Creating machine...
(swarm-master-test) Creating SSH key...
(swarm-master-test) DataCenter Created
(swarm-master-test) LAN Created
(swarm-master-test) Server Created
(swarm-master-test) Volume Attached to Server
(swarm-master-test) NIC Attached to Server
(swarm-master-test) 158.222.102.158
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with ubuntu(systemd)...
Installing Docker from: https://get.docker.com
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Configuring swarm...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: rancher-machine env swarm-master-test
```

To create a swarm child, use the command:

```
rancher-machine create -d ionoscloud --swarm --swarm-discovery token://${DOCKER_SWARM_TOKEN} swarm-child-test
```

When running `rancher-machine ls`, it should produce results similar to this:

```
NAME                ACTIVE   DRIVER       STATE     URL                          SWARM                        DOCKER     ERRORS
swarm-child-test    -        ionoscloud   Running   tcp://158.222.102.154:2376   swarm-master-test            v20.10.5   
swarm-master-test   *        ionoscloud   Running   tcp://158.222.102.158:2376   swarm-master-test (master)   v20.10.5   
```

## Testing

To run unit tests, use:

```
make test
```

## Uninstall 

### Local Version

To uninstall a local version built with the steps from [Install Rancher Driver](#install-rancher-driver), use:

```
make clean
```

## Feature Reference 

The IONOS Cloud Rancher Driver aims to offer access to all resources in the IONOS Cloud API and also offers some additional features that make the integration easier: 
- authentication for API calls
- handling of asynchronous requests 

## FAQ
- How can I open a bug/feature request?

Bugs & feature requests can be open on the repository issues: https://github.com/ionos-cloud/rancher-driver/issues/new/choose
