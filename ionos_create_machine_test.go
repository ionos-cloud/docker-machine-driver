package ionoscloud

import (
	"context"
	"encoding/base64"
	"os"
	"testing"

	"github.com/ionos-cloud/docker-machine-driver/internal/utils"
	"github.com/stretchr/testify/assert"
)

// This is an integration test to verify that the final user data is correctly generated when using RKE provisioning.
// It ensures that the RKE provisioning script is correctly appended to the user data.
// rancher-machine creates a temp file and sets the path into flag, therefore we also have to create a temp file here
func TestGetFinalUserDataWithRKEProvision(t *testing.T) {
	driver, _ := NewTestDriver(t, "test-host", "defaultstore")
	driver.SSHUser = "root"
	driver.client = func() utils.ClientService {
		return utils.New(context.TODO(), driver.Username, driver.Password, driver.Token, driver.Endpoint, "user-agent")
	}
	driver.AppendRKEProvisionUserData = true
	driver.CloudInit = `#cloud-config
hostname: test.example.com
packages:
  - somepackage
runcmd:
- sh user_script.sh
write_files:
- path: /etc/user_script.sh
  content: some user content
`
	tmpFile, err := os.CreateTemp("", "rke-provision-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rkeContent := `#cloud-config
runcmd:
- sh /etc/rke.sh
write_files:
- path: /etc/rke.sh
  content: some install content
`
	if _, err := tmpFile.WriteString(rkeContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	driver.RKEProvisionUserData = tmpFile.Name()

	expectedResult := `#cloud-config
hostname: test.example.com
packages:
    - somepackage
runcmd:
    - sh user_script.sh
    - sh /etc/rke.sh
users:
    - create_groups: false
      lock_passwd: true
      name: root
      no_user_group: true
      ssh_authorized_keys:
        - ""
      sudo: ALL=(ALL) NOPASSWD:ALL
write_files:
    - content: some user content
      path: /etc/user_script.sh
    - content: some install content
      path: /etc/rke.sh
`
	result, err := driver.GetFinalUserData()
	assert.NoError(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(expectedResult)), result)
}

func TestGetFinalUserDataWithRKEProvisionEmptyCloudInit(t *testing.T) {
	driver, _ := NewTestDriver(t, "test-host", "defaultstore")
	driver.SSHUser = "root"
	driver.client = func() utils.ClientService {
		return utils.New(context.TODO(), driver.Username, driver.Password, driver.Token, driver.Endpoint, "user-agent")
	}
	driver.AppendRKEProvisionUserData = true

	tmpFile, err := os.CreateTemp("", "rke-provision-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rkeContent := `#cloud-config
runcmd:
- sh /etc/rke.sh
write_files:
- path: /etc/rke.sh
  content: some install content
`
	if _, err := tmpFile.WriteString(rkeContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	driver.RKEProvisionUserData = tmpFile.Name()

	expectedResult := `#cloud-config
hostname: test-host
runcmd:
    - sh /etc/rke.sh
users:
    - create_groups: false
      lock_passwd: true
      name: root
      no_user_group: true
      ssh_authorized_keys:
        - ""
      sudo: ALL=(ALL) NOPASSWD:ALL
write_files:
    - content: some install content
      path: /etc/rke.sh
`
	result, err := driver.GetFinalUserData()
	assert.NoError(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(expectedResult)), result)
}

func TestGetFinalUserDataWithRKEProvisionFlagFalse(t *testing.T) {
	driver, _ := NewTestDriver(t, "test-host", "defaultstore")
	driver.SSHUser = "root"
	driver.client = func() utils.ClientService {
		return utils.New(context.TODO(), driver.Username, driver.Password, driver.Token, driver.Endpoint, "user-agent")
	}
	driver.CloudInit = `#cloud-config
packages:
  - somepackage
runcmd:
- sh user_script.sh
write_files:
- path: /etc/user_script.sh
  content: some user content
`
	tmpFile, err := os.CreateTemp("", "rke-provision-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rkeContent := `#cloud-config
runcmd:
- sh /etc/rke.sh
write_files:
- path: /etc/rke.sh
  content: some install content
`
	expectedResult := `#cloud-config
hostname: test-host
packages:
    - somepackage
runcmd:
    - sh user_script.sh
users:
    - create_groups: false
      lock_passwd: true
      name: root
      no_user_group: true
      ssh_authorized_keys:
        - ""
      sudo: ALL=(ALL) NOPASSWD:ALL
write_files:
    - content: some user content
      path: /etc/user_script.sh
`
	if _, err := tmpFile.WriteString(rkeContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	driver.RKEProvisionUserData = tmpFile.Name()

	result, err := driver.GetFinalUserData()
	assert.NoError(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(expectedResult)), result)
}
