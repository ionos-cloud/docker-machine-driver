# Troubleshooting

Here is a quick troubleshooting guide which may help you to resolve issues you may be facing. We all know that sometimes things do not go according to the plan. 

## Unable to verify the Docker daemon is listening...

When running the following command:

```text
docker-machine create --driver ionoscloud test-machine
```

you may be getting the following results:

```text
Creating CA: /home/runner/.docker/machine/certs/ca.pem
Creating client certificate: /home/runner/.docker/machine/certs/cert.pem
Running pre-create checks...
Creating machine...
(test-machine) Creating SSH key...
(test-machine) DataCenter Created
(test-machine) LAN Created
(test-machine) Server Created
(test-machine) Volume Attached to Server
(test-machine) NIC Attached to Server
(test-machine) 158.222.102.181
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with ubuntu(systemd)...
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Error creating machine: Error running provisioning: Unable to verify the Docker daemon is listening: Maximum number of retries (10) exceeded
```

The resources on Ionos Cloud are created, but the machine is unable to verify the Docker daemon, even after a number of retries.

If you run `docker-machine ls`, the output will probably be:

```text
NAME           ACTIVE   DRIVER       STATE     URL                          SWARM   DOCKER    ERRORS
test-machine   -        ionoscloud   Running   tcp://158.222.102.183:2376           Unknown   Unable to query docker version: Cannot connect to the docker engine endpoint
```

The problem is incompatibilities between Docker version 20.10.0+ and Docker Machine version 0.16. This is the [official issue](https://github.com/docker/machine/issues/4858) on [GitHub Repository](https://github.com/docker/machine) of Docker Machine.

To install an older version of Docker when using Docker Machine, please run the following command: 

```text
docker-machine create --driver ionoscloud --engine-install-url "https://releases.rancher.com/install-docker/19.03.9.sh" test-machine
```

The output should be similar to this:

```text
Creating CA: /home/runner/.docker/machine/certs/ca.pem
Creating client certificate: /home/runner/.docker/machine/certs/cert.pem
Running pre-create checks...
Creating machine...
(test-machine) Creating SSH key...
(test-machine) DataCenter Created
(test-machine) LAN Created
(test-machine) Server Created
(test-machine) Volume Attached to Server
(test-machine) NIC Attached to Server
(test-machine) 158.222.102.181
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

When running `docker-machine ls`, the output should be similar to this:

```text
NAME           ACTIVE   DRIVER       STATE     URL                          SWARM   DOCKER     ERRORS
test-machine   -        ionoscloud   Running   tcp://158.222.102.185:2376           v19.03.9  
```

## Debug Option

Docker Machine has a `--debug` or `-D` option in order to get detailed output about the command running:

```text
docker-machine -D create --driver ionoscloud test-machine
```
