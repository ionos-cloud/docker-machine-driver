package ionoscloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var dcId = ""

const (
	testStoreDir    = ".store-test"
	machineTestName = "docker machine unit tests public 8"
)

type DriverOptionsMock struct {
	Data map[string]interface{}
}

func (d DriverOptionsMock) String(key string) string {
	return d.Data[key].(string)
}

func (d DriverOptionsMock) StringSlice(key string) []string {
	return d.Data[key].([]string)
}

func (d DriverOptionsMock) Int(key string) int {
	return d.Data[key].(int)
}

func (d DriverOptionsMock) Bool(key string) bool {
	return d.Data[key].(bool)
}

func cleanup() error {
	return os.RemoveAll(testStoreDir)
}

func getTestStorePath() (string, error) {
	tmpDir, err := ioutil.TempDir("", "machine-test-")
	if err != nil {
		return "", err
	}
	os.Setenv("MACHINE_STORAGE_PATH", tmpDir)
	return tmpDir, nil
}

func getDefaultTestDriverFlags() *DriverOptionsMock {
	return &DriverOptionsMock{
		Data: map[string]interface{}{
			"profitbricks-endpoint":                 "https://api.profitbricks.com/cloudapi/v4",
			"profitbricks-username":                 os.Getenv("PROFITBRICKS_USERNAME"),
			"profitbricks-password":                 os.Getenv("PROFITBRICKS_PASSWORD"),
			"profitbricks-disk-type":                "HDD",
			"profitbricks-disk-size":                5,
			"profitbricks-cpu-family":               "AMD_OPTERON",
			"profitbricks-image":                    "Ubuntu-16.04",
			"profitbricks-cores":                    1,
			"profitbricks-ram":                      1024,
			"profitbricks-location":                 "us/las",
			"profitbricks-datacenter-id":            "",
			"profitbricks-volume-availability-zone": "AUTO",
			"profitbricks-server-availability-zone": "AUTO",
			"profitbricks-ssh-key":                  ``,
			"swarm-master":                          true,
			"swarm-host":                            "2",
			"swarm-discovery":                       "3",
		},
	}
}

func getTestDriver() (*Driver, error) {
	storePath, err := getTestStorePath()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	d := NewDriver(machineTestName, storePath)

	/*if err != nil {
		return nil, err
	}*/
	d.SetConfigFromFlags(getDefaultTestDriverFlags())
	drv := d.(*Driver)
	return drv, nil
}

func TestCreate(t *testing.T) {
	d, _ := getTestDriver()

	createerr := d.Create()
	if createerr != nil {
		t.Error(createerr)
	}

	state, err := d.GetState()
	dcId = d.DatacenterId

	fmt.Println(state)
	if err != nil {
		t.Error(err)
	}
}

func TestGetMachineName(t *testing.T) {
	d, _ := getTestDriver()
	if d.MachineName == "" {
		t.Fatal("Machine name not suplied.")
	}
	fmt.Println(d.GetMachineName())
}

func TestKill(t *testing.T) {
	d, _ := getTestDriver()
	d.Kill()
}

func TestRemove(t *testing.T) {
	d, _ := getTestDriver()
	d.DatacenterId = dcId
	d.Remove()
}

func TestGetImageName(t *testing.T) {
	d, _ := getTestDriver()
	res := d.getImageId("Debian-8-server1")

	fmt.Println(res == "")
}
