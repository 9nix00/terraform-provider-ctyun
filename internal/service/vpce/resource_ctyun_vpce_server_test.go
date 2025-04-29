package vpce_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunVpceServer(t *testing.T) {
	if !initMain {
		err := initSharedResources()
		t.Error(err)
		return
	}

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_vpce_server." + rnd
	datasourceName := "data.ctyun_vpce_servers." + dnd
	resourceFile := "resource_ctyun_vpce_server.tf"
	datasourceFile := "datasource_ctyun_vpce_servers.tf"

	initName := "init"
	initEndpointPort := "1"
	initWhitelistEmail := `lity9@chinatelecom.cn`
	updatedName := "updated"
	updatedEndpointPort := "2"
	updatedWhitelistEmail := `yunguan_ops@chinatelecom.cn`

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
			// Read testing
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, initEndpointPort, initWhitelistEmail, sharedVpcID, sharedSubnetID, sharedEcsID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "rules.0.endpoint_port", initEndpointPort),
					resource.TestCheckTypeSetElemAttr(resourceName, "whitelist_email.*", initWhitelistEmail),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedEndpointPort, updatedWhitelistEmail, sharedVpcID, sharedSubnetID, sharedEcsID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "rules.0.endpoint_port", updatedEndpointPort),
					resource.TestCheckTypeSetElemAttr(resourceName, "whitelist_email.*", updatedWhitelistEmail),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedEndpointPort, updatedWhitelistEmail, sharedVpcID, sharedSubnetID, sharedEcsID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "vpce_servers.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "vpce_servers.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "vpce_servers.0.rules.0.endpoint_port", updatedEndpointPort),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id = ds.ID
					regionId := ds.Attributes["region_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"subnet_id",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedEndpointPort, updatedWhitelistEmail, sharedVpcID, sharedSubnetID, sharedEcsID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	},
	)
}
