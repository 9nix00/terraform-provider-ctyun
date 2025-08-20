provider "ctyun" {
  env = "prod"
}

data "ctyun_ports" "%[1]s" {
}
data "ctyun_ports" "%[2]s_filtered" {
  region_id = "bb9fdb42-4eb9-45cc-8976-7a2a09412624"
  vpc_id    = "vpc-aj1ukew6eh"
}