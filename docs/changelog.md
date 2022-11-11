# Changelog

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
