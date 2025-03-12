package ebm_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/provider"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNewCtyunEbmDeviceTypes(t *testing.T) {
	fmt.Println(123)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: provider.TestConfig + `data "ctyun_ebm_device_types" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_types.test", "device_types.#", "2"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_types.test", "device_types.0.az_name", "cn-huadong1-jsnj1A-public-ctcloud"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_types.test", "device_types.0.device_type", "physical.h7ns.4xlarge25"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_types.test", "device_types.0.cpu_manufacturer", "Intel"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_types.test", "device_types.0.name_en", "Inspur NF5688M6&X660 G45"),
					resource.TestCheckResourceAttr("data.ctyun_ebm_device_types.test", "device_types.0.gpu_manufacturer", "NVIDIA"),
				),
			},
		},
	})
}
