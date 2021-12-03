# Authentication

## Credentials

Before you use Ionos Cloud Docker Machine Driver, you need to authenticate with your Ionos Cloud credentials. These would be the same username and password that you use to log into the [Ionos Cloud DCD](https://dcd.ionos.com/latest/).

It is possible to pass your credentials:

* using environment variables:

```text
export IONOSCLOUD_USERNAME="ionoscloud_username"
export IONOSCLOUD_PASSWORD="ionoscloud_password"
```

or 

```text
export IONOSCLOUD_TOKEN="ionoscloud_token"
```

* on command-line using `--ionoscloud-username` and `--ionoscloud-password`:

```text
docker-machine create --driver=ionoscloud --ionoscloud-username="ionoscloud_username" --ionoscloud-password="ionoscloud_password" test-machine
```

or

```text
rancher-machine create --driver=ionoscloud --ionoscloud-username="ionoscloud_username" --ionoscloud-password="ionoscloud_password" test-machine
```

* on command-line using `--ionoscloud-token`:

```text
docker-machine create --driver=ionoscloud --ionoscloud-token="ionoscloud_token" test-machine
```

or

```text
rancher-machine create --driver=ionoscloud --ionoscloud-token="ionoscloud_token" test-machine
```

* on Rancher UI, when creating a new Node Template.

## API Endpoint

If you want to authenticate against a different API endpoint, you can set that:

* on command-line using `--ionoscloud-endpoint`

* using environment variable:

```text
export IONOSCLOUD_ENDPOINT="ionoscloud_endpoint"
```

* on Rancher UI, when creating a new Node Template.

It is recommended to use `api.ionos.com` as `IONOSCLOUD_ENDPOINT`, for flexibility across versions.

_Note_: SDK Go will check if the `/cloudapi/v6` suffix is set at the end of the API endpoint, and if not set, it will set it automatically.
