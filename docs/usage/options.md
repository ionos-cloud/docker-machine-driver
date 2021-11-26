# Options

To get more detailed information about the options and the environment variables available, run the command:

```text
docker-machine create --help --driver ionoscloud
```

or

```text
rancher-machine create --help --driver ionoscloud
```

## Options

Available Options for the IONOS Cloud Docker Machine Driver:

| Option | Description |
| :--- | :--- |
| `--driver, -d` | Driver to create machine with |
| `--ionoscloud-username` | Ionos Cloud Username |
| `--ionoscloud-password` | Ionos Cloud Password |
| `--ionoscloud-endpoint` | Ionos Cloud API Endpoint. It is recommended to be set to `https://api.ionos.com` or `https://api.ionos.com/cloudapi/v5`. The SDK will automatically put the `/cloudapi/v5` suffix if not set. |
| `--ionoscloud-datacenter-id` | Ionos Cloud Virtual Data Center Id |
| `--ionoscloud-disk-size` | Ionos Cloud Volume Disk-Size \(10, 50, 100, 200, 400\) |
| `--ionoscloud-disk-type` | Ionos Cloud Volume Disk-Type \(HDD, SSD\) |
| `--ionoscloud-image` | Ionos Cloud Image Alias \(ubuntu:latest, ubuntu:20.04\) |
| `--ionoscloud-image-password` | Ionos Cloud Image Password to be able to access the server from DCD platform |
| `--ionoscloud-location` | Ionos Cloud Location |
| `--ionoscloud-server-availability-zone` | Ionos Cloud Server Availability Zone \(AUTO, ZONE\_1, ZONE\_2, ZONE\_3\) |
| `--ionoscloud-cores` | Ionos Cloud Server Cores \(2, 3, 4, 5, 6, etc.\) |
| `--ionoscloud-cpu-family` | Ionos Cloud Server CPU families \(AMD\_OPTERON,INTEL\_XEON, INTEL\_SKYLAKE\) |
| `--ionoscloud-ram` | Ionos Cloud Server Ram \(1024, 2048, 3072, 4096, etc.\) |
| `--ionoscloud-volume-availability-zone` | Ionos Cloud Volume Availability Zone \(AUTO, ZONE\_1, ZONE\_2, ZONE\_3\) |
| `--swarm` | Configure Machine to join a Swarm cluster |
| `--swarm-addr` | addr to advertise for Swarm \(default: detect and use the machine IP\) |
| `--swarm-discovery` | Discovery service to use with Swarm |
| `--swarm-experimental` | Enable Swarm experimental features |
| `--swarm-host` | ip/socket to listen on for Swarm master |
| `--swarm-image` | Specify Docker image to use for Swarm |
| `--swarm-join-opt` | Define arbitrary flags for Swarm join |
| `--swarm-master` | Configure Machine to be a Swarm master |
| `--swarm-opt` | Define arbitrary flags for Swarm master |
| `--swarm-strategy` | Define a default scheduling strategy for Swarm |
| `--engine-env` | Specify environment variables to set in the engine |
| `--engine-insecure-registry` | Specify insecure registries to allow with the created engine |
| `--engine-install-url` | Custom URL to use for engine installation |
| `--engine-label` | Specify labels for the created engine |
| `--engine-opt` | Specify arbitrary flags to include with the created engine in the form flag=value |
| `--engine-registry-mirror` | Specify registry mirrors to use |
| `--engine-storage-driver` | Specify a storage driver to use with the engine |
| `--tls-san` | Support extra SANs for TLS certs |

## Environment variables

Environment variables are also supported for setting options. This is a list of the environment variables available for Docker Machine Driver.

| Option | Environment variable |
| :--- | :--- |
| `--ionoscloud-username` | `IONOSCLOUD_USERNAME` |
| `--ionoscloud-password` | `IONOSCLOUD_PASSWORD` |
| `--ionoscloud-endpoint` | `IONOSCLOUD_ENDPOINT` |
| `--ionoscloud-datacenter-id` | `IONOSCLOUD_DATACENTER_ID` |
| `--ionoscloud-disk-size` | `IONOSCLOUD_DISK_SIZE` |
| `--ionoscloud-disk-type` | `IONOSCLOUD_DISK_TYPE` |
| `--ionoscloud-image` | `IONOSCLOUD_IMAGE` |
| `--ionoscloud-image-password` | `IONOSCLOUD_IMAGE_PASSWORD` |
| `--ionoscloud-location` | `IONOSCLOUD_LOCATION` |
| `--ionoscloud-server-availability-zone` | `IONOSCLOUD_SERVER_ZONE` |
| `--ionoscloud-cores` | `IONOSCLOUD_CORES` |
| `--ionoscloud-cpu-family` | `IONOSCLOUD_CPU_FAMILY` |
| `--ionoscloud-ram` | `IONOSCLOUD_RAM` |
| `--ionoscloud-volume-availability-zone` | `IONOSCLOUD_VOLUME_ZONE` |

