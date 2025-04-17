package nat_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"terraform-provider-ctyun/internal/service"
	"testing"
)

func TestAccNewCtyunNatResource(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	resourceName := "ctyun_nat.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "ctyun" {
					  region_id = "200000002530"
					  az_name = "az1"
					}

					resource "ctyun_nat" "test"{
						region_id = "200000002530"
						vpc_id = "vpc-wf029jgx2d"
						spec = 1
						name = "nat-terraform-test"
						description = "terraform测试"
						cycle_type = "on_demand"
						az_name = "cn-huanan2-1A-public-ctcloud"
					}			
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "master_resource_status", "started"),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"), //判断时候有该属性
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			//{
			//	Config: `
			//
			//     `
			//},
		},
	})
}
