package mongodb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunMongodbClusterInstance(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	resourceFile := "resource_ctyun_mongodb_instance_single_cycle_no_az_os.tf"
	cycleType := "month"
	cycleCount := 1
	//autoRenew := false
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-month" + utils.GenerateRandomString()
	password := "Kqjwyk123="
	prodId := "Single40"
	readPort := 12345
	updatedReadPort := 12356
	flavorName := "s7.large.2"
	updateName := "tf-mongodb-new" + utils.GenerateRandomString()
	storageType := "SATA"
	storageSpace := 100
	backupStorageType := "OS"
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
			// 创建mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId,
					readPort, storageType, storageSpace, backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "prod_id", prodId),
				),
			},
			// 更新mongodb实例名称和端口号
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, flavorName, subnetID, securityGroupID, updateName, password, prodId,
					updatedReadPort, storageType, storageSpace, backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "prod_id", prodId),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
				),
			},
			// destroy
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, flavorName, subnetID, securityGroupID, updateName, password, prodId,
					updatedReadPort, storageType, storageSpace, backupStorageType),
				Destroy: true,
			},
		},
	})
}
