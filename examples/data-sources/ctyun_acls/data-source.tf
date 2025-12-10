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

data "ctyun_acls" "example" {
  id         = %[2]s
  project_id = "%[3]s"
  name       = "%[4]s"
  page_no    = 1
  page_size  = 20
}

output "ctyun_acls_example" {
  value = data.ctyun_acls.example
}