

provider "ctyun" {
  alias           = "test_accpet"
  ak              = "4f996ee252e84b358a4c3f9b042f646b"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_AK
  sk              = "d22076b320d7424399e4235cb626f07a"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_SK
  env             = "prod"                                     # 如果此值不填，则默认读取环境变量中的CTYUN_ENV
}

resource "ctyun_vpc_peer_connection_attach" "%[1]s" {
  provider = ctyun.test_accpet
  peer_connection_id = %[2]s
  operation          = "%[3]s"
}

