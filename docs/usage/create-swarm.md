# Create a Swarm

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
