resource "ctyun_redis_instance" "%[1]s" {
  instance_name = "%[2]s"
  version = "%[3]s"
  edition = "%[4]s"
  engine_version = "6.0"
  shard_count = 3
  copies_count = 2
  shard_mem_size = 8
  vpc_id = "%[5]s"
  subnet_id = "%[6]s"
  security_group_id = "%[7]s"
  password = "P@ss3edc"
  cycle_type = "month"
  cycle_count = 1
  auto_renew = true
  auto_renew_cycle_count = 12
}