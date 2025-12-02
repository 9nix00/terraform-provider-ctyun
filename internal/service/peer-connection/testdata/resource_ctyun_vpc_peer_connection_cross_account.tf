
provider "ctyun" {
  alias           = "default"
  ak              = "0e302bf7a4ce433c9763a1d8bcf9f05c"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_AK
  sk              = "cf2b7f46b9d2479fa6a20a62655635ce"
  env             = "prod"
}

resource "ctyun_vpc_peer_connection" "%[1]s" {
  # provider = ctyun.default
  project_id     = "%[2]s"
  name           = "%[3]s"
  description    = "%[4]s"
  request_vpc_id = "%[5]s"
  accept_vpc_id  = "%[6]s"
  accept_email   = "%[7]s"
}

