package vpce_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var initMain bool

func TestMain(m *testing.M) {
	initMain = true
	err := initSharedResources()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	code := m.Run()

	err = destroySharedResources()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(code)
}

var sharedEcsID string
var sharedVpcID string
var sharedSubnetID string
var sharedVpceServerID string

func initSharedResources() error {
	// 应用配置
	cmd := exec.Command("terraform", "apply", "-auto-approve", "-input=false")
	cmd.Dir = "testdata/shared"
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apply failed: %s\n%s", err, string(out))
	}

	// 获取输出值
	cmd = exec.Command("terraform", "output", "-json")
	cmd.Dir = "testdata/shared"
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("output failed: %s", err)
	}

	// 解析输出到环境变量
	var outputs map[string]struct {
		Value string `json:"value"`
	}
	if err = json.Unmarshal(out, &outputs); err != nil {
		return fmt.Errorf("output parsing failed: %s", err)
	}
	sharedVpcID, sharedSubnetID, sharedEcsID, sharedVpceServerID =
		outputs["vpc_id"].Value, outputs["subnet_id"].Value, outputs["ecs_id"].Value, outputs["vpce_server_id"].Value

	return nil
}

func destroySharedResources() error {
	cmd := exec.Command("terraform", "destroy", "-auto-approve", "-input=false")
	cmd.Dir = "testdata/shared"
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("destroy failed: %s\n%s", err, string(out))
	}
	return nil
}
