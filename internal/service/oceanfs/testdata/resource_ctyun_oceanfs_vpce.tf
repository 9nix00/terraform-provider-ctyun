resource "ctyun_oceanfs" "%[1]s" {
  sfs_protocol = "%[2]s"
  name         = "%[3]s"
  sfs_size     = "%[4]d"
  cycle_type   = "%[5]s"
  az_name      = "%[6]s"
  vpc_id       = "%[7]s"
  subnet_id    = "%[8]s"
  is_vpce      = "%[9]t"
}