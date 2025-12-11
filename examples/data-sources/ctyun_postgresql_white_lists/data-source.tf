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

# instance id根据实例情况进行修改
data "ctyun_mysql_white_lists" "examples" {
  instance_id = "678ff914cf80469d86bdd663ee9b6377"
}

