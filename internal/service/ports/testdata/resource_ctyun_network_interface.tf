provider "ctyun" {
  env = "prod"
}


resource "ctyun_port" "%[1]s" {
  name                       = "%[2]s"
  description                = "%[3]s"
  subnet_id                  = "subnet-mph0lz50tg"
  secondary_private_ip_count = 1
}
