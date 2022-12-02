# Changelog

## \[6.1.0\]

* Multiline text for cloud init:
  * changed: `ionoscloud-user-data` parameter now takes multiline text as input. We recommend using the UI Driver for this parameter
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
