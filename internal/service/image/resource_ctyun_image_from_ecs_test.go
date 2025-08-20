package image_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunImageFromEcs_basic(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_system_disk.tf"

	imageName := "tf-image-test-" + utils.GenerateRandomString()
	updatedImageName := "tf-image-test-updated-" + utils.GenerateRandomString()

	description := "测试镜像描述"
	updatedDescription := description + "-updated"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// 测试创建
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					imageName,
					description,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
					// 如果标签信息也无法通过API获取，也可以添加
					// "labels",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunImageFromEcs_dataDisk(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_data_disk.tf"

	imageName := "tf-image-data-" + utils.GenerateRandomString()
	updatedImageName := "tf-image-data-updated-" + utils.GenerateRandomString()

	description := "测试数据盘创建镜像描述"
	updatedDescription := description + "-updated"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// 测试数据盘创建镜像
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					imageName,
					description,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunImageFromEcs_entire(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_entire.tf"

	imageName := "tf-image-test-" + utils.GenerateRandomString()
	updatedImageName := "tf-image-test-updated-" + utils.GenerateRandomString()

	description := "测试整机镜像描述"
	updatedDescription := description + "-updated"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// 测试创建
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					imageName,
					description,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
				),
				Destroy: true,
			},
		},
	})
}
