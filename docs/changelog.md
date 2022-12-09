# Changelog

## \[6.1.0-rc.1\]

### Added

* Added the [**IONOS UI Driver**](https://github.com/ionos-cloud/ui-driver-ionoscloud), for users of the Rancher docker image. To use the custom UI, use following fields when adding the driver:
    * Custom UI URL:  https://cdn.jsdelivr.net/gh/ionos-cloud/ui-driver-ionoscloud@0.1.0/releases/v0.1.0/component.js
    * Whitelist Domains: https://cdn.jsdelivr.net
    
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
