# Rancher Machine Commands 

## Commands

[Rancher Machine](https://github.com/rancher/machine) has multiple commands in order to manage a machine. 

| Command | Description |
| :------ | :---------- |
| active		| Print which machine is active |
| config		| Print the connection config for machine |
| create		| Create a machine |
| env			| Display the commands to set up the environment for the Docker client |
| inspect		| Inspect information about a machine |
| ip			| Get the IP address of a machine |
| kill			| Kill a machine |
| ls			| List machines |
| provision		| Re-provision existing machines |
| regenerate-certs	| Regenerate TLS Certificates for a machine |
| restart		| Restart a machine |
| rm			| Remove a machine |
| ssh			| Log into or run a command on a machine with SSH. |
| scp			| Copy files between machines |
| mount		    | Mount or unmount a directory from a machine with SSHFS. |
| start		    | Start a machine |
| status		| Get the status of a machine |
| stop			| Stop a machine |
| upgrade		| Upgrade a machine to the latest version of Docker |
| url			| Get the URL of a machine |
| version		| Show the Docker Machine version or a machine docker version |
| help			| Shows a list of commands or help for one command |

For more available options to manage a Rancher Machine, use `rancher-machine help`.

## Examples

### List Machines

To list the machines you have created, use the command:

```
rancher-machine ls
```

It will return information about your machines, similar to this:

```
NAME           ACTIVE   DRIVER         STATE     URL                         SWARM   DOCKER    ERRORS
test-machine   *        ionoscloud     Running   tcp://158.222.102.154:2376           v20.10.5
```

### Start a Machine

To start a Rancher Machine, run: 

```
rancher-machine start test-machine
```

Expected output:

```
Starting "test-machine"...
Machine "test-machine" was started.
Waiting for SSH to be available...
Detecting the provisioner...
Started machines may have new IP addresses. You may need to re-run the `rancher-machine env` command.
```

### Stop a Machine

To stop a Rancher Machine, run: 

```
rancher-machine stop test-machine
```

Expected output:

```
Stopping "test-machine"...
Machine "test-machine" was stopped.
```

### Restart a Machine

To restart a Rancher Machine, run: 

```
rancher-machine restart test-machine
```

Expected output:

```
Restarting "test-machine"...
Waiting for SSH to be available...
Detecting the provisioner...
Restarted machines may have new IP addresses. You may need to re-run the `rancher-machine env` command.
```

### Get Status

To get the status of a Rancher Machine created, run: 

```
rancher-machine status test-machine
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
