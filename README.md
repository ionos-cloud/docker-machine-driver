# Ionos Cloud Rancher Driver

This is the official Rancher driver for use with the Ionos Cloud.

## Table of Contents
* [Install Docker Machine](#install-docker-machine)
* [Install Go](#install-go)
* [Install Driver](#install-driver)
* [Create a Machine](#create-a-machine)
* [Create a Swarm](#create-a-swarm)
* [Support](#support)

## Install Rancher

This ProfitBricks plugin will only work with Docker Machine. Before we continue, you will need to install [Docker Machine](https://docs.docker.com/machine/install-machine/). Docker Machine is included as part of Docker Toolbox. You can gain access to `docker-machine` by installing Docker Toolbox on Mac OS X or Windows. It is also possible to install just Docker Machine without the rest of the components of Docker Toolbox.

## Install Go

The ProfitBricks Docker Machine Driver is written in the Go programming language. Your system will need to have Go installed. Please refer to the [Go Install Documentation](https://golang.org/doc/install) if you do not have Go installed and configured for your system.

Remember to set `$GOPATH` and update `$PATH`. The following are just examples using the `export` command, you will need to adjust the paths for your particular installation.

    export GOPATH=/usr/local/go
    export PATH=$PATH:/usr/local/go/bin

## Install Driver

With those prerequisites taken care of, will need to run the following commands to install the ProfitBricks Docker Machine driver:

    go get github.com/profitbricks/docker-machine-driver-profitbricks

If you just installed Go, you may get an error indicating the need to configure the `$GOPATH` environment variable. Once `$GOPATH` is set properly, the command should complete successfully.

Next we need to build and install the driver itself.

    cd $GOPATH/src/github.com/profitbricks/docker-machine-driver-profitbricks
    make install

When successful, we will end up with a newly created `docker-machine-driver-profitbricks` binary in `$GOPATH/bin/`.

## Create a Machine

Before you create a ProfitBricks machine you will need to set two environment variables containing your ProfitBricks credentials. These would be the same username and password that you use to log into the ProfitBricks DCD (Data Center Designer):

    export PROFITBRICKS_USERNAME="profitbricks_username"
    export PROFITBRICKS_PASSWORD="profitbricks_password"

It is possible to pass your credentials on the command-line using `--profitbricks-username` and `--profitbricks-password` if you prefer.

Now run `docker-machine create` with the relevant parameters. This example will use mostly default values and will therefore be created in the `us/las` location.

    docker-machine create --driver profitbricks test-machine

It should produce results similar to this:

```
Running pre-create checks...
Creating machine...
(test-machine) Datacenter Created
(test-machine) Server Created
(test-machine) Volume Created
(test-machine) Attached a volume  to a server.
(test-machine) LAN Created
(test-machine) NIC created
(test-machine) Updated server's boot image
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with ubuntu(systemd)...
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env test-machine

```

---

To get detailed information about the possible options,  run the command:

`docker-machine create --help --driver profitbricks`

### Available Options



#### --driver, -d "profitbricks"

Driver to create machine with. [$MACHINE_DRIVER]

#### --engine-env [--engine-env option --engine-env option]

 Specify environment variables to set in the engine

#### --engine-insecure-registry [--engine-insecure-registry option --engine-insecure-registry option]

 Specify insecure registries to allow with the created engine
   --engine-install-url "https://get.docker.com"                                                        Custom URL to use for engine installation [$MACHINE_DOCKER_INSTALL_URL]

#### --engine-label [--engine-label option --engine-label 	option]    

Specify labels for the created engine

#### --engine-opt [--engine-opt option --engine-opt option]                                              

Specify arbitrary flags to include with the created engine in the form flag=value

#### --engine-registry-mirror [--engine-registry-mirror option --engine-registry-mirror option]          

Specify registry mirrors to use [$ENGINE_REGISTRY_MIRROR]

#### --engine-storage-driver                                                                             

Specify a storage driver to use with the engine                                    

#### --profitbricks-cores "4"

ProfitBricks cores (2, 3, 4, 5, 6, etc.) [$PROFITBRICKS_CORES]

#### --profitbricks-cpu-family "AMD_OPTERON"                                                                                                                                         

ProfitBricks CPU families (AMD_OPTERON,INTEL_XEON) [$PROFITBRICKS_CPU_FAMILY]

#### --profitbricks-datacenter-id                                                                        

ProfitBricks Virtual Data Center Id

#### --profitbricks-disk-size "50"                                                                       

ProfitBricks disk size (10, 50, 100, 200, 400) [$PROFITBRICKS_DISK_SIZE]

#### --profitbricks-disk-type "HDD"                                                                      

ProfitBricks disk type (HDD, SSD) [$PROFITBRICKS_DISK_TYPE]

#### --profitbricks-endpoint "https://api.profitbricks.com/cloudapi/v4"                                  

ProfitBricks API endpoint [$PROFITBRICKS_ENDPOINT]

#### --profitbricks-image "Ubuntu-16.04" 

ProfitBricks image [$PROFITBRICKS_IMAGE], you can use the image alias "Ubuntu:latest" or the image name "Ubuntu-16.04".                                                                  

#### --profitbricks-location "us/las"                                                                    

ProfitBricks location [$PROFITBRICKS_LOCATION]

#### --profitbricks-password                                                                             

profitbricks password [$PROFITBRICKS_PASSWORD]

#### --profitbricks-ram "2048"                                                                           

ProfitBricks ram (1024, 2048, 3072, 4096, etc.) [$PROFITBRICKS_RAM]

#### --profitbricks-server-availability-zone "AUTO"                                                      

ProfitBricks Server Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)

#### --profitbricks-username                                                                             

ProfitBricks username [$PROFITBRICKS_USERNAME]

#### --profitbricks-volume-availability-zone "AUTO"                                                      

ProfitBricks Volume Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)

#### --swarm                                                                                             

Configure Machine to join a Swarm cluster

#### --swarm-addr                                                                                        

addr to advertise for Swarm (default: detect and use the machine IP)

#### --swarm-discovery                                                                                   

Discovery service to use with Swarm

#### --swarm-experimental

Enable Swarm experimental features

#### --swarm-host "tcp://0.0.0.0:3376"

ip/socket to listen on for Swarm master                                                                                

#### --swarm-image "swarm:latest"                                                                        

Specify Docker image to use for Swarm [$MACHINE_SWARM_IMAGE]

#### --swarm-join-opt [--swarm-join-opt option --swarm-join-opt option]                                  

Define arbitrary flags for Swarm join

#### --swarm-master                                                                                      

Configure Machine to be a Swarm master

#### --swarm-opt [--swarm-opt option --swarm-opt option]                                                 

Define arbitrary flags for Swarm master

#### --swarm-strategy "spread"                                                                           

Define a default scheduling strategy for Swarm

#### --tls-san [--tls-san option --tls-san option]                                                       

Support extra SANs for TLS certs

---


To list the machines you have created, use the command:

    docker-machine ls

It will return information about your machines, similar to this:

```
NAME           ACTIVE   DRIVER         STATE     URL                         SWARM   DOCKER    ERRORS
default        -        virtualbox     Running   tcp://192.168.99.100:2376           v1.10.2
test-machine   -        profitbricks   Running   tcp://162.254.26.156:2376           v1.10.3

```

# Create a Swarm

Before you create a swarm of ProfitBricks machines, run this command:

    docker run --rm swarm create

Then use the output to create the swarm and set a swarm master:

    docker-machine create -d profitbricks --swarm --swarm-master --swarm-discovery token://f3a75db19a03589ac28550834457bfc3 swarm-master-test

To create a swarm child, use the command:

```docker-machine create -d profitbricks --swarm --swarm-discovery token://f3a75db19a03589ac28550834457bfc3 swarm-child-test```

## Support

You are welcome to contact us with questions or comments at [ProfitBricks DevOps Central](https://devops.profitbricks.com/). Please report any issues via [GitHub's issue tracker](https://github.com/profitbricks/docker-machine-driver-profitbricks/issues).