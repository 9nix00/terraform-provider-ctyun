terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}



resource "ctyun_ccse_node_pool_scaling_policy" "example1" {
  cluster_id = "8bbcedd5bcf7426c9b444dda39512898"
  values_yaml = <<EOF
kind: HorizontalNodeAutoscaler
apiVersion: autoscaler.ccse.ctyun.cn/v1
metadata:
  name: default-8bbcedd5bcf7426c9b444dda39512898  # default-a2f1e09a01da444d989f7861a55c1c5e
spec:
  disable: false  # 禁用状态，false表示开启
  rules:
    - action:
        type: ScaleUp      # ScaleDown表示缩容，ScaleUp表示扩容
        unit: Node   # 当前仅支持Node
        value: 1    # 扩缩节点的数量
      cronTrigger:    # 当type为Cron时，传cronTrigger；type为Metric时传metricTrigger；type为Alarm时传alarmTrigger
        schedule: "8 20 * * *"   # 定时的cron表达式，建议用引号包裹
      disable: false
      ruleName: rule04130    # 每一个rule的ruleName唯一即可
      type: Cron     # 策略类型，包括Cron（定时）、Metric（监控）、Alarm（告警）
    - action:
        type: ScaleUp
        unit: Node
        value: 1
      disable: false  # 禁用状态
      metricTrigger:            # type为Metric时传metricTrigger
        metricName: cpu_util    # 指标名称，支持cpu_util、mem_util
        metricOperation: gt    # gt大于,lt：小于
        metricValue: "80"     # 阈值，用引号包裹字符串类型数值
      ruleName: rule09974
      type: Metric
    - action:
        type: ScaleUp
        unit: Node
        value: 1
      disable: false
      alarmTrigger:
        evaluationCount: 2   # 1-99 ，连续出现n次
        fun: avg  # avg（平均值）、max（最大值）、min（最小值）
        metric: cpu_util   # 磁盘、内存、网络等指标
        operator: ge  # 取值范围：eq：等于。gt：大于。ge：大于等于。lt：小于。le：小于等于
        value: "80"   # 阈值，用引号包裹
        period: 5m  # 监控周期
      ruleName: rule42857
      type: Alarm
  targetNodepools:
    - default  # 节点池的名称
  coolDown: 3  # 冷却时间，单位为min
EOF
}