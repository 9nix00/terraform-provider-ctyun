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
  region_id	= "bb9fdb42056f11eda1610242ac110002"
  az_name = "cn-huadong1-jsnj2A-public-ctcloud"
  env = "prod"
}

resource "ctyun_ebm" "test" {
  device_type = "physical.s5.2xlarge4"
  instance_name = "ebm-25-tf"
  hostname = "ebm-25-tf"
  image_uuid = "im-xevpi6apqilz1bixmogofyref9qm"
  password = "P@ss12345"
  security_group_id = "sg-vrp4x1lm7p"
  vpc_id = "vpc-5o8oe0oci6"
  ext_ip = "not_use"
  system_volume_raid_uuid = "r-wtzluqacgzzxgunnabdkpnpjew3d"
  data_volume_raid_uuid = "r-qytwf9r5h0yn9x4evjkyr0n1cwyb"
  instance_charge_type = "order_on_demand"
  network_card_list = [{
    master = true,
    subnet_id = "subnet-n7zbsy4b91"
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctyun_ebm.test", "device_type", "physical.s5.2xlarge4"),
					resource.TestCheckResourceAttr("ctyun_ebm.test", "status", "running"),
					resource.TestCheckResourceAttrSet("ctyun_ebm.test", "id"),
					resource.TestCheckResourceAttrSet("ctyun_ebm.test", "master_order_id"),
				),
			},
			{
				Config: `
provider "ctyun" {
  region_id	= "bb9fdb42056f11eda1610242ac110002"
  az_name = "cn-huadong1-jsnj2A-public-ctcloud"
  env = "prod"
}

resource "ctyun_ebm" "test" {
  device_type = "physical.s5.2xlarge4"
  instance_name = "ebm-0324-tf"
  hostname = "ebm-0324-hostname"
  image_uuid = "im-xevpi6apqilz1bixmogofyref9qm"
  password = "P@ss12345"
  security_group_id = "sg-vrp4x1lm7p"
  vpc_id = "vpc-5o8oe0oci6"
  ext_ip = "not_use"
  status = "stopped"
  system_volume_raid_uuid = "r-wtzluqacgzzxgunnabdkpnpjew3d"
  data_volume_raid_uuid = "r-qytwf9r5h0yn9x4evjkyr0n1cwyb"
  instance_charge_type = "order_on_demand"
  network_card_list = [{
    master = true,
    subnet_id = "subnet-n7zbsy4b91"
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctyun_ebm.test", "device_type", "physical.s5.2xlarge4"),
					resource.TestCheckResourceAttr("ctyun_ebm.test", "status", "stopped"),
					resource.TestCheckResourceAttr("ctyun_ebm.test", "instance_name", "ebm-0324-tf"),
					resource.TestCheckResourceAttr("ctyun_ebm.test", "hostname", "ebm-0324-hostname"),
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
