package vpce_test

import (
	"fmt"
	"os"
	"terraform-provider-ctyun/internal/extend/terraform"
	"testing"
)

const sharedDir = "testdata/shared"

var initMain bool
var sharedEcsID string
var sharedVpcID string
var sharedSubnetID string
var sharedVpceServerID string

func TestMain(m *testing.M) {
	initMain = true
	// 初始化共享资源
	err := initSharedResources()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 执行测试用例
	code := m.Run()

	// 清理共享资源
	err = terraform.DestroyResources(sharedDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(code)
}

func initSharedResources() error {
	outputs, err := terraform.ApplyResources(sharedDir)
	if err != nil {
		return err
	}
	sharedEcsID, sharedVpcID, sharedSubnetID, sharedVpceServerID =
		outputs["ecs_id"].Value, outputs["vpc_id"].Value, outputs["subnet_id"].Value, outputs["vpce_server_id"].Value
	return nil
}
