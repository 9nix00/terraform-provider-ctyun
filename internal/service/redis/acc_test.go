package redis_test

import (
	"fmt"
	"os"
	"terraform-provider-ctyun/internal/extend/terraform"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID              string
	subnetID           string
	securityGroupID    string
	eipAddress         string
	redisVersion       string
	redisEngineEdition string
}

var dependence Dependence

func TestMain(m *testing.M) {
	// 初始化依赖资源
	fmt.Println("开始初始化依赖资源")
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dependence = Dependence{
		vpcID:              outputs["vpc_id"].Value,
		subnetID:           outputs["subnet_id"].Value,
		securityGroupID:    outputs["security_group_id"].Value,
		eipAddress:         outputs["eip_address"].Value,
		redisVersion:       outputs["redis_version"].Value,
		redisEngineEdition: outputs["redis_engine_edition"].Value,
	}
	fmt.Println("依赖资源初始化完毕")

	// 执行测试用例
	code := m.Run()

	fmt.Println("开始清理依赖资源")
	// 清理依赖资源
	err = terraform.DestroyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("依赖资源清理完毕")

	os.Exit(code)
}
