resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-nat"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}
resource "ctyun_elb_health_check" "test" {
  name     = "tf-hc-for-targetgroup12"
  protocol = "TCP"
}

resource "ctyun_elb_target_group" "target_group_test" {
  name      = "tf_target_group"
  vpc_id    = ctyun_vpc.vpc_test.id
  algorithm = "wrr"
  health_check_id = ctyun_elb_health_check.test.id
  session_sticky_mode = "SOURCE_IP"
  source_ip_timeout = 30
  proxy_protocol = 1
  protocol = "TCP"
}
