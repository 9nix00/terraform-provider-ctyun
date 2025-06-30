resource "ctyun_mongodb_instance" "%[1]s" {
  cycle_type        = "%[2]s"
  %[3]s // cycle_count
  %[4]s // auto_renew
  vpc_id            = "%[5]s"
  host_type         = "%[6]s"
  subnet_id         = "%[7]s"
  security_group_id = "%[8]s"
  name              = "%[9]s"
  password          = "%[10]s"
  prod_id           = "%[11]s"
  node_info_list = [%[12]s]
  %[13]s  // read_port
  %[14]s //  is_upgrade_back_up
  region_id = "bb9fdb42056f11eda1610242ac110002"
}