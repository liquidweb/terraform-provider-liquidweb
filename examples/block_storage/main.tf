variable "liquidweb_config_path" {
  type = "string"
}

provider "liquidweb" {
  config_path = "${var.liquidweb_config_path}"
}

resource "liquidweb_storage_block_volume" "testing_some_space_balls" {
  attach = "TLFA7N"
  domain = "spaceballz44"
  size   = 10
}

output "space_balls" {
  value = "${liquidweb_storage_block_volume.testing_some_space_balls.domain}"
}
