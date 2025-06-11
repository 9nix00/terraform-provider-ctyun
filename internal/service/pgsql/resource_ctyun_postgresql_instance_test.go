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

func TestAccCtyunPgsqlInstance(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_instance." + rnd
	datasourceName := "data.ctyun_mysql_instances." + dnd

	resourceFile := "resource_ctyun_pgsql_instance.tf"
	datasourceFile := "datasource_ctyun_pgsql_instances.tf"

	billMode := "2"
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
	VpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	azInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud", "availability_zone_count":1, "node_type":"master"}]`

	updatedName := "pgsql-new" + utils.GenerateRandomString()
	updatedSecurityGroupID := dependence.securityGroupID2
	updatedProdPerformanceSpce := "2C8G"
	updatedProdID := 10003012
	updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud", "availability_zone_count":1, "node_type":"slave"}]`
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
			// 1. 按需验证，单节点创建，扩容至1主1备，修改名称，修改安全组， 规格扩容。
			// create 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, hostType, prodVersion, prodId, storageType, StorageSpace, name, password, caseCensitive,
					nodeType, instSpec, prodPerformanceSpce, VpcID, subnetID, securityGroupID, azInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "bill_mode", billMode),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", prodId)),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType)),
			},
			// update验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, hostType, prodVersion, updatedProdID, storageType, StorageSpace, updatedName, password, caseCensitive,
					nodeType, instSpec, updatedProdPerformanceSpce, VpcID, subnetID, updatedSecurityGroupID, updatedAzInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "bill_mode", billMode),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", prodId)),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", updatedProdPerformanceSpce),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, hostType, prodVersion, prodId, storageType, StorageSpace, updatedName, password, caseCensitive,
					nodeType, instSpec, updatedProdPerformanceSpce, VpcID, subnetID, updatedSecurityGroupID, azInfo) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("prod_inst_name=%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.0.name", name),
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.0.bill_mode", billMode),
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.0.prod_id", fmt.Sprintf("%d", prodId)),
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.0.subnet_id", subnetID),
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.0.security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.0.host_type", hostType),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, hostType, prodVersion, prodId, storageType, StorageSpace, updatedName, password, caseCensitive,
					nodeType, instSpec, updatedProdPerformanceSpce, VpcID, subnetID, updatedSecurityGroupID, azInfo),
				Destroy: true,
			},
		},
	},
	)
}
