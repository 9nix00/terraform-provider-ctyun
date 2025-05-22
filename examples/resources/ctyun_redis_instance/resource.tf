terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_redis_instance" "test" {
  cycle_type = "on_demand"
  version = "BASIC"
  edition = "DirectCluster"
  engine_version = "6.0"
  shard_count = 3
  copies_count = 2
  vpc_id = "vpc-5o8oe0oci6"
  shard_mem_size = 8
  subnet_id = "subnet-nhfs93ju2w"
  security_group_id = "sg-vrp4x1lm7p"
  instance_name = "tf-redis-3"
  password = "P@ss3edc"
  # auto_renew = true
  # auto_renew_cycle_count = 12
}