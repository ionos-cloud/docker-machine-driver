# Rancher Cluster

IONOS Cloud Docker Machine Driver is compatible with [Rancher](https://rancher.com/).

## Installation

To install Rancher and Log in to Rancher UI, follow the first 3 steps in this [Quick Manual Setup](https://rancher.com/docs/rancher/v2.x/en/quick-start-guide/deployment/quickstart-manual-setup/).

You will create a Docker Container with the following command:

```text
sudo docker run -d --restart=unless-stopped -p 80:80 -p 443:443 --privileged rancher/rancher
```

> **_NOTE:_** Please note that versions 2.8+ do not currently allow adding the IonosCloud UI extension for RKE2

To use a specific Rancher version, check the [available docker images](https://hub.docker.com/r/rancher/rancher/tags) and add the corresponding tag to the command: 

```text
sudo docker run -d --restart=unless-stopped -p 80:80 -p 443:443 --privileged rancher/rancher:v2.7.5
```

To output the available docker containers, use:

```text
docker ps
```

To follow the output logs for the running container, use:

```text
docker logs -f container-id
```

## Prerequisites

* Your IONOS Cloud account credentials: username and password or token
* A web server accessible by your browser


## Installing Via The Rancher UI

After logging into Rancher UI, follow the next steps in order to install a cluster with IONOS Cloud as cloud provider, using IONOS Cloud Docker Machine Driver:

### RKE1

#### Adding the Node Driver

* Install Node Driver
  * Go to Tools ➜ Drivers ➜ Node Drivers
  * Click on `Add New Driver` button
  * Enter the URLs and click `Create`
    * Download URL: https://github.com/ionos-cloud/docker-machine-driver/releases/download/v<version>/docker-machine-driver-<version>-linux-amd64.tar.gz
    * Custom UI URL:  https://cdn.jsdelivr.net/gh/ionos-cloud/ui-driver-ionoscloud@main/releases/v<UI_version|latest>/component.js
    * Whitelist Domains: cdn.jsdelivr.net
  * Wait fot the machine driver to be downloaded and become `Active`
* Create Node Template
  * Go to Node Templates, from the drop-down menu for `User Settings`
  * Click on `Add Node Template` button
  * At this point, `Ionoscloud` should be on the list of `Available Hosts`. Select `Ionoscloud`
  * Configure the `IONOSCLOUD OPTIONS` as you prefer and add also your password and username for IONOS Cloud account
  * Give a name to the new Node Template and press `Create` button
* Create New Rancher Cluster
  * Go to Clusters
  * Click on `Add Cluster` button
  * In the `Create a new Kubernetes cluster` section, select `Ionoscloud`
  * Choose the name of the new cluster, the name prefix of the node and make sure you have the Node Template you just created, in the `Template` section
  * Customize your cluster: Single Node \(by selecting all etcd, Control Plane and Worker\) or Multiple Nodes
  * Click on `Create` button
  * Wait for cluster to become `Active` \(it will take some minutes\).
  
### RKE2

Using Rancher Extensions requires Rancher v2.7.0 or above.

* Install Node Driver
  * connect to the machine running rancher
  * create a yaml file containing the following information:
  ```yaml
  apiVersion: management.cattle.io/v3
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
    url: <IONOS_DRIVER_URL>
    ```
  * create the driver resource using
  ```
  kubectl create -f <FILE>
  ```

  * you can also add the old UI if you want to use RKE1
    * Go to Tools ➜ Drivers ➜ Node Drivers
    * Edit the Ionoscloud driver
    * Custom UI URL:  https://cdn.jsdelivr.net/gh/ionos-cloud/ui-driver-ionoscloud@main/releases/v<UI_version|latest>/component.js
    * Whitelist Domains: cdn.jsdelivr.net
  * Wait for the machine driver to be downloaded and become `Active`
  * Add the ionoscloud ui extension from https://github.com/ionos-cloud/ui-extensions-ionoscloud
    * Go to Cluster Management ➜ Advanced ➜ Repositories
    * Click on `Create` button
    * Select Git repository as target
    * Git Repo URL: https://github.com/ionos-cloud/ui-extensions-ionoscloud
    * Git Branch: gh-pages

* Create Cloud Credential
  * Go to Cluster Management ➜ Cloud Credentials
  * Click on `Create` button
  * At this point, `Ionoscloud` should be on the list. Select `Ionoscloud`
* Create New Rancher Cluster
  * Go to Cluster Management ➜ Clusters
  * Click on `Create` button
  * In the `Create` section, select `Ionoscloud`
  * Customize your cluster
  * Click on `Create` button
  * Wait for cluster to become `Active` \(it will take some minutes\).
  
## Support

Please submit any bugs, issues or feature requests to [ionos-cloud/docker-machine-driver](https://github.com/ionos-cloud/docker-machine-driver/issues/new/choose).

