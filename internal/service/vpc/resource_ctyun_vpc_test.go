package vpc_test

import (
	"fmt"
	"terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpc(t *testing.T) {
	resourceName := "ctyun_vpc.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `
resource "ctyun_vpc" "test" {
  name        = "vpc-test-tf"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
  region_id   = "200000001852"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vpc-test-tf"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: `
resource "ctyun_vpc" "test" {
  name        = "vpc-test-tf-1"
  cidr        = "192.168.0.0/16"
  description = "terraform"
  enable_ipv6 = true
  region_id   = "200000001852"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vpc-test-tf-1"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform"),
				),
			},
			{
				Config: `
resource "ctyun_vpc" "test" {
  name        = "vpc-test-tf-1"
  cidr        = "192.168.0.0/16"
  description = "terraform"
  enable_ipv6 = true
  region_id   = "200000001852"
}

data "ctyun_vpcs" "test1" {
  region_id = "200000001852"
  vpc_id = ctyun_vpc.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ctyun_vpcs.test1", "vpcs.#", "1"),
					resource.TestCheckResourceAttr("data.ctyun_vpcs.test1", "vpcs.0.name", "vpc-test-tf-1"),
					resource.TestCheckResourceAttr("data.ctyun_vpcs.test1", "vpcs.0.cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr("data.ctyun_vpcs.test1", "vpcs.0.enable_ipv6", "true"),
				),
			},
			{

				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					projectId := ds.Attributes["project_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s,%s", id, regionId, projectId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
