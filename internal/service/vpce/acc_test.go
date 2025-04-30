package vpce_test

import (
	"fmt"
	"os"
	"terraform-provider-ctyun/internal/extend/terraform"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	ecsID               string
	vpcID               string
	subnetID            string
	vpceServerID        string
	reverseVpceServerID string
	vpceID              string
	transitIP           string
	targetIP            string
}

var dependence Dependence

func TestMain(m *testing.M) {
	// 初始化依赖资源
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dependence = Dependence{
		ecsID:               outputs["ecs_id"].Value,
		vpcID:               outputs["vpc_id"].Value,
		subnetID:            outputs["subnet_id"].Value,
		vpceServerID:        outputs["vpce_server_id"].Value,
		reverseVpceServerID: outputs["reverse_vpce_server_id"].Value,
		vpceID:              outputs["vpce_id"].Value,
		transitIP:           outputs["vpce_server_transit_ip"].Value,
		targetIP:            outputs["ecs_fixed_ip"].Value,
	}

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
