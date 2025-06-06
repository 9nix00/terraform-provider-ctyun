package mysql_test

import (
	"fmt"
	"os"
	"terraform-provider-ctyun/internal/extend/terraform"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID           string
	subnetID        string
	securityGroupID string
	eipID           string
	eipAddress      string
	mysqlID         string
}

var dependence Dependence

func TestMain(m *testing.M) {
	fmt.Println("开始初始化依赖资源")
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		terraform.DestroyResource(dependenceDir)
		os.Exit(1)
	}
	dependence = Dependence{
		vpcID:           outputs["vpc_id"].Value,
		subnetID:        outputs["subnet_id"].Value,
		securityGroupID: outputs["security_group_id"].Value,
		eipID:           outputs["eip_id"].Value,
		eipAddress:      outputs["eip_address"].Value,
		mysqlID:         outputs["mysql_id"].Value,
		//eipID:    outputs["eip_id"].Value,
		//vpcID:          outputs["vpc_id"].Value,
		//subnetID:       outputs["subnet_id"].Value,
		//loadBalanceID:  outputs["loadbalancer_id"].Value,
		//loadBalanceID2: outputs["loadbalancer_id_rule"].Value,
		//healthCheckID:  outputs["health_check_id"].Value,
		//targetGroupID:  outputs["target_group_id"].Value,
		//targetGroupID2: outputs["target_group_id2"].Value,
		//targetGroupID3: outputs["target_group_id3"].Value,
		//listenerID:     outputs["listener_id"].Value,
		//instanceID:     outputs["instance_id"].Value,
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
