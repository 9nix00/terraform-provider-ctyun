package elb_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunElbListener(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_elb_listener." + rnd
	resourceFile := "resource_ctyun_elb_listener.tf"

	datasourceName := "data.ctyun_elb_listeners." + dnd
	datasourceFile := "datasource_ctyun_elb_listeners.tf"
	loadbalanceID := dependence.loadBalanceID2
	name := "listener-" + utils.GenerateRandomString()

	updatedName := "listener-new" + utils.GenerateRandomString()

	protocolTCP := "TCP"
	//protocolUDP := "UDP"
	ProtocolHTTP := "HTTP"
	//ProtocolHTTPS := "HTTPS"
	//
	protocolPort := utils.GenerateRandomPort(1, 65535)
	defaultActionType := "forward"
	// 当default action type = forward, target_groups 必填。
	// 当default action type = redirect, redirectListenerID必填
	// target_groups 和redirectListenerID 用%[7]s来控制
	//targetGroupIds := fmt.Sprintf(`{target_group_id="%s"},{target_group_id="%s"}`, dependence.targetGroupID, dependence.targetGroupID2)
	targetGroupIds := fmt.Sprintf(`{target_group_id="%s"}`, dependence.targetGroupID2)
	updatedTargetGroupIds := fmt.Sprintf(`{target_group_id="%s"}`, dependence.targetGroupID3)

	tfTargetGroupID := fmt.Sprintf(`target_groups=[%s]`, targetGroupIds)
	updatedTargetGroupID := fmt.Sprintf(`target_groups=[%s]`, updatedTargetGroupIds)

	// nat64 需要开始ipv6

	// qps,支持http/https
	tfQPS := fmt.Sprintf(`listener_qps=%d`, 1)
	// cps 支持tcp/udp
	tfCPS := fmt.Sprintf(`listener_cps=%d`, 1)
	// establish_timeout, 仅支持tcp，建立连接超时时间，单位秒，取值范围： 1 - 1800
	tfEstablishTimeout := fmt.Sprintf(`establish_timeout=%d`, 100)
	// idle_timeout, 支持http/https，链接空闲断开超时时间，单位秒，取值范围：1 - 300
	tfIdleTimeout := fmt.Sprintf(`idle_timeout=%d`, 100)
	// response_timeout，支持http/https
	tfResponseTimeout := fmt.Sprintf(`response_timeout=%d`, 100)

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
			// 1. protocol=TCP， defaultActionType=forward, targetGroupID必填
			// 1.1 Create验证
			{
				//Create验证
				Config: utils.LoadTestCase(resourceFile, rnd, loadbalanceID, name, protocolTCP, protocolPort, defaultActionType, tfTargetGroupID, "", "", "", "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocolTCP),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
					resource.TestCheckResourceAttr(resourceName, "default_action_type", defaultActionType),
				),
			},
			// 1.2 update验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, loadbalanceID, updatedName, protocolTCP, protocolPort, defaultActionType, updatedTargetGroupID, "", "", tfCPS, tfEstablishTimeout, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocolTCP),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
					resource.TestCheckResourceAttr(resourceName, "default_action_type", defaultActionType),
					resource.TestCheckResourceAttr(resourceName, "listener_cps", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "establish_timeout", strconv.Itoa(100)),
				),
			},
			// 1.3 datasource验证
			{
				Config: utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "listeners.#", "2"),
				),
			},
			// 1.4 destroy验证
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, loadbalanceID, updatedName, protocolTCP, protocolPort, defaultActionType, tfTargetGroupID, "", "", tfCPS, tfEstablishTimeout, "", ""),
				Destroy: true,
			},
			// 2 详细信息验证，protocol=HTTP， defaultActionType=forward, targetGroupID必填
			// 2.1 Create
			{
				Config: utils.LoadTestCase(resourceFile, rnd, loadbalanceID, name, ProtocolHTTP, protocolPort, defaultActionType, tfTargetGroupID, "", "", "", "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", ProtocolHTTP),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
					resource.TestCheckResourceAttr(resourceName, "default_action_type", defaultActionType),
				),
			},
			// 2.2 update
			{
				Config: utils.LoadTestCase(resourceFile, rnd, loadbalanceID, updatedName, ProtocolHTTP, protocolPort, defaultActionType, updatedTargetGroupID, "", tfQPS, "", "", tfIdleTimeout, tfResponseTimeout),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "protocol", ProtocolHTTP),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
					resource.TestCheckResourceAttr(resourceName, "default_action_type", defaultActionType),
					resource.TestCheckResourceAttr(resourceName, "listener_qps", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "idle_timeout", strconv.Itoa(100)),
					resource.TestCheckResourceAttr(resourceName, "response_timeout", strconv.Itoa(100)),
				),
			},
			// 2.3 destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, loadbalanceID, updatedName, protocolTCP, protocolPort, defaultActionType, updatedTargetGroupID, "", tfQPS, "", "", tfIdleTimeout, tfResponseTimeout),
				Destroy: true,
			},
		},
	})
}
