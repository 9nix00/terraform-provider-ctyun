provider "ctyun" {
  env = "prod"
}

# 从云主机快照创建私有镜像
resource "ctyun_image_from_ecs" "%[1]s" {
  # 必选参数：镜像名称（2~32字符，仅数字、字母、-组成，不以数字或-开头/结尾）
  image_name = "%[2]s"

  # 必选参数：快照ID（状态需为available）
  snapshot_id = "snap-1234567890abcdef"
  description = "%[3]s"

  # 可选参数：企业项目ID（默认0，即default项目）
  project_id = "0"

  # 可选参数：标签列表（最多10个，键值不可重复）
  # labels = [
  #   {
  #     label_key   = "environment"
  #     label_value = "production"
  #   },
  #   {
  #     label_key   = "department"
  #     label_value = "devops"
  #   }
  # ]

  # 可选参数：是否启用镜像完整性校验（默认false，仅部分资源池支持）
  enable_image_integrity_check = false

  # 可选参数：最小内存限制（GiB，0表示不限制，取值：0/1/2/4/8/16/32/64/128/256/512）
  minimum_ram = 2

  # 可选参数：最大内存限制（GiB，需≥最小内存，取值同上）
  maximum_ram = 16
}