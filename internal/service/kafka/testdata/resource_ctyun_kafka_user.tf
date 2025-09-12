resource "ctyun_kafka_user" "%[1]s" {
  name = "%[2]s"
  prod_inst_id = "%[3]s"
  password  = "%[4]s"
  %[5]s
}
