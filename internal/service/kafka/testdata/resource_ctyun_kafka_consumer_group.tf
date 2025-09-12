resource "ctyun_kafka_consumer_group" "%[1]s" {
  name = "%[2]s"
  prod_inst_id = "%[3]s"
  description  = "%[4]s"
  %[5]s
}
