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

| Option                                  | Description                                                                                                                                                                                   |
|:----------------------------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--driver, -d`                          | Driver to create machine with                                                                                                                                                                 |
| `--ionoscloud-username`                 | Ionos Cloud Username                                                                                                                                                                          |
| `--ionoscloud-password`                 | Ionos Cloud Password                                                                                                                                                                          |
| `--ionoscloud-token`                    | Ionos Cloud Token                                                                                                                                                                             |
| `--ionoscloud-endpoint`                 | Ionos Cloud API Endpoint. It is recommended to be set to `https://api.ionos.com` or `https://api.ionos.com/cloudapi/v6`. The SDK will automatically put the `/cloudapi/v6` suffix if not set. |
| `--ionoscloud-datacenter-id`            | Existing Ionos Cloud Virtual Data Center ID (UUID-4) in which to create the Docker Host                                                                                                       |
| `--ionoscloud-datacenter-name`          | Existing Ionos Cloud Virtual Data Center Name (string) in which to create the Docker Host                                                                                                     |
| `--ionoscloud-lan-id`                   | Existing Ionos Cloud LAN ID (numeric) in which to create the Docker Host                                                                                                                      |
| `--ionoscloud-lan-name`                 | Existing Ionos Cloud LAN Name (string) in which to create the Docker Host                                                                                                                     |
| `--ionoscloud-additional-lans`          | Names of existing IONOS Lans to connect the machine to. Names that are not found are ignored                                                                                                                     |
| `--ionoscloud-disk-size`                | Ionos Cloud Volume Disk-Size in GB \(10, 50, 100, 200, 400\)                                                                                                                                  |
| `--ionoscloud-disk-type`                | Ionos Cloud Volume Disk-Type \(HDD, SSD\)                                                                                                                                                     |
| `--ionoscloud-image`                    | Ionos Cloud Image Id or Alias \(ubuntu:latest, debian:latest, etc.\). If Image Id is set, please make sure the disk type supports the image type.                                                    |
| `--ionoscloud-image-password`           | Ionos Cloud Image Password to be able to access the server from DCD platform                                                                                                                  |
| `--ionoscloud-location`                 | Ionos Cloud Location                                                                                                                                                                          |
| `--ionoscloud-server-type`              | Ionos Cloud Server Type (ENTERPRISE or CUBE)                                                                                                                                                  |
| `--ionoscloud-template`                 | Ionos Cloud Template name (CUBES XS, CUBES S, etc.)                                                                                                                                           |
| `--ionoscloud-server-availability-zone` | Ionos Cloud Server Availability Zone \(AUTO, ZONE\_1, ZONE\_2, ZONE\_3\)                                                                                                                      |
| `--ionoscloud-cores`                    | Ionos Cloud Server Cores \(2, 3, 4, 5, 6, etc.\)                                                                                                                                              |
| `--ionoscloud-cpu-family`               | Ionos Cloud Server CPU families \(INTEL\_XEON, INTEL\_SKYLAKE, INTEL\_ICELAKE, AMD\_EPYC, INTEL\_SIERRAFOREST)                                                                                                                  |
| `--ionoscloud-ram`                      | Ionos Cloud Server Ram in MB \(1024, 2048, 3072, 4096, etc.\)                                                                                                                                 |
| `--ionoscloud-volume-availability-zone` | Ionos Cloud Volume Availability Zone \(AUTO, ZONE\_1, ZONE\_2, ZONE\_3\)                                                                                                                      |
| `--ionoscloud-cloud-init`               | The cloud-init configuration for the volume as multiline text                                                                                                                                 |
| `--ionoscloud-cloud-init-b64`           | The cloud-init configuration for the volume as base64 encoded string. Prioritized                                                                                                             |
| `--ionoscloud-nic-multi-queue`                 | Activate or deactivate the Multi Queue feature on all NICs of this server., defaults to false |
| `--ionoscloud-nic-dhcp`                 | Wether the created NIC should have DHCP set, defaults to false |
| `--ionoscloud-nic-ips`                  | The ips used for the nic                                                                                                                                                           |
| `--ionoscloud-wait-for-ip-change`                  | Should the driver wait for the NIC IP to be set by external sources?                                                                                                                                                           |
| `--ionoscloud-wait-for-ip-change-timeout`                  | Timeout used when waiting for NIC IP changes                                                                                                                                                           |
| `--ionoscloud-nat-id`                   | Use an existing NAT via its ID                                                                                                                                                                |
| `--ionoscloud-nat-name`                 | Use an existing NAT via its name                                                                                                                                                              |
| `--ionoscloud-create-nat`               | Create a new NAT with some default open ports                                                                                                                                                 |
| `--ionoscloud-nat-public-ips`           | If --ionoscloud-create-nat is set, change the NAT's public IPs to these values                                                                                                                |
| `--ionoscloud-nat-lans-to-gateways`     | If --ionoscloud-create-nat is set, change the NAT's mappings of LANs to Gateway IPs to these values. Must respect format `1=10.0.0.1,10.0.0.2:2=10.0.0.10`                                    |
| `--ionoscloud-nat-flowlogs`     | If --ionoscloud-create-nat is set, add flowlogs to the nat. Must respect format `name:action:direction:bucket`,                                    |
| `--ionoscloud-nat-rules`     | If --ionoscloud-create-nat is set, add rules to the NAT. Must respect format `name:type:protocol:public_ip:source_subnet:target_subnet:target_port_range_start:target_port_range_end`, to skip providing an optional value just omit it (`name:type:protocol::source_subnet:::`), not setting public IP will use the public IP of the NAT for the rule, not setting source subnet will use the first ip on the NIC with mask 24                                    |
| `--ionoscloud-skip-default-nat-rules`                 | Should the driver skip creating default nat rules if creating a NAT, creating only the specified rules, the UI drivers always set this flag                                                                                                                                                                |
| `--ionoscloud-ssh-user`                 | The user to connect to via SSH                                                                                                                                                                |
| `--ionoscloud-ssh-in-cloud-init`        | Should the driver only add the SSH info in the user data? (required for custom images)                                                                                                                                                                |
| `--ionoscloud-rancher-provision-user-data`        | Placeholder flag for rancher machine creation flow to populate with rke2 install user-data instructions                                                                                                                                                                |
| `--ionoscloud-append-rke-userdata`        | Should the driver append the rke user-data to the user-data sent to the ionos server                                                                                                                                                                |
| `--swarm`                               | Configure Machine to join a Swarm cluster                                                                                                                                                     |
| `--swarm-addr`                          | addr to advertise for Swarm \(default: detect and use the machine IP\)                                                                                                                        |
| `--swarm-discovery`                     | Discovery service to use with Swarm                                                                                                                                                           |
| `--swarm-experimental`                  | Enable Swarm experimental features                                                                                                                                                            |
| `--swarm-host`                          | ip/socket to listen on for Swarm master                                                                                                                                                       |
| `--swarm-image`                         | Specify Docker image to use for Swarm                                                                                                                                                         |
| `--swarm-join-opt`                      | Define arbitrary flags for Swarm join                                                                                                                                                         |
| `--swarm-master`                        | Configure Machine to be a Swarm master                                                                                                                                                        |
| `--swarm-opt`                           | Define arbitrary flags for Swarm master                                                                                                                                                       |
| `--swarm-strategy`                      | Define a default scheduling strategy for Swarm                                                                                                                                                |
| `--engine-env`                          | Specify environment variables to set in the engine                                                                                                                                            |
| `--engine-insecure-registry`            | Specify insecure registries to allow with the created engine                                                                                                                                  |
| `--engine-install-url`                  | Custom URL to use for engine installation                                                                                                                                                     |
| `--engine-label`                        | Specify labels for the created engine                                                                                                                                                         |
| `--engine-opt`                          | Specify arbitrary flags to include with the created engine in the form flag=value                                                                                                             |
| `--engine-registry-mirror`              | Specify registry mirrors to use                                                                                                                                                               |
| `--engine-storage-driver`               | Specify a storage driver to use with the engine                                                                                                                                               |
| `--tls-san`                             | Support extra SANs for TLS certs                                                                                                                                                              |

## Environment variables

Environment variables are also supported for setting options. This is a list of the environment variables available for Docker Machine Driver.

| Option                                    | Environment variable                   |
|:------------------------------------------|:---------------------------------------|
| `--ionoscloud-username`                   | `IONOSCLOUD_USERNAME`                  |
| `--ionoscloud-password`                   | `IONOSCLOUD_PASSWORD`                  |
| `--ionoscloud-token`                      | `IONOSCLOUD_TOKEN`                     |
| `--ionoscloud-endpoint`                   | `IONOSCLOUD_ENDPOINT`                  | 
| `--ionoscloud-datacenter-id`              | `IONOSCLOUD_DATACENTER_ID`             |
| `--ionoscloud-datacenter-name`            | `IONOSCLOUD_DATACENTER_NAME`           |
| `--ionoscloud-lan-id`                     | `IONOSCLOUD_LAN_ID`                    |
| `--ionoscloud-lan-name`                   | `IONOSCLOUD_LAN_NAME`                  |
| `--ionoscloud-additional-lans`            | `IONOSCLOUD_ADDITIONAL_LANS`           |
| `--ionoscloud-disk-size`                  | `IONOSCLOUD_DISK_SIZE`                 |
| `--ionoscloud-disk-type`                  | `IONOSCLOUD_DISK_TYPE`                 |
| `--ionoscloud-image`                      | `IONOSCLOUD_IMAGE`                     |
| `--ionoscloud-image-password`             | `IONOSCLOUD_IMAGE_PASSWORD`            |
| `--ionoscloud-server-type`                | `IONOSCLOUD_SERVER_TYPE`               |
| `--ionoscloud-template`                   | `IONOSCLOUD_TEMPLATE`                  |
| `--ionoscloud-location`                   | `IONOSCLOUD_LOCATION`                  |
| `--ionoscloud-server-availability-zone`   | `IONOSCLOUD_SERVER_AVAILABILITY_ZONE`  |
| `--ionoscloud-cores`                      | `IONOSCLOUD_CORES`                     |
| `--ionoscloud-cpu-family`                 | `IONOSCLOUD_CPU_FAMILY`                |
| `--ionoscloud-ram`                        | `IONOSCLOUD_RAM`                       |
| `--ionoscloud-volume-availability-zone`   | `IONOSCLOUD_VOLUME_AVAILABILITY_ZONE`  |
| `--ionoscloud-cloud-init`                 | `IONOSCLOUD_CLOUD_INIT`                |
| `--ionoscloud-cloud-init-b64`             | `IONOSCLOUD_CLOUD_INIT_B64`            |
| `--ionoscloud-nic-multi-queue`            | `IONOSCLOUD_NIC_MULTI_QUEUE`           |
| `--ionoscloud-nic-dhcp`                   | `IONOSCLOUD_NIC_DHCP`                  |
| `--ionoscloud-nic-ips`                    | `IONOSCLOUD_NIC_IPS`                   |
| `--ionoscloud-wait-for-ip-change`         | `IONOSCLOUD_WAIT_FOR_IP_CHANGE`        |
| `--ionoscloud-wait-for-ip-change-timeout` | `IONOSCLOUD_WAIT_FOR_IP_CHANGE_TIMEOUT`|
| `--ionoscloud-create-nat`                 | `IONOSCLOUD_CREATE_NAT`                |
| `--ionoscloud-nat-name`                   | `IONOSCLOUD_NAT_NAME`                  |
| `--ionoscloud-nat-public-ips`             | `IONOSCLOUD_NAT_PUBLIC_IPS`            |
| `--ionoscloud-nat-lans-to-gateways`       | `IONOSCLOUD_NAT_LANS_TO_GATEWAYS`      |
| `--ionoscloud-nat-flowlogs`               | `IONOSCLOUD_NAT_FLOWLOG`               |
| `--ionoscloud-nat-rules`                  | `IONOSCLOUD_NAT_RULES`                 |
| `--ionoscloud-skip-default-nat-rules`     | `IONOSCLOUD_SKIP_DEFAULT_NAT_RULES`    |
| `--ionoscloud-private-lan`                | `IONOSCLOUD_PRIVATE_LAN`               |
| `--ionoscloud-ssh-user`                   | `IONOSCLOUD_SSH_USER`                  |
| `--ionoscloud-ssh-in-cloud-init`          | `IONOSCLOUD_SSH_IN_CLOUD_INIT`         |
| `--ionoscloud-rancher-provision-user-data`| `IONOSCLOUD_RANCHER_PROVISION_USERDATA`|
| `--ionoscloud-append-rke-userdata`        | `IONOSCLOUD_APPEND_RKE_USERDATA`       |
