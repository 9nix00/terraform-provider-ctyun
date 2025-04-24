package nat_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunSNat(t *testing.T) {

	rnd1 := utils.GenerateRandomString()
	rnd2 := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName1 := "ctyun_nat_snat." + rnd1
	resourceName2 := "ctyun_nat_snat." + rnd2
	datasourceName := "data.ctyun_nat_snats." + dnd
	resourceFile1 := "resource_ctyun_nat_snat1.tf"
	resourceFile2 := "resource_ctyun_nat_snat2.tf"
	datasourceFile := "datasource_ctyun_snat.tf"

	initSourceCidr := "192.168.0.0/24"
	updatedSourceCidr := "192.168.128.0/25"
	sourceSubnetId := "subnet-ysrcfdvli9"
	updatedSubnetId := "subnet-syq0nr9yyq"

	natGateWayId := "natgw-asdsmh8scy"
	snatIps := "[eip-s7vhil3y30]"

	//updateDescription := utils.GenerateRandomString()
	//var id string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// 1.resource create验证1:
				// subnetType = 1(有vpcId 的子网情况),sourceSubnetId必传
				Config: utils.LoadTestCase(resourceFile1, rnd1, natGateWayId, sourceSubnetId, snatIps),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "nat_gateway_id", natGateWayId),
					resource.TestCheckResourceAttr(resourceName1, "source_subnet_id", sourceSubnetId),
					resource.TestCheckResourceAttr(resourceName1, "snat_ips", snatIps),
					resource.TestCheckResourceAttrSet(resourceName1, "snat_id"),
				),
			},
			{
				// 3. resource update source_subnet_id验证
				Config: utils.LoadTestCase(resourceFile1, rnd1, natGateWayId, updatedSubnetId, snatIps),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "source_subnet_id", updatedSubnetId),
				),
			},
			{
				// 2.resource create验证2:
				// subnetType = 0(自定义情况),sourceCIDR必传
				Config: utils.LoadTestCase(resourceFile2, rnd2, natGateWayId, initSourceCidr, snatIps),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName2, "nat_gateway_id", natGateWayId),
					resource.TestCheckResourceAttr(resourceName2, "source_cidr", initSourceCidr),
					resource.TestCheckResourceAttr(resourceName2, "snat_ips", snatIps),
					resource.TestCheckResourceAttrSet(resourceName2, "snat_id"),
				),
			},
			{
				// 4. resource update source_cidr验证
				Config: utils.LoadTestCase(resourceFile2, rnd2, natGateWayId, updatedSourceCidr, snatIps),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName2, "source_cidr", updatedSourceCidr),
				),
			},
			{
				// 5. datasource验证
				Config: utils.LoadTestCase(resourceFile1, rnd1, natGateWayId, updatedSubnetId, snatIps) +
					//utils.LoadTestCase(resourceFile2, rnd2, natGateWayId, updatedSourceCidr, snatIps) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "snats.#", "1"),
					//resource.TestCheckResourceAttr(datasourceName, "snats.0.subnet_id", updatedSubnetId),
				),
			},
			//{
			//	ResourceName: resourceName1,
			//	ImportState:  true,
			//	ImportStateIdFunc: func(s *terraform.State) (string, error) {
			//		ds := s.RootModule().Resources[resourceName1].Primary
			//		id := ds.ID
			//		regionId := ds.Attributes["region_id"]
			//		snat_id := ds.Attributes["snat_id"]
			//		if id == "" || snat_id == "" {
			//			return "", fmt.Errorf("ID/snat_id or regionID cannot be empty")
			//		}
			//		return fmt.Sprintf("%s,%s,%s", id, snat_id, regionId), nil
			//	},
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{},
			//},
		},
	})
}
