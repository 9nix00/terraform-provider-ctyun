package elb_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"terraform-provider-ctyun/internal/service"
	"testing"
)

func TestAccNewCtyunElbLoadBalancersDataSource(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					 provider "ctyun" {
					   region_id = "bb9fdb42056f11eda1610242ac110002"
					 }
					
					 data "ctyun_elb_loadbalancers" "test"{
						 region_id = "bb9fdb42056f11eda1610242ac110002"
					 }
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ctyun_elb_loadbalancers.test", "region_id", "200000002530"),
				),
			},
		},
	})
}
