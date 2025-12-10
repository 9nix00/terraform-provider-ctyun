terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
}


resource "ctyun_ccse_node_pool_scaling_policy" "example" {
  cluster_id               = "dd92f3a6b034431bb7dceb849aed1220"
  values_yaml  = <<EOF
apiVersion: autoscaler.ccse.ctyun.cn/v1
kind: HorizontalNodeAutoscaler
metadata:
  name: example-nodepool-example-cluster-id
spec:
  disable: false
  rules:
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
      ruleName: rule1
      type: Alarm
  targetNodepools:
    - example-nodepool
  coolDown: 3
EOF
}