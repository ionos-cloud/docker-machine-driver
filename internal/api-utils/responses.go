package api_utils

import (
	"github.com/docker/machine/libmachine/log"
	api_utils "github.com/ionos-cloud/docker-machine-driver/pkg/api-utils"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func SanitizerAdapter(response ionoscloud.APIResponse) error {
	api_utils.SanitizeResponse(response, log.Infof)
}
