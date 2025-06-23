package pgsql_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunPgsqlInstanceCycle(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_instance." + rnd

	resourceFile := "resource_ctyun_pgsql_instance.tf"
	billMode := "month"
	hostType := "S7"
	prodVersion := "12.22"
	prodId := 10003011
	storageType := "SATA"
	StorageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	password := "VqOcfgJ6Nf2houSe5C9sxgM4ycExVK+F0bBZwBGdiy8DCVXoSyck0lPxw9XMRgHur2lQYenOJ5K/FxZ30qlwbKG3NfgNoPq+AXDeSDdycGTqa1TzLdGnYwAeC/hEa8pyUKS9LdlW7nnM1nGUvGCXkGdzJP8lbHCwonzazEnF3RI="
	caseCensitive := "0"
	nodeType := "master"
	instSpec := "1"
	prodPerformanceSpce := "2C4G"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	azInfo := `[{"availability_zone_name":"cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count":1, "node_type":"master"}]`
	osType := "11"
	cpuType := "30"
	period := fmt.Sprint(`cycle_count=1`)
	purchaseCount := fmt.Sprint(`purchase_count=1`)

	updatedProdID := 10003024
	updatedAzInfo := `[{"availability_zone_name":"cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count":2, "node_type":"master"}]`

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
			// 按包周期创建单节点，并测试单节点->1主2备扩容
			// Create
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, hostType, prodVersion, prodId, storageType, StorageSpace, name, password, caseCensitive, nodeType, instSpec,
					prodPerformanceSpce, vpcID, subnetID, securityGroupID, azInfo, "", false, false, false, osType, cpuType, period, purchaseCount),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", billMode),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", prodId)),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
				),
			},
			// 升配至1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, hostType, prodVersion, updatedProdID, storageType, StorageSpace, name, password, caseCensitive, nodeType, instSpec,
					prodPerformanceSpce, vpcID, subnetID, securityGroupID, updatedAzInfo, "", false, false, false, osType, cpuType, period, purchaseCount),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", billMode),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", updatedProdID)),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
				),
			},
			// destroy
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, hostType, prodVersion, updatedProdID, storageType, StorageSpace, name, password, caseCensitive, nodeType, instSpec,
					prodPerformanceSpce, vpcID, subnetID, securityGroupID, updatedAzInfo, "", false, false, false, osType, cpuType, period, purchaseCount),
				Destroy: true,
			},
		},
	})
}
