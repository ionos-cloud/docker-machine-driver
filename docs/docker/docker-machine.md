# Docker Machine

## Authentication

Before you create a Docker Machine with Ionos Cloud Docker Machine Driver, make sure you followed the [Authentication](../usage/authentication.md) steps.

## Docker Machine Create Command

In order to create Docker Machine with Ionos Cloud Docker Machine Driver, run:

```text
docker-machine create --driver ionoscloud test-machine
```

For Windows users:

```text
docker-machine create --driver ionoscloud.exe test-machine
```

It should produce results similar to this:

```text
Running pre-create checks...
Creating machine...
(test-machine) Creating SSH key...
(test-machine) DataCenter Created
(test-machine) IPBlock Reserved
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
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env test-machine
```

After creating a machine it is recommended to configure your shell, to set up your environment for the Docker client:

```text
eval $(docker-machine env test-machine)
```

All the resources created will be named with the machine name, in this example `test-machine`.

The example above uses mostly the default values and the resources will therefore be created in the `us/las` location. To change that or to see more options that can be used with this command, check the [Options](../usage/options.md) section.

For more available commands and examples on how to manage a Docker Machine, check the [Commands](../usage/commands.md) section.

For more details about possible issues, check the [Troubleshooting](troubleshooting.md) section.

