# Developing the Rancher driver and UI

## Download the repositories
* Driver: https://github.com/ionos-cloud/docker-machine-driver
* Old UI (RKE1): https://github.com/ionos-cloud/ui-driver-ionoscloud
* UI extensions (RKE2): https://github.com/ionos-cloud/ui-extensions-ionoscloud

## Install Rancher
You can install a Rancher using docker with the following command
```
docker run -d --name=rancher-server --restart=unless-stopped -p 80:80 -p 443:443 --privileged rancher/rancher:v2.7.5
```

> **_NOTE:_** Note that version 2.7+ is required for Rancher UI Extension support

> **_NOTE:_** Please note that versions 2.8+ do not currently allow adding the IonosCloud UI extension for RKE2

### Configure Rancher (for RKE2)
* Rancher URL

use the following command to receive the IP if the container
```
docker inspect <CONTAINER_ID> | grep IPAddress
```
Go to Global Setting and scroll down to `server-url`
change that address to the container IP
* Enabling developer load

Click the profile in the top right corner and go to preferences, check the box labeled `Enable Extension developer features`. This is required for loading local extensions.

## Serving the driver and UIs
* Old UI

Go the repo and run
```
npm start
```

* UI Extension

Go the repo build the extensions and serve them
```
yarn build-pkg node-driver
yarn serve-pkgs
```

* Driver

For Rancher to have access to the driver we need to open a local ftp server, for that you can use https://www.npmjs.com/package/serve.

go to the driver repo and build the driver
```
make build
```

the built driver will be in the bin folder, we can serve in using

```
serve bin/
```


## Adding the driver

Because adding the driver using the ui results in it having a random id generated for it, the driver must be added from the container using kubectl for RKE2. If you only need RKE1 feel free to add the driver using the UI.

First we must

Connect to the docker machine using
```
docker exec -it <CONTAINER_ID> bash
```

Create a YAML file containing the driver info
```
echo "apiVersion: management.cattle.io/v3
kind: NodeDriver
metadata:
  annotations:
    lifecycle.cattle.io/create.node-driver-controller: "true"
    privateCredentialFields: "token,username,password,endpoint"
  name: ionoscloud 
spec:
  active: false
  addCloudCredential: false
  builtin: false
  checksum: ""
  description: ""
  displayName: ionoscloud
  externalId: ""
  uiUrl: ""
  url: https://github.com/ionos-cloud/docker-machine-driver/releases/download/v6.1.4/docker-machine-driver-6.1.4-linux-amd64.tar.gz" > ionos.yaml
```
use kubetl to create a node driver resource
  ```
  kubectl create -f ionos.yaml
  ```

> **_NOTE:_** 'privateCredentialFields' must be set for Cloud credentials to work and must not be set for RKE1 template to work

## Adding the UI drivers

The old UI can the added by editing the driver in the node drivers list and changing the UI URL to where the component is available

The extension can be added by going to Extensions and clicking the 3 dots in the upper right of the page. Next click Developer Load and insert in url where the extensions is available. Check the persist box and click Load. A prompt ("Extensions changed - reload required") should appear in the upper right of the page, click on reload.
