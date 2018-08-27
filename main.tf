variable "storm_config_path" {
  type = "string"
}

variable "api_server_password" {
  type = "string"
}

provider "storm" {
  config_path = "${var.storm_config_path}"
}

resource "storm_server" "api_servers" {
  config_id      = 1090
  template       = "UBUNTU_1804_UNMANAGED"      // ubuntu 18.04
  domain         = "terraform.dev"
  password       = "${var.api_server_password}"
  public_ssh_key = "${file("./devkey.pub")}"
  zone           = 12
}
