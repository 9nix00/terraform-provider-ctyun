

provider "ctyun" {
  env = "prod"
}

resource "ctyun_image_from_ecs" "%[1]s" {
  name         = "%[2]s"
  file_source  = "%[3]s"
  os_distro    = "%[4]s"
  os_version   = "%[5]s"
  architecture = "%[6]s"
  boot_mode    = "%[7]s"
  description  = "%[8]s"
  disk_size    = "%[9]s"
  type         = "%[10]s"
  project_id   = "%[11]s"
}
