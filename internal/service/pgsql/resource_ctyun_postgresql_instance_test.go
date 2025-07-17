package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
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
	datasourceName := "data.ctyun_postgresql_instances." + dnd

	resourceFile := "resource_ctyun_pgsql_instance.tf"
	datasourceFile := "datasource_ctyun_pgsql_instances.tf"

	cycleType := "on_demand"
	hostType := "S7"
	prodId := "Single1222"
	storageType := "SATA"
	backupStorageType := `backup_storage_type="SATA"`
	StorageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	//password := "VqOcfgJ6Nf2houSe5C9sxgM4ycExVK+F0bBZwBGdiy8DCVXoSyck0lPxw9XMRgHur2lQYenOJ5K/FxZ30qlwbKG3NfgNoPq+AXDeSDdycGTqa1TzLdGnYwAeC/hEa8pyUKS9LdlW7nnM1nGUvGCXkGdzJP8lbHCwonzazEnF3RI="
	password := "Kqjwyk123="
	caseCensitive := true
	instanceSeries := "S"
	prodPerformanceSpce := "2C4G"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	azInfo := `[{"availability_zone_name":"cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count":1, "node_type":"master"}]`
	osType := "ctyunos"
	cpuType := "Intel"

	updatedName := "pgsql-new" + utils.GenerateRandomString()
	updatedSecurityGroupID := dependence.securityGroupID2
	updatedProdPerformanceSpce := "2C8G"
	updatedProdID := "MasterSlave1222"
	updatedStorageSpace := 120
	updatedAzInfo := `[{"availability_zone_name":"cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count":1, "node_type":"slave"}]`
	updatedBackupStorageSpace := fmt.Sprintf(`backup_storage_space="%d"`, updatedStorageSpace)

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
			// 1. 按需验证，单节点创建，扩容至1主1备，修改名称，修改安全组， 规格扩容,磁盘扩容。
			// create 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, prodId, storageType, StorageSpace, name, password, caseCensitive,
					instanceSeries, prodPerformanceSpce, vpcID, subnetID, securityGroupID, azInfo, `backup_storage_space=100`, "", osType, cpuType, "", backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType)),
			},
			// update验证--姓名, 安全组，规格扩容
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, prodId, storageType, StorageSpace, updatedName, password, caseCensitive,
					instanceSeries, updatedProdPerformanceSpce, vpcID, subnetID, updatedSecurityGroupID, azInfo, "", "", osType, cpuType, "", backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", updatedProdPerformanceSpce),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", StorageSpace)),
				),
			},
			// update验证--backup磁盘
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, prodId, storageType, StorageSpace, updatedName, password, caseCensitive,
					instanceSeries, updatedProdPerformanceSpce, vpcID, subnetID, updatedSecurityGroupID, azInfo, updatedBackupStorageSpace, "", osType, cpuType, "", backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", updatedProdPerformanceSpce),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", StorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
				),
			},
			// update验证--master磁盘
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, prodId, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					instanceSeries, updatedProdPerformanceSpce, vpcID, subnetID, updatedSecurityGroupID, azInfo, "", "", osType, cpuType, "", backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", updatedProdPerformanceSpce),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
				),
			},
			// update验证--主备，关机，开机，重启
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					instanceSeries, updatedProdPerformanceSpce, vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, "", `running_control="stop"`, osType, cpuType, "", backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222"),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", updatedProdPerformanceSpce),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					instanceSeries, updatedProdPerformanceSpce, vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, "", `running_control="start"`, osType, cpuType, "", backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222"),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", updatedProdPerformanceSpce),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					instanceSeries, updatedProdPerformanceSpce, vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, "", `running_control="restart"`, osType, cpuType, "", backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222"),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", updatedProdPerformanceSpce),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					instanceSeries, updatedProdPerformanceSpce, vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, "", "", osType, cpuType, "", backupStorageType) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("prod_inst_id=%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.0.name", updatedName),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, hostType, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					instanceSeries, updatedProdPerformanceSpce, vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, "", "", osType, cpuType, "", backupStorageType),
				Destroy: true,
			},
		},
	},
	)
}
