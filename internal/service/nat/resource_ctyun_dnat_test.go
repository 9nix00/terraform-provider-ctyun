package nat_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunDNat(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_nat_dnat." + rnd
	datasourceName := "data.ctyun_nat_dnats." + dnd

	resourceFile := "resource_ctyun_nat_dnat.tf"
	datasourceFile := "datasource_ctyun_nat_dnat.tf"

	natGatewayId := dependence.natID
	externalId := dependence.eipID
	virtualMachineType := 2
	internalPort := utils.GenerateRandomPort(0, 65535)
	updatedInternalPort := utils.GenerateRandomPort(0, 65535)
	externalPort := utils.GenerateRandomPort(0, 1024)
	updatedExternalPort := utils.GenerateRandomPort(0, 1024)

	internalIp := "127.0.0.1"
	updatedInternalIp := "127.0.0.2"

	protocol := "tcp"
	updatedProtocol := "udp"

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
			// 1resource create/ delete 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalId, externalPort, virtualMachineType, internalIp, internalPort, protocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(internalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(externalPort)),
					resource.TestCheckResourceAttr(resourceName, "internal_ip", internalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
					resource.TestCheckResourceAttr(resourceName, "virtual_machine_type", strconv.Itoa(virtualMachineType)),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				//	2 resource update验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalId, updatedExternalPort, virtualMachineType, updatedInternalIp, updatedInternalPort, updatedProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(resourceName, "internal_ip", updatedInternalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
					resource.TestCheckResourceAttr(resourceName, "virtual_machine_type", strconv.Itoa(virtualMachineType)),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				// 3 datasource 验证
				//Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalId, updatedExternalPort, virtualMachineType, updatedInternalIp, updatedInternalPort, updatedProtocol) +
				Config: utils.LoadTestCase(datasourceFile, dnd, natGatewayId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "dnats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.protocol", updatedProtocol),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_ip", updatedInternalIp),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalId, updatedExternalPort, virtualMachineType, updatedInternalIp, updatedInternalPort, updatedProtocol),
				Destroy: true,
			},
		},
	})
}
