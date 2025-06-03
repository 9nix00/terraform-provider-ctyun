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

func TestAccCtyunMysqlInstance1(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd

	resourceFile := "resource_ctyun_mysql_instance.tf"
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

	nodeType := "master"
	instSpec := "1"
	storageType := "SATA"
	storageSpace := 100
	prodPerformanceSpec := "2C4G"
	disks := 1
	updatedDiskAvailabilityZoneInfo := `availability_zone_info = [{"availability_zone_name":"cn-nm-het3-1a-public-ctcloud","availability_zone_count":2,"node_type":"slave"}]`

	cpuType := "30"
	osType := "11"
	// 单节点
	ProdId := 10001003
	// 单机到一主一备
	//updatedProdID := 10001001
	// 一主两备
	updatedDoubleProId := 10001002
	cycleBillMode := "1"
	NodeOneAvailabilityZoneInfo := `availability_zone_info = [{"availability_zone_name":"cn-nm-het3-1a-public-ctcloud","availability_zone_count":1,"node_type":"master"}]`
	//backupOneAvailabilityZoneInfo := `availability_zone_info=[{"availability_zone_name":"cn-nm-het3-1a-public-ctcloud","availability_zone_count":1,"node_type":"master"},{"availability_zone_name":"cn-nm-het3-1a-public-ctcloud","availability_zone_count":1,"node_type":"slave"}]`
	//ProdUpdatedMsyqlNodeInfoList := `{"node_type":"master","inst_spec":"1","storage_type":"SATA","storage_space":120,"prod_performance_spec":"2C8G","disks":1,"availability_zone_info":[{"availability_zone_name":"cn-nm-het3-1a-public-ctcloud","availability_zone_count":1,"node_type":"master"}]}`

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
			// 2 包周期创建，创建1主1备，升级为1主2备, 关机，开启，重启验证
			// 创建1主1备-》1主2备已经完成
			// 单节点-》1主2备,验证通过
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, ProdId, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, NodeOneAvailabilityZoneInfo, false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10001003)),
				),
			},
			// 升级1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, updatedDoubleProId, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10001002)),
				),
			},
			// 关机验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, updatedDoubleProId, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, false, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 6)),
				),
			},
			// 开机验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, updatedDoubleProId, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, true, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 0)),
				),
			},
			// 重启验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, prodVersion, vpcID, hostType, subnetID, securityGroupID, name, password, period, count, autoRenewStatus, updatedDoubleProId, cpuType, osType, "", nodeType, instSpec, storageType, storageSpace, prodPerformanceSpec, disks, updatedDiskAvailabilityZoneInfo, false, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 0)),
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
