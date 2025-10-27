package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunDhcpOptionSetAssociationVpc_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_dhcpoptionset_association_vpc." + rnd
	resourceFile := "resource_ctyun_dhcpoptionset_association_vpc.tf"
	dataSourceName := "data.ctyun_dhcpoptionset_association_vpcs.test"
	dataSourceFile := "datasource_ctyun_dhcpoptionset_association_vpcs.tf"

	// 测试参数
	dhcpOptionSetsId := "dopt-i543qdzbw0"
	vpcIds := fmt.Sprintf(`"%s"`, "vpc-ff96ycah87")
	updatedVpcIds := fmt.Sprintf(`"%s","%s"`, "vpc-ff96ycah87", "vpc-cgny4bplv8")

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// 测试创建DHCP选项集与VPC绑定关系
				Config: utils.LoadTestCase(resourceFile, rnd, dhcpOptionSetsId, vpcIds),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_sets_id", dhcpOptionSetsId),
					resource.TestCheckResourceAttr(resourceName, "vpc_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// 测试更新DHCP选项集与VPC绑定关系
				Config: utils.LoadTestCase(resourceFile, rnd, dhcpOptionSetsId, updatedVpcIds),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_sets_id", dhcpOptionSetsId),
					//resource.TestCheckResourceAttr(resourceName, "vpc_ids.#", "2"),
				),
			},
			//{
			//	ResourceName: resourceName,
			//	ImportState:  true,
			//	ImportStateIdFunc: func(s *terraform.State) (string, error) {
			//		ds := s.RootModule().Resources[resourceName].Primary
			//		id := ds.ID
			//		if id == "" {
			//			return "", fmt.Errorf("id is required")
			//		}
			//		return id, nil
			//	},
			//	ImportStateVerify: true,
			//	ImportStateVerifyIgnore: []string{
			//		"region_id",
			//	},
			//},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dhcpOptionSetsId, vpcIds),
				Destroy: true,
			},
			{
				// 测试数据源
				Config: utils.LoadTestCase(dataSourceFile, dhcpOptionSetsId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "vpcs.#"),
				),
			},
		},
	})
}
