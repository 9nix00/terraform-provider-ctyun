package nat_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunPrivateDNat(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_private_nat_dnat." + rnd
	datasourceName := "data.ctyun_private_dnats." + dnd

	resourceFile := "resource_ctyun_private_dnat.tf"
	datasourceFile := "datasource_ctyun_private_dnat.tf"

	natGatewayId := dependence.privateNatID
	internalPort := utils.GenerateRandomPort(0, 65535)
	updatedInternalPort := utils.GenerateRandomPort(0, 65535)
	externalPort := utils.GenerateRandomPort(0, 1024)
	updatedExternalPort := utils.GenerateRandomPort(0, 1024)

	internalIp := "192.168.128.3"
	updatedInternalIp := "192.168.128.6"
	externalIp := "192.168.128.7"
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
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, protocol, externalPort, internalPort, internalIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(internalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(externalPort)),
					resource.TestCheckResourceAttr(resourceName, "internal_ip", internalIp),
					resource.TestCheckResourceAttr(resourceName, "external_ip", externalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				//	2 resource update验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, updatedInternalIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(resourceName, "internal_ip", updatedInternalIp),
					resource.TestCheckResourceAttr(resourceName, "external_ip", externalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				// 3 datasource 验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, updatedInternalIp) +
					utils.LoadTestCase(datasourceFile, dnd, natGatewayId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "dnats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.protocol", updatedProtocol),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_ip", updatedInternalIp),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, updatedInternalIp),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunPrivateDNat2(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_private_dnat." + rnd
	datasourceName := "data.ctyun_private_dnats." + dnd

	resourceFile := "resource_ctyun_private_dnat2.tf"
	datasourceFile := "datasource_ctyun_private_dnat.tf"
	externalIp := "192.168.1.5"
	natGatewayId := dependence.natID
	internalPort := utils.GenerateRandomPort(0, 65535)
	updatedInternalPort := utils.GenerateRandomPort(0, 65535)
	externalPort := utils.GenerateRandomPort(0, 1024)
	updatedExternalPort := utils.GenerateRandomPort(0, 1024)

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
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID, protocol, externalPort, internalPort, dependence.ecsID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(internalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(externalPort)),
					resource.TestCheckResourceAttr(resourceName, "port_id", dependence.ecsID),
					resource.TestCheckResourceAttr(resourceName, "external_ip", dependence.eipID),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			// 1resource create/ delete 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, dependence.ecsID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(resourceName, "port_id", dependence.ecsID),
					resource.TestCheckResourceAttr(resourceName, "external_ip", externalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				// 3 datasource 验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, dependence.ecsID) +
					utils.LoadTestCase(datasourceFile, dnd, natGatewayId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "dnats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.protocol", updatedProtocol),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.port_id", dependence.ecsID),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, dependence.ecsID),
				Destroy: true,
			},
		},
	})
}
