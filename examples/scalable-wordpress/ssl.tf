provider "acme" {
	server_url = "https://acme-v02.api.letsencrypt.org/directory"
}

resource "tls_private_key" "private_key" {
  algorithm = "RSA"
}

resource "acme_registration" "reg" {
  account_key_pem = tls_private_key.private_key.private_key_pem
  email_address   = "nobody@${var.site_name}"
}

resource "acme_certificate" "web_cert" {
  account_key_pem           = acme_registration.reg.account_key_pem
  common_name               = "${var.site_name}"
	key_type = "4096"
  # subject_alternative_names = ["www2.example.com"]

  dns_challenge {
    provider = "liquidweb"
  }
}

output "ssl" {
	value = acme_certificate.web_cert.certificate_domain
}