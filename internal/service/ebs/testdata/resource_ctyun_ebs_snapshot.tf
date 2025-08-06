resource "ctyun_ebs_snapshot" "%[1]s" {
  snapshot_name = "%[2]s"
  disk_id = "%[3]s"
  retention_policy = "forever"
}
