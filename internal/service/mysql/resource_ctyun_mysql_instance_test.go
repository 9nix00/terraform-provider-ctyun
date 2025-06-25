package mysql_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunMysqlInstance(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd
	datasourceName := "data.ctyun_mysql_instances." + dnd

	resourceFile := "resource_ctyun_mysql_instance.tf"
	datasourceFile := "datasource_ctyun_mysql_instances.tf"

	cycleType := "on_demand"
	vpcID := dependence.vpcID
	hostType := "S7"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "terraform-provider-ctyun" + utils.GenerateRandomString()
	password := "kqjwyk111"
	cycleCount := "cycle_count=1"
	autoRenewStatus := `auto_renew=false`
	prodID := "Single57"

	instanceSeries := "S"
	storageType := "SATA"
	storageSpace := 100
	prodPerformanceSpec := "2C4G"
	availabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"}]`, dependence.azName)
	updatedDiskAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"slave"}]`, dependence.azName)

	cpuType := "Intel"
	osType := "ctyunos"

	updatedName := "terraform-provider-ctyun-new-" + utils.GenerateRandomString()
	updatedWritePort := `write_port=13306`

	// 磁盘、规格升配
	updatedStorageSpace := 120
	updatedBackupStorageSpace := `backup_storage_space=150`
	updatedProdPerformanceSpec := "2C8G"
	// 单机到一主一备
	updatedProdID := "MasterSlave57"
	// 一主两备
	updatedDoubleProId := "Master2Slave57"
	cycleBillMode := "month"
	backupOneAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info=[{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"},{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"slave"}]`, dependence.azName, dependence.azName)

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
			// 1. 按需验证，单节点创建，扩容至1主1备，修改端口，修改名称。
			// create 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", prodID, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single57"),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
				),
			},
			// update, 实例名称、写端口更新验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, updatedName, password, "", "", prodID, cpuType, osType, updatedWritePort, instanceSeries, storageType, storageSpace, prodPerformanceSpec, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single57"),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
					resource.TestCheckResourceAttr(resourceName, "write_port", "13306"),
				),
			},
			// 升配验证-升级磁盘空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, updatedName, password, "", "", prodID, cpuType, osType, updatedWritePort, instanceSeries, storageType, updatedStorageSpace, updatedProdPerformanceSpec, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "120"),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", "2C8G"),
				),
			},
			// 升配验证-升级备份磁盘空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, updatedName, password, "", "", prodID, cpuType, osType, updatedWritePort, instanceSeries, storageType, updatedStorageSpace, updatedProdPerformanceSpec, availabilityZoneInfo, "", updatedBackupStorageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", "150"),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", "2C8G"),
				),
			},
			// 升配验证-单机规格扩容->1主1备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, updatedName, password, "", "", updatedProdID, cpuType, osType, updatedWritePort, instanceSeries, storageType, updatedStorageSpace, updatedProdPerformanceSpec, updatedDiskAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave57"),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, updatedName, password, "", "", updatedProdID, cpuType, osType, updatedWritePort, instanceSeries, storageType, updatedStorageSpace, updatedProdPerformanceSpec, updatedDiskAvailabilityZoneInfo, "", "") +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("prod_inst_name=%s.name", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "mysql_instances.#", "1"),
				),
			},
			//销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, updatedName, password, "", "", updatedProdID, cpuType, osType, updatedWritePort, instanceSeries, storageType, updatedStorageSpace, updatedProdPerformanceSpec, updatedDiskAvailabilityZoneInfo, "", ""),
				Destroy: true,
			},
			// 2 包周期创建，创建1主1备，升级为1主2备,
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, cycleCount, autoRenewStatus, updatedProdID, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, backupOneAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave57"),
				),
			},
			// 升级1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, cycleCount, autoRenewStatus, updatedDoubleProId, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, updatedDiskAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, cycleCount, autoRenewStatus, updatedDoubleProId, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, updatedDiskAvailabilityZoneInfo, "", ""),
				Destroy: true,
			},
		},
	})
}
