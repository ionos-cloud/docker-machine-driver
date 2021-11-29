# Rancher Cluster

IONOS Cloud Docker Machine Driver is compatible with [Rancher](https://rancher.com/).

## Installation

To install Rancher and Log in to Rancher UI, follow the first 3 steps in this [Quick Manual Setup](https://rancher.com/docs/rancher/v2.x/en/quick-start-guide/deployment/quickstart-manual-setup/).

You will create a Docker Container with the following command:

```text
sudo docker run -d --restart=unless-stopped -p 80:80 -p 443:443 --privileged rancher/rancher
```

To output the available containers, use:

```text
docker ps
```

To follow the output logs for the running container, use:

```text
docker logs -f container-id
```

## Prerequisites

* Your IONOS Cloud account credentials: username and password or token
* A link to download the latest release of docker-machine-driver for your platform. Check the [Release Page](https://github.com/ionos-cloud/docker-machine-driver/releases) for the latest releases available. Find the corresponding one for your operating system and architecture. 

All you need to do is to copy the link location, which should look something like this: `https://github.com/ionos-cloud/docker-machine-driver/releases/download/v<version>/docker-machine-driver-<version>-linux-amd64.tar.gz`

* A web server accessible by your browser

## Installing Via The Rancher UI

After logging into Rancher UI, follow the next steps in order to install a cluster with IONOS Cloud as cloud provider, using IONOS Cloud Docker Machine Driver:

* Install Node Driver
  * Go to Tools ➜ Drivers ➜ Node Drivers
  * Click on `Add New Driver` button
  * Enter the URL for the docker-machine-driver in the `Download URL` and click `Create`
  * Wait fot the machine driver to be downloaded and become `Active`
* Create Node Template
  * Go to Node Templates, from the drop down menu for `User Settings`
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

## Support

Please submit any bugs, issues or feature requests to [ionos-cloud/docker-machine-driver](https://github.com/ionos-cloud/docker-machine-driver/issues/new/choose).

