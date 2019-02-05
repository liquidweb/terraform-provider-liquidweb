variable "liquidweb_config_path" {
  type = "string"
}

provider "liquidweb" {
  config_path = "${var.liquidweb_config_path}"
}

resource "liquidweb_storage_block_volume" "testing_some_space_balls" {
  #attach = "2GHUN4"
  domain = "spaceballz"
  size   = 5
}

output "space_balls" {
  value = "${liquidweb_storage_block_volume.testing_some_space_balls.domain}"
}
