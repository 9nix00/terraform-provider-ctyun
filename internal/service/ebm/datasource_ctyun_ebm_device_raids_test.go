package ebm_test

import (
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEbmDeviceRaids(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
data "ctyun_ebm_device_raids" "test" {
  region_id = "200000001852"
  az_name = "cn-huabei2-tj1A-public-ctcloud"
  device_type = "physical.s5.2xlarge4"
  volume_type = "system"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_raids.test", "raids.#", "1"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_raids.test", "raids.0.uuid", "r-wtzluqacgzzxgunnabdkpnpjew3d"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_raids.test", "raids.0.volume_detail", "2*480GB(SSD)"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_raids.test", "raids.0.volume_type", "SYSTEM"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_raids.test", "raids.0.name_en", "RAID1"),
				),
			},
		},
	})
}
