# Examples

## Using on command-line options

```text
docker-machine create \
    --driver ionoscloud \
    --ionoscloud-token="ionoscloud_token" \
    test-machine
```

```text
docker-machine create \
    --driver ionoscloud \
    --ionoscloud-username="ionoscloud_username" \
    --ionoscloud-password="ionoscloud_password" \
    test-machine
```

## Using environment variables

```text
export IONOSCLOUD_TOKEN="ionoscloud_token"

docker-machine create \
    --driver ionoscloud \
    test-machine
```

## Using a specific image alias

```text
docker-machine create \
    --driver ionoscloud \
    --ionoscloud-token="ionoscloud_token" \
    --ionoscloud-image="ubuntu:18.04" \
    test-machine
```
