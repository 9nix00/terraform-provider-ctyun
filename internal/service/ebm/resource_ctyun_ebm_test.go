package ebm_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/provider"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNewCtyunEbm(t *testing.T) {
	fmt.Println(234)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: provider.TestConfig + `
resource "ctyun_ebm" "test" {
  device_type = "physical.s5.2xlarge4"
  instance_name = "ebm-0312-tf"
  hostname = "ebm-0310-tf"
  image_uuid = "im-xevpi6apqilz1bixmogofyref9qm"
  password = "P@ss12345"
  security_group_id = "sg-vrp4x1lm7p"
  vpc_id = "vpc-5o8oe0oci6"
  ext_ip = "0"
  system_volume_raid_uuid = "r-wtzluqacgzzxgunnabdkpnpjew3d"
  data_volume_raid_uuid = "r-qytwf9r5h0yn9x4evjkyr0n1cwyb"
  instance_charge_type = "ORDER_ON_DEMAND"
  status = "RUNNING"
  network_card_list = [{
    master = true,
    subnet_id = "subnet-n7zbsy4b91"
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("ctyun_ebm.test", "items.#", "1"),
					// Verify first order item
					resource.TestCheckResourceAttr("ctyun_ebm.test", "items.0.device_type", "physical.s5.2xlarge4"),
					resource.TestCheckResourceAttr("ctyun_ebm.test", "items.0.status", "RUNNING"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ctyun_ebm.test", "id"),
					resource.TestCheckResourceAttrSet("ctyun_ebm.test", "master_order_id"),
				),
			},
			// Update and Read testing
			{
				Config: provider.TestConfig + `
resource "ctyun_ebm" "test" {
  device_type = "physical.s5.2xlarge4"
  instance_name = "ebm-0312-tf"
  hostname = "ebm-0310-tf"
  image_uuid = "im-xevpi6apqilz1bixmogofyref9qm"
  password = "P@ss12345"
  security_group_id = "sg-vrp4x1lm7p"
  vpc_id = "vpc-5o8oe0oci6"
  ext_ip = "0"
  system_volume_raid_uuid = "r-wtzluqacgzzxgunnabdkpnpjew3d"
  data_volume_raid_uuid = "r-qytwf9r5h0yn9x4evjkyr0n1cwyb"
  instance_charge_type = "ORDER_ON_DEMAND"
  status = "STOPPED"
  network_card_list = [{
    master = true,
    subnet_id = "subnet-n7zbsy4b91"
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("ctyun_ebm.test", "items.#", "1"),
					// Verify first order item
					resource.TestCheckResourceAttr("ctyun_ebm.test", "items.0.device_type", "physical.s5.2xlarge4"),
					resource.TestCheckResourceAttr("ctyun_ebm.test", "items.0.status", "STOPPED"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ctyun_ebm.test", "id"),
					resource.TestCheckResourceAttrSet("ctyun_ebm.test", "master_order_id"),
				),
			},

			// ImportState testing
			//{
			//	ResourceName:      "ctyun_ebm.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateVerifyIgnore: []string{"last_updated"},
			//},

			// Delete testing automatically occurs in TestCase
		},
	})
}
