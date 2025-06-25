terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

variable "password" {
  type      = string
  sensitive = true
}

# 按需
resource "ctyun_redis_instance" "tbidgqvfbs" {
  instance_name = "tf-redis-32cbppywerkb"
  version = "BASIC"
  edition = "StandardSingle"
  password = var.password
  engine_version = "7.0"
  maintenance_time = "02:00-04:00"
  protection_status = false
  shard_mem_size = 8
  vpc_id = "vpc-ewivt5nhiz"
  subnet_id = "subnet-vhyywu7mfe"
  security_group_id = "sg-ed9i3c98t2"
  cycle_type = "on_demand"
}

# 包周期
# resource "ctyun_redis_instance" "tbidgqvfbs" {
#   instance_name = "tf-redis-cbppywerkb"
#   version = "BASIC"
#   edition = "StandardSingle"
#   password = var.password
#   engine_version = "7.0"
#   maintenance_time = "02:00-04:00"
#   protection_status = false
#   shard_mem_size = 8
#   vpc_id = "vpc-ewivt5nhiz"
#   subnet_id = "subnet-vhyywu7mfe"
#   security_group_id = "sg-ed9i3c98t2"
#   cycle_type = "month"
#   cycle_count = 1
#   auto_renew = true
#   auto_renew_cycle_count = 12
# }