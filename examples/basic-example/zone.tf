
data "liquidweb_network_zone" "testing_zone" {
  name        = "Zone B"
  region_name = "US Central"
}

output "region_id" {
  value = data.liquidweb_network_zone.testing_zone.region_id
}