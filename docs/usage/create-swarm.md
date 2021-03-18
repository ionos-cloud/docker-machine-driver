# Create Swarm

You can use Docker Machine to provision Swarm clusters. 

Before you create a swarm of Ionos Cloud machines, run this command:

```text
docker swarm init
```

Then use the output `${DOCKER_SWARM_TOKEN}` to create the swarm and set a swarm master:

```text
docker-machine create --driver ionoscloud --swarm --swarm-master --swarm-discovery token://${DOCKER_SWARM_TOKEN} swarm-master-test
```

It should produce results similar to this:

```text
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
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Configuring swarm...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env swarm-master-test
```

To create a swarm child, use the command:

```text
docker-machine create -d ionoscloud --swarm --swarm-discovery token://${DOCKER_SWARM_TOKEN} swarm-child-test
```

When running `docker-machine ls`, it should produce results similar to this:

```text
NAME                ACTIVE   DRIVER       STATE     URL                          SWARM                        DOCKER     ERRORS
swarm-child-test    -        ionoscloud   Running   tcp://158.222.102.154:2376   swarm-master-test            v20.10.5   
swarm-master-test   *        ionoscloud   Running   tcp://158.222.102.158:2376   swarm-master-test (master)   v20.10.5   
```

For more details about possible issues, check the [Troubleshooting](troubleshooting.md) guide.

