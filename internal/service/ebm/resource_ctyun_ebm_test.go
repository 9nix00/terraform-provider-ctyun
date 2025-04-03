package ebm_test

import (
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNewCtyunEbm(t *testing.T) {
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
					resource.TestCheckResourceAttr("ctyun_ebm.test", "status", "running"),
					resource.TestCheckResourceAttrSet("ctyun_ebm.test", "id"),
					resource.TestCheckResourceAttrSet("ctyun_ebm.test", "master_order_id"),
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
					resource.TestCheckResourceAttr("ctyun_ebm.test", "instance_name", "ebm-tf-test-0402-1"),
					resource.TestCheckResourceAttr("ctyun_ebm.test", "hostname", "ebm-tf-test-0402-2"),
				),
			},

			{
				ResourceName:      "ctyun_ebm.test",
				ImportState:       true,
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
