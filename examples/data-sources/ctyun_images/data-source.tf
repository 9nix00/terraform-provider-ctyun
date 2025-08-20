terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  ak        = "6e9b93befbc84d499a6d705e92ddfa5b"
  sk        = "d9896d30857f412485c6bb0179472c51"
  region_id = "bb9fdb42056f11eda1610242ac110002"
  az_name   = "cn-huadong1-jsnj1A-public-ctcloud"
  env       = "prod"
}

data "ctyun_images" "ctyun_images_test" {
  name       = "Ubuntu 22.04"
  visibility = "public"
  page_size  = 50
  page_no    = 1
}

output "ctyun_image" {
  value = data.ctyun_images.ctyun_images_test
}