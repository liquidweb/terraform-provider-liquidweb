resource "random_id" "dns_rec" {
  byte_length = 1
}

resource "liquidweb_network_dns_record" "testing_a_record" {
  name  = "dns-rec-${random_id.dns_rec.hex}.us-midwest-2.hostbaitor.com"
  type  = "A"
  rdata = "127.0.0.1"
  zone  = "hostbaitor.com"
}

output "domain_a_name" {
  value = liquidweb_network_dns_record.testing_a_record.name
}
