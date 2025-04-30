package vpc_test

import (
	"fmt"
	"os"
	"terraform-provider-ctyun/internal/extend/terraform"
	"testing"
)

const dependenceDir = "testdata/dependence"

var dependenceVpcID string

func TestMain(m *testing.M) {
	// 初始化依赖资源
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dependenceVpcID = outputs["vpc_id"].Value

	// 执行测试用例
	code := m.Run()

	// 清理依赖资源
	err = terraform.DestroyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(code)
}
