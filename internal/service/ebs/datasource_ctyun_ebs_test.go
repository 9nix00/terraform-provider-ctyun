package ebs_test

import (
	"fmt"
	"strconv"
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNewCtyunEbs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "ctyun" {
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj-3a-public-ctcloud"
  env                  = "prod"
}

resource "ctyun_ebs" "test" {
  name       = "ebs-tf-test"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctyun_ebs.test", "name", "ebs-tf-test"),
					resource.TestCheckResourceAttr("ctyun_ebs.test", "size", "60"),
					resource.TestCheckResourceAttrSet("ctyun_ebs.test", "id"),
					resource.TestCheckResourceAttrSet("ctyun_ebs.test", "master_order_id"),
				),
			},
			{
				Config: `
provider "ctyun" {
  region_id            = "200000001852"
  az_name              = "cn-huabei2-tj-3a-public-ctcloud"
  env                  = "prod"
}

resource "ctyun_ebs" "test" {
  name       = "ebs-tf-test"
  mode       = "vbd"
  type       = "sata"
  size       = 100
  cycle_type = "on_demand"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctyun_ebs.test", "size", "100"),
				),
			},
			{
				Config: `
provider "ctyun" {
  env = "prod"
}

data "ctyun_ebs_volumes" "test" {
  region_id = "200000001852"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						ds := s.RootModule().Resources["data.ctyun_ebs_volumes.test"].Primary

						count, err := strconv.Atoi(ds.Attributes["volumes.#"])
						if err != nil || count == 0 {
							return fmt.Errorf("volumes 无效: %v", ds.Attributes)
						}

						for i := 0; i < count; i++ {
							if ds.Attributes[fmt.Sprintf("volumes.%d.name", i)] == "ebs-tf-test" {
								return nil
							}
						}
						return fmt.Errorf("未找到目标元素")
					},
				),
			},
		},
	})
}
