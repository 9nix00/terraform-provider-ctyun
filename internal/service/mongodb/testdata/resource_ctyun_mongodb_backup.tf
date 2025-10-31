resource "ctyun_mongodb_backup" "%[1]s" {
  instance_id = "%[2]s"
  backup_name        = "%[3]s"
  description    = "%[4]s"
}

