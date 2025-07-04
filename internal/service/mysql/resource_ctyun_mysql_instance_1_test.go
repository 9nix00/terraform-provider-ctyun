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
	vpcID := dependence.vpcID
	hostType := "S7"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "terraform-provider-ctyun" + utils.GenerateRandomString()
	password := "kqjwyk123."
	//period := 1
	//autoRenewStatus := 0

	instanceSeries := "S"
	storageType := "SATA"
	storageSpace := 100
	prodPerformanceSpec := "2C4G"
	updatedDiskAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":2,"node_type":"slave"}]`, dependence.azName)

	cpuType := "Intel"
	osType := "centos"
	// 单节点
	ProdId := "Single57"
	// 单机到一主一备
	//updatedProdID := 10001001
	// 一主两备
	updatedDoubleProId := "Master2Slave57"
	cycleBillMode := "on_demand"
	NodeOneAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"}]`, dependence.azName)

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
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", ProdId, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, NodeOneAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single57"),
				),
			},
			// 升级1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", updatedDoubleProId, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, updatedDiskAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// 关机验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", updatedDoubleProId, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, updatedDiskAvailabilityZoneInfo, `running_control="freeze"`, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 6)),
				),
			},
			// 开机验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", updatedDoubleProId, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, updatedDiskAvailabilityZoneInfo, `running_control="unfreeze"`, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 0)),
				),
			},
			// 重启验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", updatedDoubleProId, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, updatedDiskAvailabilityZoneInfo, `running_control="restart"`, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 0)),
				),
			},
			// 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, hostType, subnetID, securityGroupID, name, password, "", "", updatedDoubleProId, cpuType, osType, "", instanceSeries, storageType, storageSpace, prodPerformanceSpec, updatedDiskAvailabilityZoneInfo, "", ""),
				Destroy: true,
			},
		},
	})
}
