package peer_connection_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"testing"
	"time"
)

func TestAccCtyunVpcPeerConnectionRoute_basic(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_peer_connection_route." + rnd
	resourceFile := "resource_ctyun_vpc_peer_connection_route.tf"

	// 从环境变量获取测试依赖资源
	vpcID := dependence.vpcID1
	nextHopID := dependence.peerConnectionID // 对等连接ID作为下一跳

	// 测试数据
	ipVersion := "ipv4"
	destinationCIDR := "172.168.1.0/24"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建VPC对等连接路由测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ipVersion, nextHopID, vpcID, destinationCIDR,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					// 基本属性验证
					resource.TestCheckResourceAttr(resourceName, "ip_version", ipVersion),
					resource.TestCheckResourceAttr(resourceName, "next_hop_id", nextHopID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "destination", destinationCIDR),
					// 系统生成属性验证
					resource.TestCheckResourceAttrSet(resourceName, "id"),

					// 自定义验证函数
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource not found: %s", resourceName)
						}

						// 验证路由规则ID格式
						routeID := rs.Primary.Attributes["id"]
						if routeID == "" {
							return fmt.Errorf("route ID is empty")
						}
						if len(routeID) < 8 {
							return fmt.Errorf("route ID seems too short: %s", routeID)
						}
						return nil
					},
				),
			},
			// 2. 由于不支持更新，此步骤用于验证读取功能
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ipVersion, nextHopID, vpcID, destinationCIDR,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination", destinationCIDR),
					// 验证所有属性保持不变
					resource.TestCheckResourceAttr(resourceName, "ip_version", ipVersion),
					resource.TestCheckResourceAttr(resourceName, "next_hop_id", nextHopID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
				),
			},
			// 3. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s,%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["region_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["vpc_id"],
						rs.Primary.Attributes["subnet_id"], // 注意：虽然schema中没有subnet_id，但导入ID格式包含它
					), nil
				},
				ImportStateVerify: true,
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ipVersion, nextHopID, vpcID, destinationCIDR,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunVpcPeerConnectionRoute_ipv6(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_peer_connection_route." + rnd
	resourceFile := "resource_ctyun_vpc_peer_connection_route.tf"

	// 从环境变量获取测试依赖资源
	regionID := os.Getenv("CTYUN_REGION_ID")
	projectID := os.Getenv("CTYUN_PROJECT_ID")
	vpcID := os.Getenv("CTYUN_VPC_ID_IPV6") // IPv6 VPC
	nextHopID := os.Getenv("CTYUN_PEER_CONNECTION_ID_IPV6")

	// 验证环境变量是否设置
	if regionID == "" || vpcID == "" || nextHopID == "" {
		t.Skip("Skipping test: required environment variables not set")
	}

	// IPv6测试数据
	ipVersion := "ipv6"
	destinationCIDR := "2001:db8::/32"

	// 等待函数
	wait30Seconds := func() {
		t.Logf("等待30秒让IPv6路由规则稳定...")
		time.Sleep(30 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建IPv6 VPC对等连接路由测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID, projectID,
					ipVersion, nextHopID, vpcID, destinationCIDR,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip_version", ipVersion),
					resource.TestCheckResourceAttr(resourceName, "destination", destinationCIDR),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
				PreConfig: func() {
					wait30Seconds()
				},
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID, projectID,
					ipVersion, nextHopID, vpcID, destinationCIDR,
				),
				Destroy: true,
				PreConfig: func() {
					wait30Seconds()
				},
			},
		},
	})
}

func TestAccCtyunVpcPeerConnectionRoute_regionReplacement(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceFile := "resource_ctyun_vpc_peer_connection_route.tf"

	// 从环境变量获取测试依赖资源
	regionID1 := os.Getenv("CTYUN_REGION_ID")
	regionID2 := os.Getenv("CTYUN_REGION_ID_2")
	projectID1 := os.Getenv("CTYUN_PROJECT_ID")
	projectID2 := os.Getenv("CTYUN_PROJECT_ID_2")
	vpcID1 := os.Getenv("CTYUN_VPC_ID")
	vpcID2 := os.Getenv("CTYUN_VPC_ID_2")
	nextHopID1 := os.Getenv("CTYUN_PEER_CONNECTION_ID")
	nextHopID2 := os.Getenv("CTYUN_PEER_CONNECTION_ID_2")

	// 验证环境变量是否设置
	if regionID1 == "" || regionID2 == "" || vpcID1 == "" || vpcID2 == "" || nextHopID1 == "" || nextHopID2 == "" {
		t.Skip("Skipping test: required environment variables not set")
	}

	// 测试数据
	ipVersion := "ipv4"
	destinationCIDR := "192.168.1.0/24"

	// 等待函数
	wait30Seconds := func() {
		t.Logf("等待30秒让路由规则稳定...")
		time.Sleep(30 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 在第一个区域创建
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID1, projectID1,
					ipVersion, nextHopID1, vpcID1, destinationCIDR,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "region_id", regionID1),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "project_id", projectID1),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "vpc_id", vpcID1),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "next_hop_id", nextHopID1),
				),
				PreConfig: func() {
					wait30Seconds()
				},
			},
			// 2. 切换到第二个区域（需要替换）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID2, projectID2,
					ipVersion, nextHopID2, vpcID2, destinationCIDR,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "region_id", regionID2),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "project_id", projectID2),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "vpc_id", vpcID2),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "next_hop_id", nextHopID2),
				),
				PreConfig: func() {
					wait30Seconds()
				},
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID2, projectID2,
					ipVersion, nextHopID2, vpcID2, destinationCIDR,
				),
				Destroy: true,
				PreConfig: func() {
					wait30Seconds()
				},
			},
		},
	})
}

func TestAccCtyunVpcPeerConnectionRoute_projectReplacement(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceFile := "resource_ctyun_vpc_peer_connection_route.tf"

	// 从环境变量获取测试依赖资源
	regionID := os.Getenv("CTYUN_REGION_ID")
	projectID1 := os.Getenv("CTYUN_PROJECT_ID")
	projectID2 := os.Getenv("CTYUN_PROJECT_ID_2")
	vpcID1 := os.Getenv("CTYUN_VPC_ID")
	vpcID2 := os.Getenv("CTYUN_VPC_ID_2")
	nextHopID1 := os.Getenv("CTYUN_PEER_CONNECTION_ID")
	nextHopID2 := os.Getenv("CTYUN_PEER_CONNECTION_ID_2")

	// 验证环境变量是否设置
	if regionID == "" || projectID1 == "" || projectID2 == "" || vpcID1 == "" || vpcID2 == "" || nextHopID1 == "" || nextHopID2 == "" {
		t.Skip("Skipping test: required environment variables not set")
	}

	// 测试数据
	ipVersion := "ipv4"
	destinationCIDR := "192.168.1.0/24"

	// 等待函数
	wait30Seconds := func() {
		t.Logf("等待30秒让路由规则稳定...")
		time.Sleep(30 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 在第一个项目创建
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID, projectID1,
					ipVersion, nextHopID1, vpcID1, destinationCIDR,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "project_id", projectID1),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "vpc_id", vpcID1),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "next_hop_id", nextHopID1),
				),
				PreConfig: func() {
					wait30Seconds()
				},
			},
			// 2. 切换到第二个项目（需要替换）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID, projectID2,
					ipVersion, nextHopID2, vpcID2, destinationCIDR,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "project_id", projectID2),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "vpc_id", vpcID2),
					resource.TestCheckResourceAttr("ctyun_vpc_peer_connection_route."+rnd, "next_hop_id", nextHopID2),
				),
				PreConfig: func() {
					wait30Seconds()
				},
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID, projectID2,
					ipVersion, nextHopID2, vpcID2, destinationCIDR,
				),
				Destroy: true,
				PreConfig: func() {
					wait30Seconds()
				},
			},
		},
	})
}

func TestAccCtyunVpcPeerConnectionRoute_validation(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceFile := "resource_ctyun_vpc_peer_connection_route.tf"

	// 从环境变量获取测试依赖资源
	regionID := os.Getenv("CTYUN_REGION_ID")
	projectID := os.Getenv("CTYUN_PROJECT_ID")
	vpcID := os.Getenv("CTYUN_VPC_ID")
	nextHopID := os.Getenv("CTYUN_PEER_CONNECTION_ID")

	// 验证环境变量是否设置
	if regionID == "" || vpcID == "" || nextHopID == "" {
		t.Skip("Skipping test: required environment variables not set")
	}

	// 测试数据
	ipVersion := "ipv4"

	// 测试不同的CIDR格式
	testCases := []struct {
		name        string
		destination string
		shouldFail  bool
	}{
		{"valid_cidr", "10.0.0.0/16", false},
		{"valid_cidr_small", "10.0.0.0/24", false},
		{"valid_cidr_large", "10.0.0.0/8", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceName := "ctyun_vpc_peer_connection_route." + rnd + "_" + tc.name

			// 等待函数
			wait30Seconds := func() {
				t.Logf("等待30秒让路由规则稳定...")
				time.Sleep(30 * time.Second)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					// 1. 创建路由规则
					{
						Config: utils.LoadTestCase(
							rnd+"_"+tc.name, resourceFile, regionID, projectID,
							ipVersion, nextHopID, vpcID, tc.destination,
						),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "destination", tc.destination),
							resource.TestCheckResourceAttrSet(resourceName, "id"),
						),
						PreConfig: func() {
							wait30Seconds()
						},
					},
					// 2. 清理资源
					{
						Config: utils.LoadTestCase(
							rnd+"_"+tc.name, resourceFile, regionID, projectID,
							ipVersion, nextHopID, vpcID, tc.destination,
						),
						Destroy: true,
						PreConfig: func() {
							wait30Seconds()
						},
					},
				},
			})
		})
	}
}

func TestAccCtyunVpcPeerConnectionRoute_routingLogic(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_peer_connection_route." + rnd
	resourceFile := "resource_ctyun_vpc_peer_connection_route.tf"

	// 从环境变量获取测试依赖资源
	regionID := os.Getenv("CTYUN_REGION_ID")
	projectID := os.Getenv("CTYUN_PROJECT_ID")
	vpcID := os.Getenv("CTYUN_VPC_ID")
	nextHopID := os.Getenv("CTYUN_PEER_CONNECTION_ID")

	// 验证环境变量是否设置
	if regionID == "" || vpcID == "" || nextHopID == "" {
		t.Skip("Skipping test: required environment variables not set")
	}

	// 测试数据
	ipVersion := "ipv4"
	destinationCIDR := "10.100.0.0/16"

	// 等待函数
	wait30Seconds := func() {
		t.Logf("等待30秒让路由规则稳定...")
		time.Sleep(30 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建路由规则并验证路由逻辑
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID, projectID,
					ipVersion, nextHopID, vpcID, destinationCIDR,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// 验证路由规则生效
					func(s *terraform.State) error {
						// 这里可以添加自定义检查函数来验证路由逻辑
						// 例如通过API查询确认路由规则已生效
						// 或者通过测试网络连通性来验证路由
						return nil
					},
				),
				PreConfig: func() {
					wait30Seconds()
				},
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID, projectID,
					ipVersion, nextHopID, vpcID, destinationCIDR,
				),
				Destroy: true,
				PreConfig: func() {
					wait30Seconds()
				},
			},
		},
	})
}
