resource "ctyun_oceanfs" "%[1]s" {
  sfs_protocol = "%[2]s"
  name         = "%[3]s"
  sfs_size     = "%[4]d"
  cycle_type   = "%[5]s"
  vpc_id       = "%[6]s"
  subnet_id    = "%[7]s"
  is_vpce      = "%[8]t"
}

