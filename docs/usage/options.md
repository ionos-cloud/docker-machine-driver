# Options

To get more detailed information about the options below, run the command:

```text
rancher-machine create --help --driver ionoscloud
```

Available Options for the `rancher-machine create` command with the Rancher Driver: 

| Option | Description |
| :--- | :--- |
| --driver, -d | Driver to create machine with |
| --ionoscloud-datacenter-id | Ionos Cloud Virtual Data Center Id |
| --ionoscloud-disk-size | Ionos Cloud Volume Disk-Size \(10, 50, 100, 200, 400\) |
| --ionoscloud-disk-type | Ionos Cloud Volume Disk-Type \(HDD, SSD\) |
| --ionoscloud-endpoint | Ionos Cloud API Endpoint |
| --ionoscloud-image | Ionos Cloud Image Alias |
| --ionoscloud-location | Ionos Cloud Location |
| --ionoscloud-password | Ionos Cloud Password |
| --ionoscloud-server-availability-zone | Ionos Cloud Server Availability Zone \(AUTO, ZONE\_1, ZONE\_2, ZONE\_3\) |
| --ionoscloud-server-cores | Ionos Cloud Server Cores \(2, 3, 4, 5, 6, etc.\) |
| --ionoscloud-server-cpu-family | Ionos Cloud Server CPU families \(AMD\_OPTERON,INTEL\_XEON, INTEL\_SKYLAKE\) |
| --ionoscloud-server-ram | Ionos Cloud Server Ram \(1024, 2048, 3072, 4096, etc.\) |
| --ionoscloud-username | Ionos Cloud Username |
| --ionoscloud-volume-availability-zone | Ionos Cloud Volume Availability Zone \(AUTO, ZONE\_1, ZONE\_2, ZONE\_3\) |
| --swarm | Configure Machine to join a Swarm cluster |
| --swarm-addr | addr to advertise for Swarm \(default: detect and use the machine IP\) |
| --swarm-discovery | Discovery service to use with Swarm |
| --swarm-experimental | Enable Swarm experimental features |
| --swarm-host | ip/socket to listen on for Swarm master |
| --swarm-image | Specify Docker image to use for Swarm |
| --swarm-join-opt | Define arbitrary flags for Swarm join |
| --swarm-master | Configure Machine to be a Swarm master |
| --swarm-opt | Define arbitrary flags for Swarm master |
| --swarm-strategy | Define a default scheduling strategy for Swarm |
| --engine-env | Specify environment variables to set in the engine |
| --engine-insecure-registry | Specify insecure registries to allow with the created engine |
| --engine-install-url | Custom URL to use for engine installation |
| --engine-label | Specify labels for the created engine |
| --engine-opt | Specify arbitrary flags to include with the created engine in the form flag=value |
| --engine-registry-mirror | Specify registry mirrors to use |
| --engine-storage-driver | Specify a storage driver to use with the engine |
| --tls-san | Support extra SANs for TLS certs |

