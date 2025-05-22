package ccse_test

import (
	"fmt"
	"os"
	"terraform-provider-ctyun/internal/extend/terraform"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID      string
	subnetID   string
	flavorName string
	clusterID  string
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
		vpcID:      outputs["vpc_id"].Value,
		subnetID:   outputs["subnet_id"].Value,
		flavorName: outputs["flavor_name"].Value,
		clusterID:  outputs["cluster_id"].Value,
	}
	fmt.Println("依赖资源初始化完毕")

	// 执行测试用例
	code := m.Run()
	terraform.DestroyResource(dependenceDir)

	// ccse依赖的子网无法马上删除
	//fmt.Println("开始清理依赖资源")
	//// 清理依赖资源
	//err = terraform.DestroyResource(dependenceDir)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println("依赖资源清理完毕")

	os.Exit(code)
}
