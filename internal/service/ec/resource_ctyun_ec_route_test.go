package ec_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunExpressConnectRoute(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_express_connect_route." + rnd
	resourceFile := "resource_ctyun_ec_route.tf"
	resourceFile1 := "resource_ctyun_ec_route_blackhole.tf"

	// 从环境变量获取测试依赖资源
	ecID := "49410d6d-fd53-48b3-9f78-cb28da38d7be"
	cgwID := "85de16c1-12d8-4608-aea1-eae75843af25"
	rtbID := "beacf1e4-952a-451b-b7be-4df122b36df8"
	cidr := "192.168.1.3/32"
	ipVersion := "ipv4"
	nextHopID := "vpc-obo8cwurdi"
	nextHopType := "vpc"

	initDescription := "init description for route test"
	initDescription1 := "init black hole description for route test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建普通路由测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID, cidr, ipVersion, initDescription,
					false, nextHopType, nextHopID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "cgw_id", cgwID),
					resource.TestCheckResourceAttr(resourceName, "rtb_id", rtbID),
					resource.TestCheckResourceAttr(resourceName, "cidr", cidr),
					resource.TestCheckResourceAttr(resourceName, "next_hop_type", nextHopType),
					resource.TestCheckResourceAttr(resourceName, "next_hop_id", nextHopID),
					resource.TestCheckResourceAttr(resourceName, "ip_version", ipVersion),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "is_black_hole_route", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID, cidr, ipVersion, initDescription,
					false, nextHopType, nextHopID),
				Destroy: true,
			},

			// 2. 创建黑洞路由测试
			{
				Config: utils.LoadTestCase(
					resourceFile1, rnd,
					ecID, cgwID, rtbID, cidr, ipVersion, initDescription1,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cidr", cidr),
					resource.TestCheckResourceAttr(resourceName, "ip_version", ipVersion),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription1),
					resource.TestCheckResourceAttr(resourceName, "is_black_hole_route", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile1, rnd,
					ecID, cgwID, rtbID, cidr, ipVersion, initDescription1,
					true,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: IPv6路由测试
func TestAccCtyunExpressConnectRouteIPv6(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_express_connect_route." + rnd
	resourceFile := "resource_ctyun_ec_route.tf"

	ecID := "49410d6d-fd53-48b3-9f78-cb28da38d7be"
	cgwID := "85de16c1-12d8-4608-aea1-eae75843af25"
	rtbID := "beacf1e4-952a-451b-b7be-4df122b36df8"
	cidr := "2001:db8::/32"
	ipVersion := "ipv6"
	description := "IPv6 route description"
	isBlackHole := false
	nextHopType := "vpc"
	nextHopID := "vpc-obo8cwurdi"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建IPv6路由测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID,
					cidr, ipVersion, description, isBlackHole, nextHopType, nextHopID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip_version", "ipv6"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "2001:db8::/32"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID,
					cidr, ipVersion, description, isBlackHole, nextHopType, nextHopID),
				Destroy: true,
			},
		},
	})
}

// 测试用例4: 不同下一跳类型测试
func TestAccCtyunExpressConnectRouteNextHopTypes(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceFile := "resource_ctyun_ec_route.tf"

	datasourceFile := "datasource_ctyun_ec_routes.tf"
	datasourceName := "data.ctyun_express_connect_routes." + dnd
	// 从环境变量获取测试依赖资源
	ecID := "49410d6d-fd53-48b3-9f78-cb28da38d7be"
	cgwID := "85de16c1-12d8-4608-aea1-eae75843af25"
	rtbID := "beacf1e4-952a-451b-b7be-4df122b36df8"

	// 测试不同的下一跳类型
	nextHopTypesAndIDMap := map[string]string{
		"vpc": "vpc-obo8cwurdi",
		//"cda":   "",
		//"vpn":   "",
		//"cross": ""
	}
	i := 1
	for nextHopType, nextHopID := range nextHopTypesAndIDMap {
		t.Run(nextHopType, func(t *testing.T) {
			resourceName := "ctyun_express_connect_route." + rnd + "_" + nextHopType
			cidr := fmt.Sprintf("192.168.1.%d/32", i)
			ipVersion := "ipv4"
			description := fmt.Sprintf("Route with %s next hop", nextHopType)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					// 1. 创建路由测试
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd+"_"+nextHopType,
							ecID, cgwID, rtbID,
							cidr, ipVersion, description, false, nextHopType, nextHopID,
						),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "next_hop_type", nextHopType),
							resource.TestCheckResourceAttrSet(resourceName, "id"),
						),
					},
					{
						Config: utils.LoadTestCase(datasourceFile, dnd,
							ecID, cgwID, rtbID),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttrSet(datasourceName, "routes.#")),
					},
					// 2. 清理资源
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd+"_"+nextHopType,
							ecID, cgwID, rtbID,
							cidr, ipVersion, description, false, nextHopType, nextHopID,
						),
						Destroy: true,
					},
				},
			})
		})
		i = i + 1
	}
}
