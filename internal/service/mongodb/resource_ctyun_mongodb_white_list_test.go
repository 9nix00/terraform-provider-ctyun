package mongodb_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMongodbWhiteList_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_white_list." + rnd
	resourceFile := "resource_ctyun_mongodb_white_list.tf"

	instance_id := "98f229ddda674219a0930da14cb78394"
	groupName := utils.GenerateRandomString()
	ipType := "ipv4"
	whiteListType := "2"
	ipList := "10.138.16.8"
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
			// 基本功能验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, groupName, ipType, whiteListType, ipList),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "ipType", ipType),
					resource.TestCheckResourceAttr(resourceName, "whiteListType", whiteListType),
					resource.TestCheckResourceAttr(resourceName, "ipList", ipList),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, groupName, ipType, whiteListType, ipList),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "ipType", ipType),
					resource.TestCheckResourceAttr(resourceName, "whiteListType", whiteListType),
					resource.TestCheckResourceAttr(resourceName, "ipList", ipList),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, groupName, ipType, whiteListType, ipList),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "ipType", ipType),
					resource.TestCheckResourceAttr(resourceName, "whiteListType", whiteListType),
					resource.TestCheckResourceAttr(resourceName, "ipList", ipList),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, instance_id, groupName, ipType, whiteListType, ipList),
				Destroy: true,
			},
		},
	})
}

func TestAccMongodbWhiteList_update(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_account." + rnd
	resourceFile := "resource_ctyun_mongodb_account.tf"

	instId := dependence.mongodbID

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
				Config: utils.LoadTestCase(resourceFile, rnd, instId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "privileges", "read"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "privileges", "readWrite"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, instId),
				Destroy: true,
			},
		},
	})
}

func TestMongodbWhiteValidation(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceFile := "resource_ctyun_mongodb_account.tf"
	instId := dependence.mongodbID

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      utils.LoadTestCase(resourceFile, rnd, instId),
				ExpectError: regexp.MustCompile("实例名称不符合规范"),
			},
			{
				Config:      utils.LoadTestCase(resourceFile, rnd, instId),
				ExpectError: regexp.MustCompile("密码必须包含大写字母、小写字母、数字、特殊字符中的至少三种"),
			},
		},
	})
}
