package scaling_test

//
//import (
//	"fmt"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
//	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//	"github.com/hashicorp/terraform-plugin-testing/terraform"
//	"os"
//	"strconv"
//	"testing"
//)
//
//func TestAccCtyunScalingEcs(t *testing.T) {
//	// 设置环境变量以启用TF_ACC测试
//	err := os.Setenv("TF_ACC", "1")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// 生成随机名称避免冲突
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_scaling_ecs." + rnd
//	resourceFile := "resource_ctyun_scaling_ecs.tf"
//
//	// 测试依赖项（需要提前创建好）
//	groupID, err := strconv.ParseInt(dependence.scalingGroupID, 10, 64)
//	if err != nil {
//		fmt.Println(err)
//	}
//	// 替换为实际伸缩组ID
//	instanceUUIDList := fmt.Sprintf(`["%s","%s"]`, dependence.instanceUUID, dependence.instanceUUID1) // 替换为实际云主机ID列表
//
//	// 创建参数
//	protectStatus := "enable"
//
//	// 更新参数
//	updatedProtectStatus := "disable"
//	updatedInstanceUUIDList := fmt.Sprintf(`["%s"]`, dependence.instanceUUID) // 移除一个实例
//
//	resource.Test(t, resource.TestCase{
//		// 检查资源是否被销毁
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource %s still exists", resourceName)
//			}
//			return nil
//		},
//		// 使用ProtoV6ProviderFactories
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 创建伸缩组云主机关联
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, groupID, instanceUUIDList, protectStatus),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "group_id", fmt.Sprintf("%d", groupID)),
//					resource.TestCheckResourceAttr(resourceName, "protect_status", protectStatus),
//					resource.TestCheckResourceAttr(resourceName, "instance_uuid_list.#", "2"),
//				),
//			},
//			// 更新保护状态和实例列表
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, groupID, updatedInstanceUUIDList, updatedProtectStatus),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "protect_status", updatedProtectStatus),
//					resource.TestCheckResourceAttr(resourceName, "instance_uuid_list.#", "1"),
//				),
//			},
//			// 销毁资源（通过空配置）
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, groupID, updatedInstanceUUIDList, updatedProtectStatus),
//				Destroy: false,
//			},
//		},
//	})
//}
