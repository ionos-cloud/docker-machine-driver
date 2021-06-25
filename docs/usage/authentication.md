# Authentication

Before you use Ionos Cloud Docker Machine Driver, you need to authenticate with your Ionos Cloud credentials. These would be the same username and password that you use to log into the [Ionos Cloud DCD](https://dcd.ionos.com/latest/).

It is possible to pass your credentials:

* using environment variables:

```text
export IONOSCLOUD_USERNAME="ionoscloud_username"
export IONOSCLOUD_PASSWORD="ionoscloud_password"
```

* on command-line using `--ionoscloud-username` and `--ionoscloud-password`:

```text
docker-machine create --driver=ionoscloud --ionoscloud-username="ionoscloud_username" --ionoscloud-password="ionoscloud_password" test-machine
```

or

```text
rancher-machine create --driver=ionoscloud --ionoscloud-username="ionoscloud_username" --ionoscloud-password="ionoscloud_password" test-machine
```

* on [Rancher UI](), when creating a new Node Template.

