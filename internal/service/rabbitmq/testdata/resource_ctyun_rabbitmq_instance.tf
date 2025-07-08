resource "ctyun_rabbitmq_instance" "%[1]s" {
  instance_name = "%[2]s"
  cpu_num = %[3]d
  mem_size = %[4]d
  node_num = %[5]d
  zone_list = ["%[6]s"]
  disk_size = %[7]d
  disk_type = "%[8]s"
  vpc_id = "%[9]s"
  subnet_id = "%[10]s"
  security_group_id = "%[11]s"
  cycle_type = "month"
  cycle_count = 1
}
