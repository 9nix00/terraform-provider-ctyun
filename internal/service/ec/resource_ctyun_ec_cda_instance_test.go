package ec_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"regexp"
	"testing"
)

func TestAccCtyunEcCdaInstance_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_ec_cda_instance." + rnd
	resourceFile := "resource_ctyun_ec_cda_instance.tf"
	//datasourceFile := "datasource_ctyun_ec_cda_instances.tf"
	//datasourceName := "data.ctyun_ec_cda_instances." + dnd

	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	cgwID := dependence.cloudGatewayId
	cdaID := "VTDQY51LO0UWDZ96ALY4"
	cdaName := "shih14@chinatelecom.cn-GATEWAY华北2-11"
	rtbID := dependence.rtbID
	initialCidrV4List := `"192.168.1.0/24"`
	updatedCidrV4List := `"192.168.1.0/24"`
	////// CIDR列表
	//initialCidrV4List := `"192.168.1.0/24"`
	//updatedCidrV4List := `"192.168.1.0/24", "192.168.2.0/24"`

	// CDA信息（JSON格式）
	cdaInfo := `"{\"cdaCidrV4List\":[\"192.168.1.0/24\"]}"`
	updatedCdaInfo := `"{\"cdaCidrV4List\":[\"192.168.1.0/24\",\"192.168.2.0/24\"]}"`

	account := "shih14@chinatelecom.cn"
	weights := 50
	routeLearn := 1
	routeSync := 1

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建CDA网络实例测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, cdaID, cdaName, initialCidrV4List, rtbID, cdaInfo, account, weights, routeLearn, routeSync,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "cgw_id", cgwID),
					resource.TestCheckResourceAttr(resourceName, "cda_id", cdaID),
					resource.TestCheckResourceAttr(resourceName, "cda_name", cdaName),
					resource.TestCheckResourceAttr(resourceName, "rtb_id", rtbID),
					resource.TestCheckResourceAttr(resourceName, "cda_cidr_v4_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cda_cidr_v4_list.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(resourceName, "weights", "50"),
					resource.TestCheckResourceAttr(resourceName, "route_learn", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_sync", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新CIDR列表测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, cdaID, cdaName, updatedCidrV4List, rtbID, updatedCdaInfo, account, weights, routeLearn, routeSync,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cda_cidr_v4_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cda_cidr_v4_list.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(resourceName, "cda_cidr_v4_list.1", "192.168.2.0/24"),
				),
			},
			// 3. 数据源测试（暂时注释，因为还没有实现数据源）
			//{
			//	Config: utils.LoadTestCase(
			//		resourceFile, rnd,
			//		ecID, cgwID, cdaID, cdaName, updatedCidrV4List, rtbID, updatedCdaInfo, weights, routeLearn, routeSync,
			//	) + utils.LoadTestCase(datasourceFile, dnd, ecID),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttrSet(datasourceName, "cda_instances.#")),
			//},
			// 4. 导入测试（当前资源不支持导入）
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return rs.Primary.ID, nil
				},
				ExpectError: regexp.MustCompile("This resource does not support import"),
			},
			// 5. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, cdaID, cdaName, updatedCidrV4List, rtbID, updatedCdaInfo, account, weights, routeLearn, routeSync,
				),
				Destroy: true,
			},
		},
	})
}
