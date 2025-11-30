resource "ctyun_sfs" "%[1]s" {
  sfs_type     = "%[2]s"
  sfs_protocol = "%[3]s"
  name         = "%[4]s"
  sfs_size     = %[5]d
  cycle_type   = "%[6]s"
  vpc_id       = "%[7]s"
  subnet_id    = "%[8]s"
}

