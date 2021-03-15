# Authentication

Before you create a Ionos Cloud Rancher Machine you need to set two environment variables containing your Ionos Cloud credentials. 

These would be the same username and password that you use to log into the [Ionos Cloud DCD](https://dcd.ionos.com/latest/):

```text
export IONOSCLOUD_USERNAME="ionoscloud_username"
export IONOSCLOUD_PASSWORD="ionoscloud_password"
```

It is also possible to pass your credentials on the command-line using `--ionoscloud-username` and `--ionoscloud-password`:

```text
rancher-machine create --driver=ionoscloud --ionoscloud-username="ionoscloud_username" --ionoscloud-password="ionoscloud_password" test-machine
```

