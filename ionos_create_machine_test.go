package ionoscloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/ionos-cloud/docker-machine-driver/internal/utils"
	sdkgo "github.com/ionos-cloud/sdk-go/v6"
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
	driver.AppendRKECloudInit = true
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
	driver.AppendRKECloudInit = true

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

func TestGetParentLocation(t *testing.T) {
	tests := []struct {
		name     string
		location string
		want     string
	}{
		{"top-level location has no parent", "de/fra", ""},
		{"child location returns its parent", "de/fra/2", "de/fra"},
		{"deeper segments still take first two", "de/fra/2/extra", "de/fra"},
		{"empty location returns empty", "", ""},
		{"malformed (single segment) returns empty", "de", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &Driver{Location: tc.location}
			assert.Equal(t, tc.want, d.getParentLocation())
		})
	}
}

// When neither the imageId lookup nor the extended image-name search resolves
// anything (e.g. 403 from both endpoints in a child location with no contract
// visibility into the image catalog), the driver must fall back to passing the
// user-supplied name as an alias so the API can resolve it server-side.
func TestGetImageIdOrAlias_FallsBackToAliasWhenAllLookupsFail(t *testing.T) {
	driver, clientMock := NewTestDriver(t, "test-host", "defaultstore")
	driver.Location = "de/fra/2"

	clientMock.EXPECT().GetImageById("ubuntu:latest").Return(nil, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&sdkgo.Images{Items: &[]sdkgo.Image{}}, nil)

	got, err := driver.getImageIdOrAlias("ubuntu:latest")
	assert.NoError(t, err)
	assert.Equal(t, "ubuntu:latest", got)
	assert.True(t, driver.UseAlias)
}

// GetImages failing (any error shape — we don't restrict on the message) must
// not propagate: the alias fallback is what makes child locations work even
// when the contract has no visibility into the image catalog.
func TestGetImageIdOrAlias_FallsBackWhenGetImagesErrors(t *testing.T) {
	driver, clientMock := NewTestDriver(t, "test-host", "defaultstore")
	driver.Location = "de/fra/2"

	clientMock.EXPECT().GetImageById("ubuntu:latest").Return(nil, fmt.Errorf("403 forbidden"))
	clientMock.EXPECT().GetImages().Return(nil, fmt.Errorf("403 forbidden"))

	got, err := driver.getImageIdOrAlias("ubuntu:latest")
	assert.NoError(t, err)
	assert.Equal(t, "ubuntu:latest", got)
	assert.True(t, driver.UseAlias)
}
