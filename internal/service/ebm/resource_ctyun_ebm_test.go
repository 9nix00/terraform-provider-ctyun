package ebm_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunEbm(t *testing.T) {
	resourceName := "ctyun_ebm.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "ctyun" {
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj1A-public-ctcloud"
  env                  = "prod"
}

resource "ctyun_ebm" "test" {
  device_type = "physical.s5.2xlarge1"
  instance_name = "ebm-0323-tf"
  hostname = "ebm-03221-tf"
  image_uuid = "im-xevpi6apqilz1bixmogofyref9qm"
  password = "P@ss132345"
  security_group_ids = ["sg-hsqwzeythj","sg-t0ae11aig1"]
  vpc_id = "vpc-6zxqwrg1r6"
  ext_ip = "not_use"
  system_volume_raid_uuid = ""
  instance_charge_type = "order_on_demand"
  status = "running"
  disk_list = [{
    disk_type = "system"
    size = "100"
    type = "sata"
  }]
  network_card_list = [{
    master = true,
    subnet_id = "subnet-43z7cqmjlp"
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config: `
provider "ctyun" {
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj1A-public-ctcloud"
  env                  = "prod"
}

resource "ctyun_ebm" "test" {
  device_type = "physical.s5.2xlarge1"
  instance_name = "ebm-tf-test-0402-1"
  hostname = "ebm-tf-test-0402-2"
  image_uuid = "im-xevpi6apqilz1bixmogofyref9qm"
  password = "P@ss12345"
  security_group_ids = ["sg-hsqwzeythj","sg-t0ae11aig1"]
  vpc_id = "vpc-6zxqwrg1r6"
  ext_ip = "not_use"
  system_volume_raid_uuid = ""
  instance_charge_type = "order_on_demand"
  status = "stopped"
  disk_list = [{
    disk_type = "system"
    size = "100"
    type = "sata"
  }]
  network_card_list = [{
    master = true,
    subnet_id = "subnet-43z7cqmjlp"
  }]
}

data "ctyun_ebms" "data_test" {
	instance_id_list = ctyun_ebm.test.instance_id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ctyun_ebms.data_test", "instances.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_name", "ebm-tf-test-0402-1"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "ebm-tf-test-0402-2"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionID := ds.Attributes["region_id"]
					azName := ds.Attributes["az_name"]
					return fmt.Sprintf("%s,%s,%s", id, regionID, azName), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auto_renew_status", // 查询接口没返回
					"band_width",        // 查询接口没返回
					"ext_ip",            // 查询接口没返回
					"ip_type",           // 查询接口没返回
					"master_order_id",   // 查询接口没返回
					"password",          // 查询接口没返回
					"project_id",        // 查询接口没返回
					"user_data",
				},
			},
		},
	})
}
