package ccse_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunClusterStandard(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_cluster." + rnd
	datasourceName := "data.ctyun_ccse_clusters." + dnd
	resourceFile := "resource_ctyun_ccse_cluster_standard.tf"
	datasourceFile := "datasource_ctyun_ccse_clusters.tf"

	clusterName := "tf-" + utils.GenerateRandomString()
	clusterSeries := "cce.standard"

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
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_series", clusterSeries),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "records.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_series", clusterSeries),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunClusterManaged(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_cluster." + rnd
	datasourceName := "data.ctyun_ccse_clusters." + dnd
	resourceFile := "resource_ctyun_ccse_cluster_managed.tf"
	datasourceFile := "datasource_ctyun_ccse_clusters.tf"

	clusterName := "tf-" + rnd
	clusterSeries := "cce.managed"

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
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_series", clusterSeries),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "records.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_series", clusterSeries),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
				Destroy: true,
			},
		},
	})
}
