package ebm_test

import (
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEbmDeviceTypes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "ctyun" {
  env = "prod"
}

data "ctyun_ebm_device_types" "test" {
  region_id = "200000001852"
  az_name = "cn-huabei2-tj-3a-public-ctcloud"
}`,
			},
		},
	})
}
