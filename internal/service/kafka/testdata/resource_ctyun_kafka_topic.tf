resource "ctyun_kafka_topic" "%[1]s" {
  name = "%[2]s"
  prod_inst_id = "%[3]s"
  partition_num  = %[4]d
}
