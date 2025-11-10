package ccse_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunCcseScalingNodePoolPolicy(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ccse_node_pool_scaling_policy." + rnd
	resourceFile := "resource_ctyun_ccse_node_pool_scaling_policy.tf"

	nodePoolName := "default"
	name := fmt.Sprintf("%s-%s", nodePoolName, dependence.clusterID)
	valuesYaml := fmt.Sprintf(`kind: HorizontalNodeAutoscaler
apiVersion: autoscaler.ccse.ctyun.cn/v1
metadata:
  name: default-%s
spec:
  disable: false
  rules:
    - action:
        type: ScaleUp
        unit: Node
        value: 1
      cronTrigger:
        schedule: "%s"
      disable: false
      ruleName: rule04130
      type: Cron
    - action:
        type: ScaleUp
        unit: Node
        value: 1
      disable: false
      metricTrigger:
        metricName: cpu_util
        metricOperation: gt
        metricValue: "80"
      ruleName: rule09974
      type: Metric
    - action:
        type: ScaleUp
        unit: Node
        value: 1
      disable: false
      alarmTrigger:
        evaluationCount: 2
        fun: avg
        metric: cpu_util
        operator: ge
        value: "80"
        period: 5m
      ruleName: rule42857
      type: Alarm
  targetNodepools:
    - default
  coolDown: 3
`, dependence.clusterID, "8 20 * * *")
	valuesYaml2 := fmt.Sprintf(`kind: HorizontalNodeAutoscaler
apiVersion: autoscaler.ccse.ctyun.cn/v1
metadata:
  name: default-%s
spec:
  disable: false
  rules:
    - action:
        type: ScaleUp
        unit: Node
        value: 1
      cronTrigger:
        schedule: "%s"
      disable: false
      ruleName: rule04130
      type: Cron
    - action:
        type: ScaleUp
        unit: Node
        value: 1
      disable: false
      metricTrigger:
        metricName: cpu_util
        metricOperation: gt
        metricValue: "80"
      ruleName: rule09974
      type: Metric
    - action:
        type: ScaleUp
        unit: Node
        value: 1
      disable: false
      alarmTrigger:
        evaluationCount: 2
        fun: avg
        metric: cpu_util
        operator: ge
        value: "80"
        period: 5m
      ruleName: rule42857
      type: Alarm
  targetNodepools:
    - default
  coolDown: 3
`, dependence.clusterID, "8 22 * * *")

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
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					valuesYaml,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", dependence.clusterID),
					resource.TestCheckResourceAttr(resourceName, "values_yaml", valuesYaml),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "node_pool_name", nodePoolName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "actual_config"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					valuesYaml2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", dependence.clusterID),
					resource.TestCheckResourceAttr(resourceName, "values_yaml", valuesYaml2),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "node_pool_name", nodePoolName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "actual_config"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"values_yaml",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					valuesYaml2,
				),
				Destroy: true,
			},
		},
	})
}
