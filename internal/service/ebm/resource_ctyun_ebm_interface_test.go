package ebm_test

import (
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEbmInterface(t *testing.T) {
	resourceName := "ctyun_ebm_interface.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "ctyun" {
  env                  = "prod"
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj-3a-public-ctcloud"
}

resource "ctyun_ebm_interface" "test" {
  security_group_ids = ["sg-t0ae11aig1"]
  instance_id = "ss-uadmwtxinfp4tkbhvwp52vnzl2kn"
  ipv4 = "192.168.0.18"
  subnet_id = "subnet-43z7cqmjlp"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", "sg-t0ae11aig1"),
					resource.TestCheckResourceAttrSet(resourceName, "interface_id"),
				),
			},
			{
				Config: `
provider "ctyun" {
  env                  = "prod"
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj-3a-public-ctcloud"
}

resource "ctyun_ebm_interface" "test" {
  security_group_ids = ["sg-t0ae11aig1", "sg-hsqwzeythj"]
  instance_id = "ss-uadmwtxinfp4tkbhvwp52vnzl2kn"
  ipv4 = "192.168.0.18"
  subnet_id = "subnet-43z7cqmjlp"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
