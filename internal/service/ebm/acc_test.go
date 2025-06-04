package ebm_test

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
	deviceType      string
	smartNicExist   bool
	supportCloud    bool
	cloudBoot       bool
	systemRaid      string
	dataRaid        string
	imageUUID       string
}

var dependence Dependence

func TestMain(m *testing.M) {
	// 初始化依赖资源
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
		deviceType:      outputs["device_type"].Value,
		systemRaid:      outputs["system_raid"].Value,
		dataRaid:        outputs["data_raid"].Value,
		imageUUID:       outputs["image_uuid"].Value,
	}
	if outputs["smart_nic_exist "].Value == "true" {
		dependence.smartNicExist = true
	}
	if outputs["support_cloud"].Value == "true" {
		dependence.supportCloud = true
	}
	if outputs["cloud_boot"].Value == "true" {
		dependence.cloudBoot = true
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
