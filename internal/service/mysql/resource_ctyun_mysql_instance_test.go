package mysql_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunMysqlInstance(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd
	datasourceName := "data.ctyun_mysql_instances." + dnd

	resourceFile := "resource_ctyun_mysql_instance.tf"
	datasourceFile := "datasource_ctyun_mysql_instances.tf"

	billMode := "2"
	prodVersion := "5.7"
	vpcID := dependence.vpcID
	hostType := "S7"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "terraform-provider-ctyun" + utils.GenerateRandomString()
	password := "kqjwyk"
	period := 1
	count := 1
	autoRenewStatus := 0
	prodID := 10001003

	nodeType := "master"
	instSpec := "1"
	storageType := "SATA"
	storageSpace := 100
	prodPerformanceSpec := "2C4G"
	disks := 1
	availabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"}]`, dependence.azName)
	updatedDiskAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"slave"}]`, dependence.azName)

	cpuType := "30"
	osType := "11"

	updatedName := "terraform-provider-ctyun-new-" + utils.GenerateRandomString()
	updatedWritePort := `write_port="13306"`

	// 磁盘、规格升配
	updatedStorageSpace := 120
	updatedProdPerformanceSpec := "2C8G"
	// 单机到一主一备
	updatedProdID := 10001001
	// 一主两备
	updatedDoubleProId := 10001002
	cycleBillMode := "1"
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
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, prodID, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, availabilityZoneInfo, false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "bill_mode", billMode),
					resource.TestCheckResourceAttr(resourceName, "prod_version", prodVersion),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "auto_renew_status", strconv.Itoa(autoRenewStatus)),
					resource.TestCheckResourceAttr(resourceName, "prod_id", strconv.Itoa(prodID)),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
				),
			},
			// update, 实例名称、写端口更新验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, updatedName, password, period, count, autoRenewStatus, prodID, cpuType, osType, updatedWritePort, nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, availabilityZoneInfo, false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "bill_mode", billMode),
					resource.TestCheckResourceAttr(resourceName, "prod_version", prodVersion),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "auto_renew_status", strconv.Itoa(autoRenewStatus)),
					resource.TestCheckResourceAttr(resourceName, "prod_id", strconv.Itoa(prodID)),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
					resource.TestCheckResourceAttr(resourceName, "write_port", "13306"),
				),
			},
			// 升配验证-升级磁盘空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, updatedName, password, period, count, autoRenewStatus, prodID, cpuType, osType, updatedWritePort, nodeType, instSpec, storageType, updatedStorageSpace, updatedProdPerformanceSpec, disks, availabilityZoneInfo, false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "120"),
					resource.TestCheckResourceAttr(resourceName, "prod_performance_spec", "2C8G"),
				),
			},
			// 升配验证-单机规格扩容->1主1备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, updatedName, password, period, count, autoRenewStatus, updatedProdID, cpuType, osType, updatedWritePort, nodeType, instSpec, storageType, updatedStorageSpace, updatedProdPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", updatedProdID)),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, billMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, updatedName, password, period, count, autoRenewStatus, updatedProdID, cpuType, osType, updatedWritePort, nodeType, instSpec, storageType, updatedStorageSpace, updatedProdPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, false, false, false) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("prod_inst_name=%s.name", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "mysql_instances.#", "1"),
				),
			},
			//销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, billMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, updatedName, password, period, count, autoRenewStatus, updatedProdID, cpuType, osType, updatedWritePort, nodeType, instSpec, storageType, updatedStorageSpace, updatedProdPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, false, false, false),
				Destroy: true,
			},
			// 2 包周期创建，创建1主1备，升级为1主2备,
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, updatedProdID, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, backupOneAvailabilityZoneInfo, false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", updatedProdID)),
				),
			},
			// 升级1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, updatedDoubleProId, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", updatedDoubleProId)),
				),
			},
			// 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleBillMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, updatedDoubleProId, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, false, false, false),
				Destroy: true,
			},
		},
	})
}
