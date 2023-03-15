# Changelog

## \[6.1.3\]
### Added
- Added `ionoscloud-nic-dhcp` and  `ionoscloud-nic-ips` which allow to change the properties of the NIC which will be created
### Fixes
- Fixed some issues regarding IP allocation when creating a NAT Gateway
### Changed
* default NIC DHCP is now set to false


## \[6.1.2\]
### Added
- Added `ionoscloud-nat-id` and  `ionoscloud-nat-name` which allows using a pre-configured NAT Gateway if it exists, by searching for the NAT with the given name in the given Datacenter. Setting the id will ignore the name flag
 - Added `ionoscloud-create-nat` which creates a NAT if set
- Added customization options for the NAT `ionoscloud-nat-public-ips` (a list of Public IPs) and `ionoscloud-nat-lans-to-gateways` (mappings of Lans to Gateway IPs) which are optional (used in conjunction with `create-nat`)
- Added `ionoscloud-private-lan` which, if set to True, will make the default LAN be private. Note that creating a NAT required a private LAN, so either set this to true, or provide an already existing private LAN via `ionoscloud-lan-id` or `ionoscloud-lan-name`.
## Use with Rancher
    * Download URL: https://github.com/ionos-cloud/docker-machine-driver/releases/download/v6.1.2/docker-machine-driver-6.1.2-linux-amd64.tar.gz
    * Custom UI URL:  https://cdn.jsdelivr.net/gh/ionos-cloud/ui-driver-ionoscloud@main/releases/latest/component.js
    * Whitelist Domains: cdn.jsdelivr.ne


## \[6.1.1\]
### Fixed
- Fixed `failed to create server due to error: [(root).entities.volumes.items.[0].properties.sshKeys] Invalid SSH key. Maximum allowed key size is 8K (8192 characters) and it can not be empty. Given ssh key length: 0 characters`

## \[6.1.0\]

### Added

* Added the [**IONOS UI Driver**](https://github.com/ionos-cloud/ui-driver-ionoscloud), for users of the Rancher docker image. To use the custom UI, use following fields when adding the driver:
    * Custom UI URL:  https://cdn.jsdelivr.net/gh/ionos-cloud/ui-driver-ionoscloud@main/releases/latest/component.js
    * Whitelist Domains: cdn.jsdelivr.net
    
    We highly recommend using this UI Driver if you are using the Rancher docker image.

* Added the option to customize the SSH User that Rancher uses to connect to the Docker Host (`ionoscloud-ssh-user`). [#49](https://github.com/ionos-cloud/docker-machine-driver/pull/49)
* Added the option to select an existing LAN in which to provision the Docker Host (`ionoscloud-lan-id`). Using this option requires you to set the Datacenter ID as well (`ionoscloud-datacenter-id`). [#42](https://github.com/ionos-cloud/docker-machine-driver/pull/42)
* Added support for CUBE servers (#63)
* Added ability to select existing LAN or Datacenter by name. (#54) 

### Changed

* Changed cloud-init parameter behaviour:
  * changed: `ionoscloud-user-data` parameter now takes multiline text as input.
  * added: `ionoscloud-user-data-b64` flag, which takes a b64 encoded string. This field will only be evaluated if `ionoscloud-user-data` is empty.
* Changed default image alias to ubuntu20, as currently the Docker Machine Driver only supports id-rsa ssh keys, which cannot be used to connect to ubuntu22 VMs

### Fixed

* Fixed error messages getting cut off at the newline mark
* Fixes related to user data: Cloud Config YAML was not encoded, and users would be duplicated if the ssh User wasn't root

### Known Issues

* Currently, ubuntu:22.04 (aka ubuntu:latest) is unsupported for the Rancher docker image.


## \[6.1.0-rc.2\]

### Fixes
* error messages cutoff
* user data not encoded when sshUser is root


## \[6.1.0-rc.1\]

### Added

* Added the [**IONOS UI Driver**](https://github.com/ionos-cloud/ui-driver-ionoscloud), for users of the Rancher docker image. To use the custom UI, use following fields when adding the driver:
    * Custom UI URL:  https://cdn.jsdelivr.net/gh/ionos-cloud/ui-driver-ionoscloud@0.1.0/releases/v0.1.0/component.js
    * Whitelist Domains: cdn.jsdelivr.net
    
    We highly recommend using this UI Driver if you are using the Rancher docker image.

* Added the option to customize the SSH User that Rancher uses to connect to the Docker Host (`ionoscloud-ssh-user`). [#49](https://github.com/ionos-cloud/docker-machine-driver/pull/49)
* Added the option to select an existing LAN in which to provision the Docker Host (`ionoscloud-lan-id`). Using this option requires you to set the Datacenter ID as well (`ionoscloud-datacenter-id`). [#42](https://github.com/ionos-cloud/docker-machine-driver/pull/42)

### Changed

* Changed cloud-init parameter behaviour:
  * changed: `ionoscloud-user-data` parameter now takes multiline text as input.
  * added: `ionoscloud-user-data-b64` flag, which takes a b64 encoded string. This field will only be evaluated if `ionoscloud-user-data` is empty.


### Known Issues

* Currently, ubuntu:22.04 (aka ubuntu:latest) is unsupported for the Rancher docker image.


## \[6.0.1\]

* Added: `user-data` parameter support for volumes: 
You can now export `IONOSCLOUD_USER_DATA` or use flag `ionoscloud-user-data` to set the volume's cloud init. 
Needs to be a base64 encoded string
* dependency-updates:
	- SDK Go `v6.1.3`
	- go `1.18`

## \[6.0.0\]

* first release of Docker Machine Driver using SDK Go v6 🎉
* dependency-update: added SDK Go `v6.0.0`
