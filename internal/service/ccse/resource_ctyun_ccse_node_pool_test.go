package ccse_test

import (
	"fmt"
	"strconv"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunccse_node_pool(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_node_pool." + rnd
	datasourceName := "data.ctyun_ccse_node_pools." + dnd
	resourceFile := "resource_ctyun_ccse_node_pool.tf"
	datasourceFile := "datasource_ctyun_ccse_node_pools.tf"

	initName := "init-pool"
	initAutoRenewStatus := 0
	initVisibilityPostHostScript := "YWJj"
	initVisibilityHostScript := "MTIz"
	initSysDiskType := "SATA"
	initSysDiskSize := 100
	initDataDiskType := "SATA"
	initDataDiskSize := 200
	initCycleType := "on_demand"
	initCycleCount := ""

	updatedName := "updated-pool"
	updatedAutoRenewStatus := 1
	updatedVisibilityPostHostScript := "MTIz"
	updatedVisibilityHostScript := "YWJj"
	updatedSysDiskType := "SSD"
	updatedSysDiskSize := 200
	updatedDataDiskType := "SSD"
	updatedDataDiskSize := 400
	updatedCycleType := "month"
	updatedCycleCount := "cycle_count             = 1"

	var id string
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
				Config: utils.LoadTestCase(resourceFile, rnd,
					initName,
					initAutoRenewStatus,
					initVisibilityPostHostScript,
					initVisibilityHostScript,
					initSysDiskType,
					initSysDiskSize,
					initDataDiskType,
					initDataDiskSize,
					initCycleType,
					initCycleCount,
					dependence.flavorName,
					dependence.clusterID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_pool_name", initName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew_status", strconv.Itoa(initAutoRenewStatus)),
					resource.TestCheckResourceAttr(resourceName, "visibility_post_host_script", initVisibilityPostHostScript),
					resource.TestCheckResourceAttr(resourceName, "visibility_host_script", initVisibilityHostScript),
					resource.TestCheckResourceAttr(resourceName, "sys_disk.type", initSysDiskType),
					resource.TestCheckResourceAttr(resourceName, "sys_disk.size", strconv.Itoa(initSysDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.type", initDataDiskType),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.size", strconv.Itoa(initDataDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", initCycleType),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					updatedName,
					updatedAutoRenewStatus,
					updatedVisibilityPostHostScript,
					updatedVisibilityHostScript,
					updatedSysDiskType,
					updatedSysDiskSize,
					updatedDataDiskType,
					updatedDataDiskSize,
					updatedCycleType,
					updatedCycleCount,
					dependence.flavorName,
					dependence.clusterID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_pool_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew_status", strconv.Itoa(updatedAutoRenewStatus)),
					resource.TestCheckResourceAttr(resourceName, "visibility_post_host_script", updatedVisibilityPostHostScript),
					resource.TestCheckResourceAttr(resourceName, "visibility_host_script", updatedVisibilityHostScript),
					resource.TestCheckResourceAttr(resourceName, "sys_disk.type", updatedSysDiskType),
					resource.TestCheckResourceAttr(resourceName, "sys_disk.size", strconv.Itoa(updatedSysDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.type", updatedDataDiskType),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.size", strconv.Itoa(updatedDataDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", updatedCycleType),
					resource.TestCheckResourceAttr(resourceName, "cycle_count", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					updatedName,
					updatedAutoRenewStatus,
					updatedVisibilityPostHostScript,
					updatedVisibilityHostScript,
					updatedSysDiskType,
					updatedSysDiskSize,
					updatedDataDiskType,
					updatedDataDiskSize,
					updatedCycleType,
					updatedCycleCount,
					dependence.flavorName,
					dependence.clusterID,
				) + utils.LoadTestCase(datasourceFile, dnd, dependence.clusterID, resourceName+".node_pool_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "records.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "records.0.node_pool_name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "records.0.auto_renew_status", strconv.Itoa(updatedAutoRenewStatus)),
					resource.TestCheckResourceAttr(datasourceName, "records.0.visibility_post_host_script", updatedVisibilityPostHostScript),
					resource.TestCheckResourceAttr(datasourceName, "records.0.visibility_host_script", updatedVisibilityHostScript),
					resource.TestCheckResourceAttr(datasourceName, "records.0.sys_disk.type", updatedSysDiskType),
					resource.TestCheckResourceAttr(datasourceName, "records.0.sys_disk.size", strconv.Itoa(updatedSysDiskSize)),
					resource.TestCheckResourceAttr(datasourceName, "records.0.data_disks.0.type", updatedDataDiskType),
					resource.TestCheckResourceAttr(datasourceName, "records.0.data_disks.0.size", strconv.Itoa(updatedDataDiskSize)),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cycle_type", updatedCycleType),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cycle_count", "1"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id = ds.ID
					regionId := ds.Attributes["region_id"]
					clusterId := ds.Attributes["cluster_id"]
					return fmt.Sprintf("%s,%s,%s", id, clusterId, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					updatedName,
					updatedAutoRenewStatus,
					updatedVisibilityPostHostScript,
					updatedVisibilityHostScript,
					updatedSysDiskType,
					updatedSysDiskSize,
					updatedDataDiskType,
					updatedDataDiskSize,
					updatedCycleType,
					updatedCycleCount,
					dependence.flavorName,
					dependence.clusterID,
				) + utils.LoadTestCase(datasourceFile, dnd, dependence.clusterID, resourceName+".node_pool_name"),
				Destroy: true,
			},
		},
	})
}
