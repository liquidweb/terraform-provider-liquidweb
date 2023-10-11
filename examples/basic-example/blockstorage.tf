resource "random_id" "block" {
  byte_length = 1
}

resource "liquidweb_cloud_block_storage" "testing_block_volume" {
  domain = "terraform-block${random_id.block.dec}.us-midwest-2.example.com"
  size   = 10
}

output "block_storage" {
  value = liquidweb_cloud_block_storage.testing_block_volume
}
