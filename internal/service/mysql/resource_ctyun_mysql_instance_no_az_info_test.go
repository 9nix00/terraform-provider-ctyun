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

func TestAccCtyunMysqlNoAzInfoInstance(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd

	resourceFile := "resource_ctyun_mysql_instance.tf"

	cycleType := "on_demand"
	vpcID := dependence.vpcID
	hostType := "S7"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "terraform-provider-ctyun" + utils.GenerateRandomString()
	name1 := "terraform-provider-ctyun" + utils.GenerateRandomString()
	name2 := "terraform-provider-ctyun" + utils.GenerateRandomString()
	password := "kqjwyk111*"
	prodID := "Single57"
	updateProdID := "Master2Slave57"
	MsProdID := "MasterSlave57"

	instanceSeries := "S"
	storageType := "SATA"
	storageSpace := 100
	updatedStorageSpace := 120
	prodPerformanceSpec := "2C4G"
	updatedProdPerformanceSpec := "2C8G"
	backupStorageSpace := `backup_storage_space=120`

	cpuType := "Intel"
	osType := "ctyunos"
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

			// 直接开通一个1主1备的mysql, 并进行变配磁盘和1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name1, password, "", "", MsProdID, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name1),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave57"),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name1, password, "", "", updateProdID, cpuType, osType, "", instanceSeries, storageType, updatedStorageSpace, prodPerformanceSpec, "", "", backupStorageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name1),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name1, password, "", "", updateProdID, cpuType, osType, "", instanceSeries, storageType, updatedStorageSpace, prodPerformanceSpec, "", "", backupStorageSpace),
				Destroy: true,
			},

			// 1. 按需验证，单节点创建，扩容至1主1备，修改端口，修改名称。
			// create 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", prodID, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, "", "", ""),
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
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", updateProdID, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", updateProdID, cpuType, osType, "", instanceSeries, storageType, storageSpace, updatedProdPerformanceSpec, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", updateProdID, cpuType, osType, "", instanceSeries, storageType, storageSpace, updatedProdPerformanceSpec, "", "", ""),
				Destroy: true,
			},

			// 直接开通一个1主2备的mysql，变配规格
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name2, password, "", "", updateProdID, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
				),
			},
			// 变配规格2c4g->2c8g
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name2, password, "", "", updateProdID, cpuType, osType, "", instanceSeries, storageType, storageSpace, updatedProdPerformanceSpec, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "host_type", hostType),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
					resource.TestCheckResourceAttr(resourceName, "cpu_type", cpuType),
					resource.TestCheckResourceAttr(resourceName, "os_type", osType),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, hostType, subnetID, securityGroupID, name2, password, "", "", updateProdID, cpuType, osType, "", instanceSeries, storageType, storageSpace, updatedProdPerformanceSpec, "", "", ""),
				Destroy: true,
			},
		},
	})
}
